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

const (
	plantsPerGroupLimit = 50
)

func AddPlantGroupCallback(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		groupID, err := strconv.Atoi(context.Data())
		if err != nil {
			logger.Error("Failed to parse groupID", "Error", err)

			return err
		}

		groupPlantsCount, err := useCases.CountGroupPlants(groupID)
		if err != nil {
			return err
		}

		if groupPlantsCount >= plantsPerGroupLimit {
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
					Text:       fmt.Sprintf(texts.PlantsPerGroupLimit, plantsPerGroupLimit),
				},
			)
			if err != nil {
				logger.Error("Failed to send Response", "Error", err)

				return err
			}

			return nil
		}

		if err = context.Delete(); err != nil {
			logger.Error("Failed to delete message", "Error", err)

			return err
		}

		plant, err := useCases.AddPlantGroup(int(context.Sender().ID), groupID)
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
					buttons.AcceptAddPlantPhoto,
				},
				{
					buttons.RejectAddPlantPhoto,
				},
				{
					buttons.BackToAddPlantGroup,
					buttons.Menu,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.AddPlantPhotoQuestionImage),
				Caption: fmt.Sprintf(
					texts.AddPlantPhotoQuestion,
					plant.Title,
					plant.Description,
					group.Title,
					group.Title,
					plant.Title,
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

func BackToAddPlantGroupCallback(
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

		// Получаем растение для корректного отображения данных прошлых этапов:
		plant, err := temp.GetPlant()
		if err != nil {
			logger.Error("Failed to get Plant from Temporary", "Error", err)

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
			logger.Error("Failed to send message", "Error", err)

			return err
		}

		// TODO при проблемах логики следует сделать в рамках транзакции
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddPlantGroup); err != nil {
			return err
		}

		return nil
	}
}
