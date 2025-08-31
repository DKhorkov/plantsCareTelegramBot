package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
)

func OnMedia(bot interfaces.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		// Для календаря используем context.Chat().ID:
		temp, err := useCases.GetUserTemporary(int(context.Chat().ID))
		if err != nil {
			// Ошибка уже заллогирована, удаляем сообщение.
			// Может быть, когда пользователь не жал /start и отправил что-то боту:
			return Delete(bot, useCases, logger)(context)
		}

		switch temp.Step {
		case steps.AddGroupLastWateringDate:
			return AddGroupLastWateringDate(bot, useCases, logger)(context)
		default:
			return Delete(bot, useCases, logger)(context)
		}
	}
}
