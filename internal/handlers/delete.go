package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"gopkg.in/telebot.v4"
)

func Delete(_ interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		err := context.Delete()
		if err != nil {
			logger.Error("Failed to delete message", "Error", err)
		}

		return err
	}
}
