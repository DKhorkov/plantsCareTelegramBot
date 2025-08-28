package buttons

import (
	"gopkg.in/telebot.v4"
)

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

	ManagePlantChange = telebot.InlineButton{
		Unique: "managePlantChange",
		Text:   "Редактировать растение 🛠",
	}

	ManagePlantRemoval = telebot.InlineButton{
		Unique: "managePlantRemoval",
		Text:   "Удалить растение 🗑",
	}

	ConfirmPlantRemoval = telebot.InlineButton{
		Unique: "confirmPlantRemoval",
		Text:   "Подтвердить удаление ✅",
	}

	BackToManagePlantAction = telebot.InlineButton{
		Unique: "backToManagePlantAction",
		Text:   "Назад ↩️",
	}

	ManagePlantChangeTitle = telebot.InlineButton{
		Unique: "managePlantChangeTitle",
		Text:   "Изменить название растения",
	}

	ManagePlantChangeDescription = telebot.InlineButton{
		Unique: "managePlantChangeDescription",
		Text:   "Изменить заметки по растению",
	}

	ManagePlantChangeGroup = telebot.InlineButton{
		Unique: "managePlantChangeGroup",
		Text:   "Изменить сценарий полива растения",
	}

	ManagePlantChangePhoto = telebot.InlineButton{
		Unique: "managePlantChangePhoto",
		Text:   "Изменить фотографию растения",
	}

	BackToManagePlantChange = telebot.InlineButton{
		Unique: "backToManagePlantChange",
		Text:   "Назад ↩️",
	}

	ManagePlantsGroup = telebot.InlineButton{
		Unique: "managePlantsGroup",
	}

	AddPlantGroup = telebot.InlineButton{
		Unique: "addPlantGroup",
	}

	ChangePlantGroup = telebot.InlineButton{
		Unique: "changePlantGroup",
	}

	ManagePlant = telebot.InlineButton{
		Unique: "managePlant",
	}
)
