package buttons

import "gopkg.in/telebot.v4"

var (
	CreatePlantButton = telebot.InlineButton{
		Unique: "addPlant",
		Text:   "Добавить растение",
	}

	ManagePlantsButton = telebot.InlineButton{
		Unique: "managePlants",
		Text:   "Управление растениями",
	}

	SkipPlantDescriptionButton = telebot.InlineButton{
		Unique: "skipPlantDescription",
		Text:   "Пропустить",
	}

	BackToAddPlantTitleButton = telebot.InlineButton{
		Unique: "backToAddPlantTitle",
		Text:   "Назад ↩️",
	}

	BackToAddPlantDescriptionButton = telebot.InlineButton{
		Unique: "backToAddPlantDescription",
		Text:   "Назад ↩️",
	}

	BackToAddPlantGroupButton = telebot.InlineButton{
		Unique: "backToAddPlantGroup",
		Text:   "Назад ↩️",
	}

	AcceptAddPlantPhotoButton = telebot.InlineButton{
		Unique: "acceptAddPlantPhoto",
		Text:   "Да️",
	}

	RejectAddPlantPhotoButton = telebot.InlineButton{
		Unique: "rejectAddPlantPhoto",
		Text:   "Нет",
	}

	BackToAddPlantPhotoQuestionButton = telebot.InlineButton{
		Unique: "backToAddPlantPhotoQuestion",
		Text:   "Назад ↩️",
	}
)
