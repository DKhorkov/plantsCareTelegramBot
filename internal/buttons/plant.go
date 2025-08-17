package buttons

import "gopkg.in/telebot.v4"

var (
	AddPlantButton = telebot.InlineButton{
		Unique: "addPlant",
		Text:   "Добавить растение",
	}

	ManagePlantsButton = telebot.InlineButton{
		Unique: "managePlants",
		Text:   "Управление растениями",
	}
)
