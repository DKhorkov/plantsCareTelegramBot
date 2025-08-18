package handlers

import (
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

const (
	plantGroupButtonsPerRaw = 1
)

func AddPlantDescription(bot *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
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

		plant, err := useCases.AddPlantDescription(int(context.Sender().ID), context.Message().Text)
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
				Unique: utils.GenUniqueParam("plant_group"),
				Text:   group.Title,
				Data:   strconv.Itoa(group.ID),
			}

			bot.Handle(&btn, AddPlantGroupCallback(bot, useCases, logger))

			row = append(row, btn)
			if len(row) == plantGroupButtonsPerRaw {
				menu.InlineKeyboard = append(menu.InlineKeyboard, row)
				row = []telebot.InlineButton{}
			}
		}

		menu.InlineKeyboard = append(
			menu.InlineKeyboard,
			[]telebot.InlineButton{
				buttons.BackToAddPlantDescriptionButton,
				buttons.MenuButton,
			},
		)

		err = context.Send(
			&telebot.Photo{
				File:    telebot.FromDisk(paths.AddPlantGroupImagePath),
				Caption: fmt.Sprintf(texts.AddPlantGroupText, plant.Title, plant.Description, plant.Title),
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

func BackToAddPlantDescriptionCallback(
	_ *telebot.Bot,
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

		// Получаем растение для корректного отображения данных прошлых этапов:
		plant, err := temp.GetPlant()
		if err != nil {
			logger.Error("Failed to get Plant from Temporary", "Error", err)

			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.SkipPlantDescriptionButton,
				},
				{
					buttons.BackToAddPlantTitleButton,
					buttons.MenuButton,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File:    telebot.FromDisk(paths.AddPlantDescriptionImagePath),
				Caption: fmt.Sprintf(texts.AddPlantDescriptionText, plant.Title, plant.Title),
			},
			menu,
		)
		if err != nil {
			logger.Error("Failed to send message", "Error", err)

			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddPlantDescriptionStep); err != nil {
			return err
		}

		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), &msg.ID); err != nil {
			return err
		}

		return nil
	}
}
