package handlers

import (
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

var Default = map[any]interfaces.Handler{
	"/start":                                Start,
	"/help":                                 Help,
	&buttons.CreateGroup:                    AddGroupCallback,
	&buttons.ManageGroups:                   Delete,
	&buttons.CreatePlant:                    AddPlantCallback,
	&buttons.ManagePlants:                   ManagePlantsCallback,
	&buttons.BackToStart:                    BackToMenu,
	&buttons.BackToAddGroupTitle:            AddGroupCallback,
	&buttons.BackToAddGroupDescription:      BackToAddGroupDescriptionCallback,
	&buttons.SkipGroupDescription:           SkipGroupDescriptionCallback,
	&buttons.BackToAddGroupLastWateringDate: BackToAddGroupLastWateringDateCallback,
	&buttons.BackToAddGroupWateringInterval: BackToAddGroupWateringIntervalCallback,
	&buttons.ConfirmAddGroup:                ConfirmAddGroupCallback,
	&buttons.Menu:                           BackToMenu,
	&buttons.BackToAddPlantTitle:            AddPlantCallback,
	&buttons.BackToAddPlantDescription:      BackToAddPlantDescriptionCallback,
	&buttons.SkipPlantDescription:           SkipPlantDescriptionCallback,
	&buttons.BackToAddPlantGroup:            BackToAddPlantGroupCallback,
	&buttons.AcceptAddPlantPhoto:            AcceptAddPlantPhotoCallback,
	&buttons.RejectAddPlantPhoto:            RejectAddPlantPhotoCallback,
	&buttons.BackToAddPlantPhotoQuestion:    BackToAddPlantPhotoQuestionCallback,
	&buttons.ConfirmAddPlant:                ConfirmAddPlantCallback,
	&buttons.BackToAddPlantPhoto:            AcceptAddPlantPhotoCallback,
	&buttons.CreateAnotherPlant:             AddPlantCallback,
	&buttons.BackToManagePlantsChooseGroup:  ManagePlantsCallback,
	&buttons.BackToManagePlant:              BackToManagePlantCallback,
	&buttons.ChangePlant:                    Delete,
	&buttons.DeletePlant:                    Delete,
	telebot.OnText:                          OnText,
	telebot.OnPhoto:                         OnPhoto,
	telebot.OnMedia:                         OnMedia,
	telebot.OnAudio:                         Delete,
	telebot.OnAnimation:                     Delete,
	telebot.OnBoost:                         Delete,
	telebot.OnContact:                       Delete,
	telebot.OnDice:                          Delete,
	telebot.OnPoll:                          Delete,
	telebot.OnDocument:                      Delete,
	telebot.OnLocation:                      Delete,
}
