package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

const (
	loggingTraceSkipLevel = 1
)

type App struct {
	bot    interfaces.Bot
	logger logging.Logger
	crons  []interfaces.Cron
}

func New(
	bot interfaces.Bot,
	logger logging.Logger,
	crons []interfaces.Cron,
) *App {
	return &App{
		bot:    bot,
		logger: logger,
		crons:  crons,
	}
}

func (application *App) Run() {
	// Launch asynchronous for graceful shutdown purpose:
	go application.bot.Start()

	// Запускаем кроны для отправки уведомлений
	wg := new(sync.WaitGroup)
	for _, cron := range application.crons {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if err := cron.Run(); err != nil {
				application.logger.Error(
					"Error running cron job",
					"Error", err,
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)
			}
		}()
	}

	// Graceful shutdown. When system signal will be received, signal.Notify function will write it to channel.
	// After this event, main goroutine will be unblocked (<-stopChannel blocks it) and application will be
	// gracefully stopped:
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)
	<-stopChannel

	// Убиваем кроны до бота, чтобы не резать отправки:
	for _, cron := range application.crons {
		if err := cron.Stop(); err != nil {
			application.logger.Error("Error stopping cron job", "Error", err)
		}
	}

	// Дожидаемся остановки всех горутин:
	wg.Wait()

	// Останавливаем приложение:
	application.bot.Stop()
}
