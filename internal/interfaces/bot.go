package interfaces

import (
	"gopkg.in/telebot.v4"
)

//go:generate mockgen -source=bot.go -destination=../../mocks/bot/bot.go -package=mockbot
type Bot interface {
	Handle(endpoint any, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc)
	Start()
	Stop()
	ProcessUpdate(u telebot.Update)
	telebot.API
}

type Context interface {
	telebot.Context
}
