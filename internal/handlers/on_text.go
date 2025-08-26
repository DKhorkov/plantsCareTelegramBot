package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
)

func OnText(bot *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			// Ошибка уже заллогирована, удаляем сообщение.
			// Может быть, когда пользователь не жал /start и написал что-то боту:
			return Delete(bot, useCases, logger)(context)
		}

		switch temp.Step {
		case steps.AddGroupTitle:
			return AddGroupTitle(bot, useCases, logger)(context)
		case steps.AddGroupDescription:
			return AddGroupDescription(bot, useCases, logger)(context)
		case steps.ChangeGroupTitle:
			return ChangeGroupTitle(bot, useCases, logger)(context)
		case steps.ChangeGroupDescription:
			return ChangeGroupDescription(bot, useCases, logger)(context)
		case steps.AddPlantTitle:
			return AddPlantTitle(bot, useCases, logger)(context)
		case steps.AddPlantDescription:
			return AddPlantDescription(bot, useCases, logger)(context)
		case steps.ChangePlantTitle:
			return ChangePlantTitle(bot, useCases, logger)(context)
		case steps.ChangePlantDescription:
			return ChangePlantDescription(bot, useCases, logger)(context)
		default:
			return Delete(bot, useCases, logger)(context)
		}
	}
}
