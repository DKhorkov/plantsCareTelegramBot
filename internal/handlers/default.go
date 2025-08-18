package handlers

import (
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

var Default = map[any]interfaces.Handler{
	"/start":                                      Start,
	&buttons.CreateGroupButton:                    AddGroupCallback,
	&buttons.ManageGroupsButton:                   Delete,
	&buttons.CreatePlantButton:                    AddPlantCallback,
	&buttons.ManagePlantsButton:                   Delete,
	&buttons.BackToStartButton:                    BackToMenu,
	&buttons.BackToAddGroupTitleButton:            AddGroupCallback,
	&buttons.BackToAddGroupDescriptionButton:      BackToAddGroupDescriptionCallback,
	&buttons.SkipGroupDescriptionButton:           SkipGroupDescriptionCallback,
	&buttons.BackToAddGroupLastWateringDateButton: BackToAddGroupLastWateringDateCallback,
	&buttons.BackToAddGroupWateringIntervalButton: BackToAddGroupWateringIntervalCallback,
	&buttons.ConfirmAddGroupButton:                ConfirmAddGroupCallback,
	&buttons.MenuButton:                           BackToMenu,
	&buttons.BackToAddPlantTitleButton:            AddPlantCallback,
	&buttons.BackToAddPlantDescriptionButton:      BackToAddPlantDescriptionCallback,
	&buttons.SkipPlantDescriptionButton:           SkipPlantDescriptionCallback,
	telebot.OnText:                                OnText,
	// telebot.OnPhoto:     OnPhoto,
	telebot.OnMedia:     OnMedia,
	telebot.OnAudio:     Delete,
	telebot.OnAnimation: Delete,
	telebot.OnBoost:     Delete,
	telebot.OnContact:   Delete,
	telebot.OnDice:      Delete,
	telebot.OnPoll:      Delete,
	telebot.OnDocument:  Delete,
	telebot.OnLocation:  Delete,
}
