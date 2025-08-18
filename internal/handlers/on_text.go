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
			return err
		}

		switch temp.Step {
		case steps.AddGroupTitleStep:
			return AddGroupTitle(bot, useCases, logger)(context)
		case steps.AddGroupDescriptionStep:
			return AddGroupDescription(bot, useCases, logger)(context)
		case steps.AddPlantTitleStep:
			return AddPlantTitle(bot, useCases, logger)(context)
		case steps.AddPlantDescriptionStep:
			return AddPlantDescription(bot, useCases, logger)(context)
		default:
			return Delete(bot, useCases, logger)(context)
		}
	}
}
