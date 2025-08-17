package handlers

import (
	"errors"
	"fmt"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
	"gopkg.in/telebot.v4"
	"strconv"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

func AddGroupWateringInterval(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		wateringInterval, err := strconv.Atoi(context.Data())
		if err != nil {
			logger.Error("Failed to parse watering interval", "Error", err)

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
					buttons.ConfirmAddGroupButton,
				},
				{
					buttons.BackToAddGroupWateringIntervalButton,
					buttons.MenuButton,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.AddGroupConfirmDateImagePath),
				Caption: fmt.Sprintf(
					texts.AddGroupConfirmText,
					group.Title,
					group.Description,
					group.LastWateringDate.Format(dateFormat),
					getWateringIntervalText(group.WateringInterval),
					group.NextWateringDate.Format(dateFormat),
				),
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

func BackToAddGroupWateringIntervalCallback(
	bot *telebot.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		temp, err := useCases.GetUserTemporary(int(context.Sender().ID))
		if err != nil {
			return err
		}

		// Получаем группу для корректного отображения данных прошлых этапов::
		group, err := temp.GetGroup()
		if err != nil {
			logger.Error("Failed to get Group from Temporary", "Error", err)

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
				Text:   getWateringIntervalText(value),
				Data:   strconv.Itoa(value),
			}

			bot.Handle(&btn, AddGroupWateringInterval(bot, useCases, logger))

			row = append(row, btn)
			if len(row) == buttonsPerRaw {
				menu.InlineKeyboard = append(menu.InlineKeyboard, row)
				row = []telebot.InlineButton{}
			}
		}

		menu.InlineKeyboard = append(
			menu.InlineKeyboard,
			[]telebot.InlineButton{
				buttons.BackToAddGroupLastWateringDateButton,
				buttons.MenuButton,
			},
		)

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.AddGroupWateringIntervalImagePath),
				Caption: fmt.Sprintf(
					texts.AddGroupWateringIntervalText,
					group.Title,
					group.Description,
					group.LastWateringDate.Format(dateFormat),
					group.Title,
				),
			},
			menu,
		)
		if err != nil {
			logger.Error("Failed to send message", "Error", err)

			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddGroupWateringIntervalStep); err != nil {
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
			logger.Error("Failed to get Group from Temporary", "Error", err)

			return err
		}

		groupExists, err := useCases.GroupExists(*group)
		if err != nil {
			logger.Error(
				fmt.Sprintf("Failed to check Group existence for user with telegramId=%d", context.Sender().ID),
				"Error",
				err,
			)

			return err
		}

		if groupExists {
			if context.Callback() == nil {
				logger.Warn(
					"Failed to send Response due to nil callback",
					"Message",
					context.Message(),
					"Sender",
					context.Sender(),
					"Chat",
					context.Chat(),
					"Callback",
					context.Callback(),
				)

				return errors.New("failed to send Response due to nil callback")
			}

			err = context.Respond(
				&telebot.CallbackResponse{
					CallbackID: context.Callback().ID,
					Text:       fmt.Sprintf(texts.GroupAlreadyExists),
				},
			)
			if err != nil {
				logger.Error("Failed to send Response", "Error", err)

				return err
			}

			return nil
		}

		// Удаляем сообщение только если нет такой группы, чтобы отправить корректно CallbackResponse:
		if err = context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		group, err = useCases.CreateGroup(*group)
		if err != nil {
			logger.Error(
				fmt.Sprintf("Failed to create Group for user with telegramId=%d", context.Sender().ID),
				"Error",
				err,
			)

			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.AddPlantButton,
				},
				{
					buttons.MenuButton,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.GroupCreatedImagePath),
				Caption: fmt.Sprintf(
					texts.GroupCreatedText,
					group.Title,
					group.Description,
					group.LastWateringDate.Format(dateFormat),
					getWateringIntervalText(group.WateringInterval),
					group.NextWateringDate.Format(dateFormat),
				),
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
