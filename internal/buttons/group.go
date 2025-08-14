package buttons

import "gopkg.in/telebot.v4"

var (
	BackToStartButton = telebot.InlineButton{
		Unique: "backToStart",
		Text:   "Назад ↩️",
	}

	SkipGroupDescriptionButton = telebot.InlineButton{
		Unique: "skipGroupDescription",
		Text:   "Пропустить",
	}

	BackToAddGroupTitleButton = telebot.InlineButton{
		Unique: "backToAddGroupTitle",
		Text:   "Назад ↩️",
	}

	BackToAddGroupDescriptionButton = telebot.InlineButton{
		Unique: "backToAddGroupDescription",
		Text:   "Назад ↩️",
	}
)
