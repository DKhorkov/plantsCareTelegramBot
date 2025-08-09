package handlers

import "gopkg.in/telebot.v4"

var (
	createGroupButton = telebot.InlineButton{
		Unique: "createGroup",
		Text:   "Добавить сценарий полива",
	}

	manageGroupsButton = telebot.InlineButton{
		Unique: "manageGroups",
		Text:   "Управление сценариями полива",
	}

	addFlowerButton = telebot.InlineButton{
		Unique: "addFlower",
		Text:   "Добавить растение",
	}

	managePlantsButton = telebot.InlineButton{
		Unique: "managePlants",
		Text:   "Управление растениями",
	}

	backToStartButton = telebot.InlineButton{
		Unique: "backToStart",
		Text:   "Назад ↩️",
	}
)
