package buttons

import "gopkg.in/telebot.v4"

var (
	CreatePlant = telebot.InlineButton{
		Unique: "addPlant",
		Text:   "Добавить растение",
	}

	ManagePlants = telebot.InlineButton{
		Unique: "managePlants",
		Text:   "Управление растениями",
	}

	SkipPlantDescription = telebot.InlineButton{
		Unique: "skipPlantDescription",
		Text:   "Пропустить",
	}

	BackToAddPlantTitle = telebot.InlineButton{
		Unique: "backToAddPlantTitle",
		Text:   "Назад ↩️",
	}

	BackToAddPlantDescription = telebot.InlineButton{
		Unique: "backToAddPlantDescription",
		Text:   "Назад ↩️",
	}

	BackToAddPlantGroup = telebot.InlineButton{
		Unique: "backToAddPlantGroup",
		Text:   "Назад ↩️",
	}

	AcceptAddPlantPhoto = telebot.InlineButton{
		Unique: "acceptAddPlantPhoto",
		Text:   "Да️",
	}

	RejectAddPlantPhoto = telebot.InlineButton{
		Unique: "rejectAddPlantPhoto",
		Text:   "Нет",
	}

	BackToAddPlantPhotoQuestion = telebot.InlineButton{
		Unique: "backToAddPlantPhotoQuestion",
		Text:   "Назад ↩️",
	}

	ConfirmAddPlant = telebot.InlineButton{
		Unique: "confirmAddPlant",
		Text:   "Все верно ✅",
	}

	BackToAddPlantPhoto = telebot.InlineButton{
		Unique: "backToAddPlantPhoto",
		Text:   "Назад ↩️",
	}

	CreateAnotherPlant = telebot.InlineButton{
		Unique: "createAnotherPlant",
		Text:   "Добавить еще одно растение",
	}

	BackToManagePlantsChooseGroup = telebot.InlineButton{
		Unique: "backToManagePlantsChooseGroup",
		Text:   "Назад ↩️",
	}

	BackToManagePlant = telebot.InlineButton{
		Unique: "backToManagePlant",
		Text:   "Назад ↩️",
	}

	ChangePlant = telebot.InlineButton{
		Unique: "changePlant",
		Text:   "Редактировать растение",
	}

	DeletePlant = telebot.InlineButton{
		Unique: "deletePlant",
		Text:   "Удалить растение",
	}
)
