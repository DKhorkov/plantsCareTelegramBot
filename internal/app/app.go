package app

import (
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/telebot.v4"
)

type App struct {
	bot *telebot.Bot
}

func New(bot *telebot.Bot) *App {
	return &App{
		bot: bot,
	}
}

func (application *App) Run() {
	// Launch asynchronous for graceful shutdown purpose:
	go application.bot.Start()

	// Graceful shutdown. When system signal will be received, signal.Notify function will write it to channel.
	// After this event, main goroutine will be unblocked (<-stopChannel blocks it) and application will be
	// gracefully stopped:
	stopChannel := make(chan os.Signal, 1)
	signal.Notify(stopChannel, syscall.SIGINT, syscall.SIGTERM)
	<-stopChannel
	application.bot.Stop()
}
