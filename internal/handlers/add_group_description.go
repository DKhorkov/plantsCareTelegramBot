package handlers

import (
	"fmt"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"gopkg.in/telebot.v4"
)

func AddGroupDescription(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

		if temp.MessageID != nil {
			err = context.Bot().Delete(&telebot.Message{ID: *temp.MessageID, Chat: context.Chat()})
			if err != nil {
				logger.Error("Failed to delete message", "Error", err)

				return err
			}
		}

		group, err := useCases.AddGroupDescription(int(context.Sender().ID), context.Message().Text)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				//{
				//	addGroupLastWateringDateCalendar,
				//},
				{
					backToAddGroupDescriptionButton,
					menuButton,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File:    telebot.FromDisk(addGroupLastWateringDateImagePath),
				Caption: fmt.Sprintf(addGroupLastWateringDateText, group.Title, group.Description, group.Title),
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

func BackToAddGroupDescriptionCallback(useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)
			return err
		}

		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

		// Получаем группу для корректного отображения данных прошлых этапов (Title):
		group, err := temp.GetGroup()
		if err != nil {
			logger.Error("Failed to get Group from Temporary", "Error", err)
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					skipGroupDescriptionButton,
				},
				{
					backToAddGroupTitleButton,
					menuButton,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File:    telebot.FromDisk(addGroupDescriptionImagePath),
				Caption: fmt.Sprintf(addGroupDescriptionText, group.Title, group.Title),
			},
			menu,
		)

		if err != nil {
			logger.Error("Failed to send message", "Error", err)
			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.GroupDescriptionStep); err != nil {
			return err
		}

		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), msg.ID); err != nil {
			return err
		}

		return nil
	}
}
