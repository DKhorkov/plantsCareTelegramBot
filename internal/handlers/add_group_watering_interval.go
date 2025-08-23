package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
)

func AddGroupWateringIntervalCallback(
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

		wateringInterval, err := strconv.Atoi(context.Data())
		if err != nil {
			logger.Error(
				"Failed to parse watering interval",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		group, err := useCases.AddGroupWateringInterval(int(context.Sender().ID), wateringInterval)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.ConfirmAddGroup,
				},
				{
					buttons.BackToAddGroupWateringInterval,
					buttons.Menu,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.AddGroupConfirmImage),
				Caption: fmt.Sprintf(
					texts.ConfirmAddGroup,
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

		return nil
	}
}

func BackToAddGroupWateringIntervalCallback(
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

		// Получаем группу для корректного отображения данных прошлых этапов:
		group, err := temp.GetGroup()
		if err != nil {
			logger.Error(
				"Failed to get Group from Temporary",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{},
		}

		var row []telebot.InlineButton

		for _, value := range wateringIntervals {
			btn := telebot.InlineButton{
				Unique: utils.GenUniqueParam("watering_interval"),
				Text:   utils.GetWateringInterval(value),
				Data:   strconv.Itoa(value),
			}

			bot.Handle(&btn, AddGroupWateringIntervalCallback(bot, useCases, logger))

			row = append(row, btn)
			if len(row) == groupWateringIntervalButtonsPerRaw {
				menu.InlineKeyboard = append(menu.InlineKeyboard, row)
				row = []telebot.InlineButton{}
			}
		}

		menu.InlineKeyboard = append(
			menu.InlineKeyboard,
			[]telebot.InlineButton{
				buttons.BackToAddGroupLastWateringDate,
				buttons.Menu,
			},
		)

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.AddGroupWateringIntervalImage),
				Caption: fmt.Sprintf(
					texts.AddGroupWateringInterval,
					group.Title,
					group.Description,
					group.LastWateringDate.Format(dateFormat),
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
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddGroupWateringInterval); err != nil {
			return err
		}

		return nil
	}
}

func ConfirmAddGroupCallback(
	_ *telebot.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
) telebot.HandlerFunc {
	return func(context telebot.Context) error {
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

		groupExists, err := useCases.GroupExists(*group)
		if err != nil {
			logger.Error(
				fmt.Sprintf("Failed to check Group existence for user with telegramId=%d", context.Sender().ID),
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		if groupExists {
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
					Text:       texts.GroupAlreadyExists,
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

		group, err = useCases.CreateGroup(*group)
		if err != nil {
			logger.Error(
				fmt.Sprintf("Failed to create Group for user with telegramId=%d", context.Sender().ID),
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
					buttons.CreatePlant,
				},
				{
					buttons.Menu,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.GroupCreatedImage),
				Caption: fmt.Sprintf(
					texts.GroupCreated,
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

		return nil
	}
}
