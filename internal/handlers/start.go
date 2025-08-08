package handlers

import (
	"fmt"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"gopkg.in/telebot.v4"
)

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
)

func Start(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete /start message", "Error", err)
			return err
		}

		userID, err := useCases.SaveUser(
			entities.User{
				TelegramID: int(context.Sender().ID),
				Username:   context.Sender().Username,
				Firstname:  context.Sender().FirstName,
				Lastname:   context.Sender().LastName,
				IsBot:      context.Sender().IsBot,
			},
			context.Message().ID,
		)

		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					createGroupButton,
				},
			},
		}

		groupsCount, err := useCases.CountUserGroups(userID)
		if err != nil {
			return err
		}

		if groupsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{addFlowerButton})
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{manageGroupsButton})
		}

		plantsCount, err := useCases.CountUserPlants(userID)
		if err != nil {
			return err
		}

		if plantsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{managePlantsButton})
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(startImagePath),
				Caption: fmt.Sprintf("Привет!\n" +
					"Кажется, ты хочешь отрегулировать полив растений😃\n" +
					"Я тебе помогу.\n\n" +
					"Выбери действие ниже, чтобы мы могли продолжить:\n\n",
				),
			},
			menu,
		)

		if err != nil {
			logger.Error("Failed to send message", "Error", err)
			return err
		}

		return nil
	}
}

func Test(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		err := context.Delete()
		if err != nil {
			return err
		}

		return context.Send("some message")
	}
}
