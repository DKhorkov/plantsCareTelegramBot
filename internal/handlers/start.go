package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"gopkg.in/telebot.v4"
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
		)

		if err != nil {
			return err
		}

		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.StartStep); err != nil {
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
				File:    telebot.FromDisk(startImagePath),
				Caption: startMessageText,
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

func AddGroupCallback(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)
			return err
		}

		if err := useCases.SetTemporaryStep(int(context.Sender().ID), steps.GroupTitleStep); err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					backToStartButton,
				},
			},
		}

		err := context.Send(
			&telebot.Photo{
				File:    telebot.FromDisk(addGroupTitleImagePath),
				Caption: addGroupTitleText,
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
		//err := context.Delete()
		//if err != nil {
		//	return err
		//}
		//
		//return context.Send("some message")

		return context.Respond(&telebot.CallbackResponse{
			Text: "Ошибка обновления: ",
		})
	}
}
