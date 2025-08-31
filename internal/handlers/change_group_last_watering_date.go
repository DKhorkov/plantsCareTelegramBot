package handlers

import (
	"fmt"
	"time"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
)

func ChangeGroupLastWateringDate(
	_ interfaces.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		lastWateringDate, err := time.Parse(dateFormat, context.Data())
		if err != nil {
			logger.Error(
				"Failed to parse last watering date",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		// Дата последнего полива не может быть позже текущего дня:
		if time.Now().Before(lastWateringDate) {
			// Нет context.Callback() для обычного сообщения, поэтому отправляем ответ текстом:
			if err = context.Send(texts.LastWateringDateInFuture); err != nil {
				logger.Error(
					"Failed to send message",
					"Error", err,
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)

				return err
			}

			return nil
		}

		// Удаляем только если дата выбрана корректно (раньше текущего дня):
		if err = context.Delete(); err != nil {
			logger.Error(
				"Failed to delete message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		// Для календаря используем context.Chat().ID:
		temp, err := useCases.GetUserTemporary(int(context.Chat().ID))
		if err != nil {
			return err
		}

		group, err := temp.GetGroup()
		if err != nil {
			logger.Error(
				"Failed to get Group from Temporary",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		group, err = useCases.UpdateGroupLastWateringDate(group.ID, lastWateringDate)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.ManageGroupChangeTitle,
				},
				{
					buttons.ManageGroupChangeDescription,
				},
				{
					buttons.ManageGroupChangeLastWateringDate,
				},
				{
					buttons.ManageGroupChangeWateringInterval,
				},
				{
					buttons.BackToManageGroupAction,
					buttons.Menu,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.ManageGroupChangeImage),
				Caption: fmt.Sprintf(
					texts.ManageGroupChange,
					group.Title,
					group.Description,
					group.LastWateringDate.Format(dateFormat),
					utils.GetWateringInterval(group.WateringInterval),
					group.NextWateringDate.Format(dateFormat),
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
		if err = useCases.SetTemporaryStep(int(context.Chat().ID), steps.ManageGroupChange); err != nil {
			return err
		}

		return nil
	}
}
