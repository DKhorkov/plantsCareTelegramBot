package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/config"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

type App struct {
	bot                 *telebot.Bot
	logger              logging.Logger
	notificationsCron   interfaces.Cron
	notificationsConfig config.NotificationsConfig
}

func New(
	bot *telebot.Bot,
	logger logging.Logger,
	notificationsCron interfaces.Cron,
	notificationsConfig config.NotificationsConfig,
) *App {
	return &App{
		bot:                 bot,
		logger:              logger,
		notificationsCron:   notificationsCron,
		notificationsConfig: notificationsConfig,
	}
}

func (application *App) Run() {
	// Launch asynchronous for graceful shutdown purpose:
	go application.bot.Start()

	// Запускаем кроны для отправки уведомлений
	wg := new(sync.WaitGroup)
	for i := range application.notificationsConfig.CronsCount {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := application.notificationsCron.Run(
				application.notificationsConfig.GroupsLimitPerQuery,
				application.notificationsConfig.GroupsLimitPerQuery*i,
				application.notificationsConfig.CronCheckInterval,
			)
			if err != nil {
				application.logger.Error("Error running cron job", "Error", err)
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
	if err := application.notificationsCron.Stop(); err != nil {
		application.logger.Error("Error stopping cron job", "Error", err)
	}

	// Дожидаемся остановки всех горутин:
	wg.Wait()

	// Останавливаем приложение:
	application.bot.Stop()
}
