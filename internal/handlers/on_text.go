package handlers

import (
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
)

func OnText(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

		switch temp.Step {
		case steps.GroupTitleStep:
			return AddGroupTitle(useCases, logger)(context)
		case steps.GroupDescriptionStep:
		case steps.PlantTitleStep:
		case steps.PlantDescriptionStep:
		default:
			return Delete(useCases, logger)(context)
		}

		panic("unreachable")
	}
}
