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

	AddFlowerButton = telebot.InlineButton{
		Unique: "addFlower",
		Text:   "Добавить растение",
	}

	ManagePlantsButton = telebot.InlineButton{
		Unique: "managePlants",
		Text:   "Управление растениями",
	}
)
