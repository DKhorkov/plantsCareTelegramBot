package cron

import (
	"time"

	"github.com/DKhorkov/libs/logging"

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

func (c *Cron) Run() error {
	go func() {
		if r := recover(); r != nil {
			c.logger.Error("Recovered from panic", "Recovered", r)
		}
	}()

	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			return nil
		case <-ticker.C:
			if err := c.callback(); err != nil {
				return err
			}
		}
	}
}

func (c *Cron) Stop() error {
	c.stopChan <- struct{}{}

	return nil
}
