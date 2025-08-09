package middlewares

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"
)

func Logging(logger logging.Logger) func(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			if c.Callback() != nil {
				logger.Info(
					"Received new callback",
					"From",
					c.Message().Sender.ID,
					"Unique",
					c.Callback().Unique,
					"Data",
					c.Callback().Data,
				)

				return next(c) // continue execution chain
			}

			logger.Info(
				"Received new message",
				"From",
				c.Message().Sender.ID,
				"Message",
				c.Text(),
			)

			return next(c) // continue execution chain
		}
	}
}
