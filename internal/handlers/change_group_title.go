package handlers

import (
	"errors"
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	customerrors "github.com/DKhorkov/plantsCareTelegramBot/internal/errors"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
)

func ChangeGroupTitle(
	_ interfaces.Bot,
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

		if len(context.Message().Text) > groupTitleMaxLength {
			// Нет context.Callback() для обычного сообщения, поэтому отправляем ответ текстом:
			if err := context.Send(fmt.Sprintf(texts.GroupTitleTooLong, groupTitleMaxLength)); err != nil {
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

		group, err := temp.GetGroup()
		if err != nil {
			logger.Error(
				"Failed to get Group from Temporary",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		group, err = useCases.UpdateGroupTitle(group.ID, context.Message().Text)

		switch {
		case errors.Is(err, customerrors.ErrGroupAlreadyExists):
			// Нет context.Callback() для обычного сообщения, поэтому отправляем ответ текстом:
			if err = context.Send(texts.GroupAlreadyExists); err != nil {
				logger.Error(
					"Failed to send message",
					"Error", err,
					"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
				)

				return err
			}

			return nil
		case err != nil:
			return err
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
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.ManageGroupChange); err != nil {
			return err
		}

		return nil
	}
}
