package buttons

import "gopkg.in/telebot.v4"

var (
	CreateGroup = telebot.InlineButton{
		Unique: "createGroup",
		Text:   "Добавить сценарий полива",
	}

	ManageGroups = telebot.InlineButton{
		Unique: "manageGroups",
		Text:   "Управление сценариями полива",
	}

	SkipGroupDescription = telebot.InlineButton{
		Unique: "skipGroupDescription",
		Text:   "Пропустить",
	}

	BackToAddGroupTitle = telebot.InlineButton{
		Unique: "backToAddGroupTitle",
		Text:   "Назад ↩️",
	}

	BackToAddGroupDescription = telebot.InlineButton{
		Unique: "backToAddGroupDescription",
		Text:   "Назад ↩️",
	}

	BackToAddGroupLastWateringDate = telebot.InlineButton{
		Unique: "backToAddGroupLastWateringDate",
		Text:   "Назад ↩️",
	}

	BackToAddGroupWateringInterval = telebot.InlineButton{
		Unique: "backToAddGroupWateringInterval",
		Text:   "Назад ↩️",
	}

	ConfirmAddGroup = telebot.InlineButton{
		Unique: "confirmAddGroupButton",
		Text:   "Все верно ✅",
	}
)
