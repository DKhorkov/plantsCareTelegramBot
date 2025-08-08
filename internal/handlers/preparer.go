package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"gopkg.in/telebot.v4"
)

func Prepare(bot *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger, handlers map[any]Handler) {
	for cmd, h := range handlers {
		bot.Handle(cmd, h(useCases, logger))
	}

	//bot.Handle(&clb, Test)
	bot.Handle(telebot.OnMedia, func(context telebot.Context) error {
		return context.Send("invalid")
	})
}
