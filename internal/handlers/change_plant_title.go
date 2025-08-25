package handlers

import (
	"bytes"
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

func ChangePlantTitle(
	_ *telebot.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error(
				"Failed to delete message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		if len(context.Message().Text) > plantTitleMaxLength {
			if err := context.Send(fmt.Sprintf(texts.PlantTitleTooLong, plantTitleMaxLength)); err != nil {
				logger.Error(
					"Failed to send message",
					"Error", err,
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)

				return err
			}

			return nil
		}

		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

		plant, err := temp.GetPlant()
		if err != nil {
			logger.Error(
				"Failed to get Plant from Temporary",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		plant, err = useCases.GetPlant(plant.ID)
		if err != nil {
			return err
		}

		plant.Title = context.Message().Text

		plantExists, err := useCases.PlantExists(*plant)
		if err != nil {
			logger.Error(
				fmt.Sprintf("Failed to check Plant existence for user with telegramId=%d", context.Sender().ID),
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		group, err := useCases.GetGroup(plant.GroupID)
		if err != nil {
			return err
		}

		if plantExists {
			// Нет context.Callback() для обычного сообщения, поэтому отправляем ответ текстом:
			if err = context.Send(fmt.Sprintf(texts.PlantAlreadyExists, group.Title)); err != nil {
				logger.Error(
					"Failed to send message",
					"Error", err,
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)

				return err
			}

			return nil
		}

		// Удаляем сообщение только если нет идентичного растения для группы:
		if temp.MessageID != nil {
			err = context.Bot().Delete(&telebot.Message{ID: *temp.MessageID, Chat: context.Chat()})
			if err != nil {
				logger.Error(
					"Failed to delete message",
					"Error", err,
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)

				return err
			}
		}

		if err = useCases.UpdatePlant(*plant); err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.ManagePlantChangeTitle,
				},
				{
					buttons.ManagePlantChangeDescription,
				},
				{
					buttons.ManagePlantChangeGroup,
				},
				{
					buttons.ManagePlantChangePhoto,
				},
				{
					buttons.BackToManagePlantAction,
					buttons.Menu,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromReader(bytes.NewReader(plant.Photo)),
				Caption: fmt.Sprintf(
					texts.ManagePlantChange,
					plant.Title,
					plant.Description,
					group.Title,
				),
			},
			menu,
		)
		if err != nil {
			logger.Error(
				"Failed to send message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.ManagePlantChange); err != nil {
			return err
		}

		return nil
	}
}
