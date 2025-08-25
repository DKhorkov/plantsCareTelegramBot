package handlers

import (
	"bytes"
	"fmt"
	"io"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

func ChangePlantPhoto(
	bot *telebot.Bot,
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

		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

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

		photoReader, err := bot.File(context.Message().Photo.MediaFile())
		if err != nil {
			return err
		}

		defer func() {
			if err = photoReader.Close(); err != nil {
				logger.Error(
					"Failed to close photo reader",
					"Error", err,
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)
			}
		}()

		// После копирования photoReader будет пустым:
		buffer := new(bytes.Buffer)
		if _, err = io.Copy(buffer, photoReader); err != nil {
			return err
		}

		plant.Photo = buffer.Bytes()
		if err = useCases.UpdatePlant(*plant); err != nil {
			return err
		}

		group, err := useCases.GetGroup(plant.GroupID)
		if err != nil {
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
