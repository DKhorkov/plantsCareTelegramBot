package main

import (
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/app"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/bot"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/config"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/handlers"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/middlewares"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/storage"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/usecases"
)

func main() {
	settings := config.New()
	logger := logging.New(
		settings.Logging.Level,
		settings.Logging.LogFilePath,
	)

	logger.Info("test")

	dbConnector, err := db.New(
		db.BuildDsn(settings.Database),
		settings.Database.Driver,
		logger,
		db.WithMaxOpenConnections(settings.Database.Pool.MaxOpenConnections),
		db.WithMaxIdleConnections(settings.Database.Pool.MaxIdleConnections),
		db.WithMaxConnectionLifetime(settings.Database.Pool.MaxConnectionLifetime),
		db.WithMaxConnectionIdleTime(settings.Database.Pool.MaxConnectionIdleTime),
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = dbConnector.Close()
		if err != nil {
			logger.Error("Failed to close db connections pool", "Error", err)
		}
	}()

	useCases := usecases.New(
		storage.New(dbConnector, logger),
		logger,
	)

	b, err := bot.New(settings.Bot.Token, settings.Bot.PollTimeout)
	if err != nil {
		panic(err)
	}

	b.Use(middlewares.Logging(logger))
	handlers.Prepare(b, useCases, logger, handlers.Default)

	application := app.New(b)
	application.Run()
}
