package interfaces

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"
)

type Handler func(bot Bot, useCases UseCases, logger logging.Logger) telebot.HandlerFunc
