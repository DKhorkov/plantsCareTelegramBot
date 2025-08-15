package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
)

func OnMedia(bot *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		// Для календаря используем context.Chat().ID
		temp, err := useCases.GetUserTemporary(int(context.Chat().ID))
		if err != nil {
			return err
		}

		switch temp.Step {
		case steps.AddGroupLastWateringDateStep:
			return AddGroupLastWateringDate(bot, useCases, logger)(context)
		default:
			return Delete(bot, useCases, logger)(context)
		}
	}
}
