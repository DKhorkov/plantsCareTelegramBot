package buttons

import "gopkg.in/telebot.v4"

var (
	CreateGroupButton = telebot.InlineButton{
		Unique: "createGroup",
		Text:   "Добавить сценарий полива",
	}

	ManageGroupsButton = telebot.InlineButton{
		Unique: "manageGroups",
		Text:   "Управление сценариями полива",
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

	BackToAddGroupLastWateringDateButton = telebot.InlineButton{
		Unique: "backToAddGroupLastWateringDate",
		Text:   "Назад ↩️",
	}

	BackToAddGroupWateringIntervalButton = telebot.InlineButton{
		Unique: "backToAddGroupWateringInterval",
		Text:   "Назад ↩️",
	}

	ConfirmAddGroupButton = telebot.InlineButton{
		Unique: "confirmAddGroupButton",
		Text:   "Все верно ✅",
	}
)
