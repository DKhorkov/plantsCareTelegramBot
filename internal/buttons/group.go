package buttons

import (
	"gopkg.in/telebot.v4"
)

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

	ManageGroupSeePlants = telebot.InlineButton{
		Unique: "manageGroupSeePlants",
		Text:   "Просмотр растений в данном сценарии 👀",
	}

	ManageGroupChange = telebot.InlineButton{
		Unique: "manageGroupChange",
		Text:   "Редактировать сценарий полива 🛠",
	}

	ManageGroupRemoval = telebot.InlineButton{
		Unique: "manageGroupRemoval",
		Text:   "Удалить сценарий полива 🗑",
	}

	BackToManageGroup = telebot.InlineButton{
		Unique: "backToManageGroup",
		Text:   "Назад ↩️",
	}

	ConfirmGroupRemoval = telebot.InlineButton{
		Unique: "confirmGroupRemoval",
		Text:   "Подтвердить удаление ✅",
	}

	BackToManageGroupAction = telebot.InlineButton{
		Unique: "backToManageGroupAction",
		Text:   "Назад ↩️",
	}

	BackToManageGroupChange = telebot.InlineButton{
		Unique: "backToManageGroupChange",
		Text:   "Назад ↩️",
	}

	ManageGroupChangeTitle = telebot.InlineButton{
		Unique: "manageGroupChangeTitle",
		Text:   "Изменить название сценария полива",
	}

	ManageGroupChangeDescription = telebot.InlineButton{
		Unique: "manageGroupChangeDescription",
		Text:   "Изменить описание сценария полива",
	}

	ManageGroupChangeLastWateringDate = telebot.InlineButton{
		Unique: "manageGroupChangeLastWateringDate",
		Text:   "Изменить дату последнего полива сценария",
	}

	ManageGroupChangeWateringInterval = telebot.InlineButton{
		Unique: "manageGroupChangeWateringInterval",
		Text:   "Изменить интервал полива сценария",
	}

	GroupWatered = telebot.InlineButton{
		Unique: "groupWatered",
		Text:   "Растения в данном сценарии политы ✅",
	}

	ManageGroup = telebot.InlineButton{
		Unique: "manageGroup",
	}

	AddGroupWateringInterval = telebot.InlineButton{
		Unique: "addGroupWateringInterval",
	}

	ChangeGroupWateringInterval = telebot.InlineButton{
		Unique: "changeGroupWateringInterval",
	}
)
