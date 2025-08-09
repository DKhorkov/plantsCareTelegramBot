package handlers

import "gopkg.in/telebot.v4"

var Default = map[any]Handler{
	"/start":            Start,
	&createGroupButton:  AddGroupCallback,
	&manageGroupsButton: Test,
	&managePlantsButton: Test,
	&addFlowerButton:    Test,
	&backToStartButton:  BackToStart,
	telebot.OnText:      OnText,
	//telebot.OnPhoto:     OnPhoto,
	telebot.OnMedia:     Delete,
	telebot.OnAudio:     Delete,
	telebot.OnAnimation: Delete,
	telebot.OnBoost:     Delete,
	telebot.OnContact:   Delete,
	telebot.OnDice:      Delete,
	telebot.OnPoll:      Delete,
	telebot.OnDocument:  Delete,
	telebot.OnLocation:  Delete,
}
