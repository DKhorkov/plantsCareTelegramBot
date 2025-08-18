package handlers

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/calendar"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
)

const (
	dateFormat                         = "02.01.2006"
	groupWateringIntervalButtonsPerRaw = 2
)

var wateringIntervals = []int{1, 2, 3, 4, 5, 6, 7, 10, 14, 18, 21, 30}

func AddGroupLastWateringDate(
	bot *telebot.Bot,
	useCases interfaces.UseCases,
	logger logging.Logger,
) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		lastWateringDate, err := time.Parse(dateFormat, context.Data())
		if err != nil {
			logger.Error("Failed to parse last watering date", "Error", err)

			return err
		}

		// Для календаря используем context.Chat().ID
		group, err := useCases.AddGroupLastWateringDate(int(context.Chat().ID), lastWateringDate)
		if err != nil {
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

		return nil
	}
}

func BackToAddGroupLastWateringDateCallback(
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

		// Получаем группу для корректного отображения данных прошлых этапов:
		group, err := temp.GetGroup()
		if err != nil {
			logger.Error("Failed to get Group from Temporary", "Error", err)

			return err
		}

		c := calendar.NewCalendar(bot, logger, calendar.Options{Language: "ru"})
		c.SetBackButton(buttons.BackToAddGroupDescriptionButton)
		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: c.GetKeyboard(),
		}

		err = context.Send(
			&telebot.Photo{
				File:    telebot.FromDisk(paths.AddGroupLastWateringDateImagePath),
				Caption: fmt.Sprintf(texts.AddGroupLastWateringDateText, group.Title, group.Description, group.Title),
			},
			menu,
		)
		if err != nil {
			logger.Error("Failed to send message", "Error", err)

			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddGroupLastWateringDateStep); err != nil {
			return err
		}

		return nil
	}
}

func getWateringIntervalText(wateringInterval int) string {
	switch {
	case wateringInterval%10 == 1:
		return fmt.Sprintf("%d день", wateringInterval)
	case slices.Contains([]int{2, 3, 4}, wateringInterval):
		return fmt.Sprintf("%d дня", wateringInterval)
	default:
		return fmt.Sprintf("%d дней", wateringInterval)
	}
}
