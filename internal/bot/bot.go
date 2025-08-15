package bot

import (
	"time"

	"gopkg.in/telebot.v4"
)

func New(token string, pollTimeout time.Duration) (*telebot.Bot, error) {
	cfg := telebot.Settings{
		Token:     token,
		Poller:    &telebot.LongPoller{Timeout: pollTimeout},
		ParseMode: telebot.ModeHTML,
	}

	return telebot.NewBot(cfg)
}
