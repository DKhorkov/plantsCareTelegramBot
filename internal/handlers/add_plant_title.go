package handlers

import (
	"fmt"
	"strconv"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

const (
	plantTitleMaxLength = 50
)

func AddPlantTitle(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
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

		plant, err := useCases.AddPlantTitle(int(context.Sender().ID), context.Message().Text)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.SkipPlantDescription,
				},
				{
					buttons.BackToAddPlantTitle,
					buttons.Menu,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File:    telebot.FromDisk(paths.AddPlantDescriptionImage),
				Caption: fmt.Sprintf(texts.AddPlantDescription, plant.Title, plant.Title),
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

		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), &msg.ID); err != nil {
			return err
		}

		return nil
	}
}

func SkipPlantDescriptionCallback(
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

		plant, err := useCases.AddPlantDescription(int(context.Sender().ID), "➖")
		if err != nil {
			return err
		}

		user, err := useCases.GetUserByTelegramID(int(context.Sender().ID))
		if err != nil {
			return err
		}

		groups, err := useCases.GetUserGroups(user.ID)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{},
		}

		var row []telebot.InlineButton

		for _, group := range groups {
			btn := telebot.InlineButton{
				Unique: buttons.AddPlantGroup.Unique,
				Text:   group.Title,
				Data:   strconv.Itoa(group.ID),
			}

			row = append(row, btn)
			if len(row) == plantGroupButtonsPerRaw {
				menu.InlineKeyboard = append(menu.InlineKeyboard, row)
				row = []telebot.InlineButton{}
			}
		}

		menu.InlineKeyboard = append(
			menu.InlineKeyboard,
			[]telebot.InlineButton{
				buttons.BackToAddPlantDescription,
				buttons.Menu,
			},
		)

		err = context.Send(
			&telebot.Photo{
				File:    telebot.FromDisk(paths.AddPlantGroupImage),
				Caption: fmt.Sprintf(texts.AddPlantGroup, plant.Title, plant.Description, plant.Title),
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
