package handlers

import (
	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

func BackToMenu(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
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
		if err = useCases.ResetTemporary(int(context.Sender().ID)); err != nil {
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

		groupsCount, err := useCases.CountUserGroups(user.ID)
		if err != nil {
			return err
		}

		if groupsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.CreatePlantButton})
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.ManageGroupsButton})
		}

		plantsCount, err := useCases.CountUserPlants(user.ID)
		if err != nil {
			return err
		}

		if plantsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.ManageGroupsButton})
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
