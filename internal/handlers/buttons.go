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

	skipGroupDescriptionButton = telebot.InlineButton{
		Unique: "skipGroupDescription",
		Text:   "Пропустить",
	}

	backToAddGroupTitleButton = telebot.InlineButton{
		Unique: "backToAddGroupTitle",
		Text:   "Назад ↩️",
	}

	menuButton = telebot.InlineButton{
		Unique: "menu",
		Text:   "В меню 🏠",
	}

	backToAddGroupDescriptionButton = telebot.InlineButton{
		Unique: "backToAddGroupDescription",
		Text:   "Назад ↩️",
	}
)
