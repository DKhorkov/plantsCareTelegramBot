package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
)

// OnPhoto После добавления бот считает колбэк на калндарь как фото, а не медиа.
func OnPhoto(bot interfaces.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		// Для календаря используем context.Chat().ID. В случае с фото будет равен context.Sender().ID:
		temp, err := useCases.GetUserTemporary(int(context.Chat().ID))
		if err != nil {
			// Ошибка уже заллогирована, удаляем сообщение.
			// Может быть, когда пользователь не жал /start и отправил что-то боту:
			return Delete(bot, useCases, logger)(context)
		}

		switch temp.Step {
		case steps.AddGroupLastWateringDate: // Логика обработки ответа от календаря с сообщением с картинкой
			return AddGroupLastWateringDate(bot, useCases, logger)(context)
		case steps.ChangeGroupLastWateringDate: // Логика обработки ответа от календаря с сообщением с картинкой
			return ChangeGroupLastWateringDate(bot, useCases, logger)(context)
		case steps.AddPlantPhoto:
			return AddPlantPhoto(bot, useCases, logger)(context)
		case steps.ChangePlantPhoto:
			return ChangePlantPhoto(bot, useCases, logger)(context)
		default:
			return Delete(bot, useCases, logger)(context)
		}
	}
}
