package handlers

import (
	"bytes"
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

func ManagePlantRemovalCallback(
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

		plant, err = useCases.GetPlant(plant.ID)
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
					buttons.ConfirmPlantRemoval,
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
					texts.ManagePlantRemoval,
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
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.ManagePlantRemoval); err != nil {
			return err
		}

		return nil
	}
}

func BackToManagePlantActionCallback(
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

		plant, err = useCases.GetPlant(plant.ID)
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
					buttons.ManagePlantChange,
				},
				{
					buttons.ManagePlantRemoval,
				},
				{
					buttons.BackToManagePlant,
					buttons.Menu,
				},
			},
		}

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromReader(bytes.NewReader(plant.Photo)),
				Caption: fmt.Sprintf(
					texts.ManagePlantAction,
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
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.ManagePlantAction); err != nil {
			return err
		}

		return nil
	}
}

func ConfirmPlantRemovalCallback(
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

		if err = useCases.DeletePlant(plant.ID); err != nil {
			return err
		}

		err = context.Respond(
			&telebot.CallbackResponse{
				CallbackID: context.Callback().ID,
				Text:       texts.PlantDeleted,
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

		if err = context.Delete(); err != nil {
			logger.Error(
				"Failed to delete message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		user, err := useCases.GetUserByTelegramID(int(context.Sender().ID))
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.CreateGroup,
				},
			},
		}

		groupsCount, err := useCases.CountUserGroups(user.ID)
		if err != nil {
			return err
		}

		if groupsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.CreatePlant})
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.ManageGroups})
		}

		plantsCount, err := useCases.CountUserPlants(user.ID)
		if err != nil {
			return err
		}

		if plantsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.ManagePlants})
		}

		err = context.Send(
			&telebot.Photo{
				File:    telebot.FromDisk(paths.StartImage),
				Caption: texts.OnStart,
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

		// TODO при проблемах логики следует сделать в рамках транзакции.
		if err = useCases.ResetTemporary(int(context.Sender().ID)); err != nil {
			return err
		}

		return nil
	}
}
