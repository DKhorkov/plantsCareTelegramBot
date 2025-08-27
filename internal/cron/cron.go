package cron

import (
	"time"

	"github.com/DKhorkov/libs/logging"

	customerrors "github.com/DKhorkov/plantsCareTelegramBot/internal/errors"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

type Cron struct {
	logger   logging.Logger
	callback interfaces.Callback
	interval time.Duration
	stopChan chan struct{}
}

func New(logger logging.Logger, callback interfaces.Callback, interval time.Duration) *Cron {
	return &Cron{
		logger:   logger,
		callback: callback,
		interval: interval,
		stopChan: make(chan struct{}),
	}
}

func (c *Cron) Run() (err error) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return nil
		case <-ticker.C:
			// recover() работает ТОЛЬКО в той же горутине, где произошла паника.
			// Оборачиваем вызов callback в recover.
			func() {
				defer func() {
					if r := recover(); r != nil {
						c.logger.Error("Recovered from panic", "Recovered", r)

						err = customerrors.ErrPanic
					}
				}()

				if err = c.callback(); err != nil {
					return
				}
			}()

			// Возвращаем ошибку только если она существует, чтобы не прерывать цикл крона:
			if err != nil {
				return err
			}
		}
	}
}

// Stop останавливает крона.
// Лучше использовать select чтобы не было блокировки при повторном вызове.
func (c *Cron) Stop() error {
	select {
	case <-c.stopChan:
		// Уже остановлен
		return nil
	default:
		close(c.stopChan) // Отправляем сигнал остановки
	}

	return nil
}
