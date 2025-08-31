package main

import (
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/loadenv"
	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/app"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/bot"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/config"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/cron"
	cronPreparers "github.com/DKhorkov/plantsCareTelegramBot/internal/cron/preparers"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/handlers"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/middlewares"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/storage"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/usecases"
)

const (
	loggingTraceSkipLevel = 1
)

func main() {
	// Инициализируем переменные окружения для дальнейшего считывания в конфиге:
	loadenv.Init()

	cfg := config.New()
	logger := logging.New(
		cfg.Logging.Level,
		cfg.Logging.LogFilePath,
	)

	dbConnector, err := db.New(
		db.BuildDsn(cfg.Database),
		cfg.Database.Driver,
		logger,
		db.WithMaxOpenConnections(cfg.Database.Pool.MaxOpenConnections),
		db.WithMaxIdleConnections(cfg.Database.Pool.MaxIdleConnections),
		db.WithMaxConnectionLifetime(cfg.Database.Pool.MaxConnectionLifetime),
		db.WithMaxConnectionIdleTime(cfg.Database.Pool.MaxConnectionIdleTime),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = dbConnector.Close()
		if err != nil {
			logger.Error(
				"Failed to close db connections pool",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}
	}()

	useCases := usecases.New(
		storage.New(dbConnector, logger),
		logger,
	)

	b, err := bot.New(cfg.Bot.Token, cfg.Bot.PollTimeout)
	if err != nil {
		panic(err)
	}

	b.Use(middlewares.Logging(logger))
	handlers.Prepare(b, useCases, logger, handlers.Default)

	// Setup crons:
	var crons []interfaces.Cron

	for i := range cfg.Notifications.CronsCount {
		callback := cronPreparers.NewNotificationsPreparer(
			b,
			useCases,
			logger,
			cfg.Notifications.GroupsLimitPerQuery,
			cfg.Notifications.GroupsLimitPerQuery*i,
		).GetCallback()

		crons = append(
			crons,
			cron.New(
				logger,
				callback,
				cfg.Notifications.CronCheckInterval,
			),
		)
	}

	application := app.New(b, logger, crons)
	application.Run()
}
