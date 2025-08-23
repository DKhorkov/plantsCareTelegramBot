package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

func AddPlantPhoto(
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

		// Для календаря используем context.Chat().ID. В случае с фото будет равен context.Sender().ID:
		temp, err := useCases.GetUserTemporary(int(context.Chat().ID))
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

		plant, err := useCases.AddPlantPhoto(int(context.Sender().ID), buffer.Bytes())
		if err != nil {
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
					buttons.ConfirmAddPlant,
				},
				{
					buttons.BackToAddPlantPhoto,
					buttons.Menu,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromReader(bytes.NewReader(plant.Photo)),
				Caption: fmt.Sprintf(
					texts.ConfirmAddPlant,
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

		return nil
	}
}

func ConfirmAddPlantCallback(
	_ *telebot.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
) telebot.HandlerFunc {
	return func(context telebot.Context) error {
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
			if context.Callback() == nil {
				logger.Warn(
					"Failed to send Response due to nil callback",
					"Message", context.Message(),
					"Sender", context.Sender(),
					"Chat", context.Chat(),
					"Callback", context.Callback(),
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)

				return errors.New("failed to send Response due to nil callback")
			}

			err = context.Respond(
				&telebot.CallbackResponse{
					CallbackID: context.Callback().ID,
					Text:       fmt.Sprintf(texts.PlantAlreadyExists, group.Title),
				},
			)
			if err != nil {
				logger.Error(
					"Failed to send Response",
					"Error", err,
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)

				return err
			}

			return nil
		}

		// Удаляем сообщение только если нет такой группы, чтобы отправить корректно CallbackResponse:
		if err = context.Delete(); err != nil {
			logger.Error(
				"Failed to delete message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		plant, err = useCases.CreatePlant(*plant)
		if err != nil {
			logger.Error(
				fmt.Sprintf("Failed to create Plant for user with telegramId=%d", context.Sender().ID),
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		// Только логгируем, поскольку не является критической логикой:
		if err = useCases.ResetTemporary(int(context.Sender().ID)); err != nil {
			logger.Error(
				fmt.Sprintf("Failed to reset Temporary for user with telegramId=%d", context.Sender().ID),
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.CreateAnotherPlant,
				},
				{
					buttons.Menu,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromReader(bytes.NewReader(plant.Photo)),
				Caption: fmt.Sprintf(
					texts.PlantCreated,
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

		return nil
	}
}
