package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"gopkg.in/telebot.v4"
)

func BackToMenu(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)
			return err
		}

		user, err := useCases.GetUserByTelegramID(int(context.Sender().ID))
		if err != nil {
			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
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

		groupsCount, err := useCases.CountUserGroups(user.ID)
		if err != nil {
			return err
		}

		if groupsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{addFlowerButton})
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{manageGroupsButton})
		}

		plantsCount, err := useCases.CountUserPlants(user.ID)
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
