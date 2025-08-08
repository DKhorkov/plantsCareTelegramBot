package bot

import (
	"gopkg.in/telebot.v4"
	"time"
)

func New(token string, pollTimeout time.Duration) (*telebot.Bot, error) {
	botConfig := telebot.Settings{
		Token:     token,
		Poller:    &telebot.LongPoller{Timeout: pollTimeout},
		ParseMode: telebot.ModeHTML,
	}

	return telebot.NewBot(botConfig)
}
