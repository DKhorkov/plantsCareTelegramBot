package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

func Start(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
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

		// TODO при проблемах логики следует сделать в рамках транзакции.
		// TODO Тут повторяем вне юзкейсов, чтобы работало даже вне повторной регистрации.
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.StartStep); err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.CreateGroupButton,
				},
			},
		}

		groupsCount, err := useCases.CountUserGroups(userID)
		if err != nil {
			return err
		}

		if groupsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.AddFlowerButton})
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.ManageGroupsButton})
		}

		plantsCount, err := useCases.CountUserPlants(userID)
		if err != nil {
			return err
		}

		if plantsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.ManagePlantsButton})
		}

		err = context.Send(
			&telebot.Photo{
				File:    telebot.FromDisk(paths.StartImagePath),
				Caption: texts.StartMessageText,
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

func AddGroupCallback(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err := useCases.SetTemporaryStep(int(context.Sender().ID), steps.GroupTitleStep); err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.BackToStartButton,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File:    telebot.FromDisk(paths.AddGroupTitleImagePath),
				Caption: texts.AddGroupTitleText,
			},
			menu,
		)
		if err != nil {
			logger.Error("Failed to send message", "Error", err)

			return err
		}

		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), msg.ID); err != nil {
			return err
		}

		return nil
	}
}
