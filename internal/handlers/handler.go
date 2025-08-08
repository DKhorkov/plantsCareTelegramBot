package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"gopkg.in/telebot.v4"
)

type Handler = func(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc
