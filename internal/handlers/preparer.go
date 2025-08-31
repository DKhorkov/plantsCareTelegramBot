package handlers

import (
	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

func Prepare(
	bot interfaces.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
	handlers map[any]interfaces.Handler,
) {
	for cmd, h := range handlers {
		bot.Handle(cmd, h(bot, useCases, logger))
	}
}
