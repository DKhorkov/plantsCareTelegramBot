package interfaces

import (
	"io"

	"gopkg.in/telebot.v4"
)

//go:generate mockgen -source=bot.go -destination=../../mocks/bot/bot.go -package=mockbot
type Bot interface {
	Handle(endpoint any, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc)
	Start()
	Stop()
	Send(to telebot.Recipient, what any, opts ...any) (*telebot.Message, error)
	EditReplyMarkup(msg telebot.Editable, markup *telebot.ReplyMarkup) (*telebot.Message, error)
	ProcessUpdate(u telebot.Update)
	File(file *telebot.File) (io.ReadCloser, error)
}
