package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/paths"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

const (
	groupsPerUserLimit             = 5
	managePlantsGroupButtonsPerRaw = 1
	manageGroupButtonsPerRaw       = 1
)

func Start(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error(
				"Failed to delete /start message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		userID, err := useCases.SaveUser(
			entities.User{
				TelegramID: int(context.Sender().ID),
				Username:   context.Sender().Username,
				Firstname:  context.Sender().FirstName,
				Lastname:   context.Sender().LastName,
				IsBot:      context.Sender().IsBot,
			},
		)
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

		groupsCount, err := useCases.CountUserGroups(userID)
		if err != nil {
			return err
		}

		if groupsCount > 0 {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.CreatePlant})
			menu.InlineKeyboard = append(menu.InlineKeyboard, []telebot.InlineButton{buttons.ManageGroups})
		}

		plantsCount, err := useCases.CountUserPlants(userID)
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
		// TODO Тут повторяем вне юзкейсов, чтобы работало даже вне повторной регистрации.
		if err = useCases.ResetTemporary(int(context.Sender().ID)); err != nil {
			return err
		}

		return nil
	}
}

func AddGroupCallback(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		user, err := useCases.GetUserByTelegramID(int(context.Sender().ID))
		if err != nil {
			return err
		}

		groupsCount, err := useCases.CountUserGroups(user.ID)
		if err != nil {
			return err
		}

		if groupsCount >= groupsPerUserLimit {
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
					Text:       fmt.Sprintf(texts.GroupsPerUserLimit, groupsPerUserLimit),
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

		if err = context.Delete(); err != nil {
			logger.Error(
				"Failed to delete message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.BackToStart,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File:    telebot.FromDisk(paths.AddGroupTitleImage),
				Caption: texts.AddGroupTitle,
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
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddGroupTitle); err != nil {
			return err
		}

		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), &msg.ID); err != nil {
			return err
		}

		return nil
	}
}

func AddPlantCallback(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
			logger.Error(
				"Failed to delete message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{
				{
					buttons.BackToStart,
				},
			},
		}

		// Получаем бота, чтобы при отправке получить messageID для дальнейшего удаления:
		msg, err := context.Bot().Send(
			context.Chat(),
			&telebot.Photo{
				File:    telebot.FromDisk(paths.AddPlantTitleImage),
				Caption: texts.AddPlantTitle,
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
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.AddPlantTitle); err != nil {
			return err
		}

		if err = useCases.SetTemporaryMessage(int(context.Sender().ID), &msg.ID); err != nil {
			return err
		}

		return nil
	}
}

func ManagePlantsCallback(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
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

		groups, err := useCases.GetUserGroups(user.ID)
		if err != nil {
			return err
		}

		// Показываем только те сценарии, где есть растения:
		var groupsWithPlants []entities.Group

		for _, group := range groups {
			plants, err := useCases.GetGroupPlants(group.ID)
			if err != nil {
				return err
			}

			if len(plants) > 0 {
				groupsWithPlants = append(groupsWithPlants, group)
			}
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{},
		}

		var row []telebot.InlineButton

		for _, group := range groupsWithPlants {
			btn := telebot.InlineButton{
				Unique: buttons.ManagePlantsGroup.Unique,
				Text:   group.Title,
				Data:   strconv.Itoa(group.ID),
			}

			row = append(row, btn)
			if len(row) == managePlantsGroupButtonsPerRaw {
				menu.InlineKeyboard = append(menu.InlineKeyboard, row)
				row = []telebot.InlineButton{}
			}
		}

		menu.InlineKeyboard = append(
			menu.InlineKeyboard,
			[]telebot.InlineButton{
				buttons.BackToStart,
			},
		)

		err = context.Send(
			&telebot.Photo{
				File:    telebot.FromDisk(paths.ManagePlantsChooseGroupImage),
				Caption: texts.ManagePlantsChooseGroup,
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
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.ManagePlantsChooseGroup); err != nil {
			return err
		}

		return nil
	}
}

func ManageGroupsCallback(_ *telebot.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if err := context.Delete(); err != nil {
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
				Unique: buttons.ManageGroup.Unique,
				Text:   group.Title,
				Data:   strconv.Itoa(group.ID),
			}

			row = append(row, btn)
			if len(row) == manageGroupButtonsPerRaw {
				menu.InlineKeyboard = append(menu.InlineKeyboard, row)
				row = []telebot.InlineButton{}
			}
		}

		menu.InlineKeyboard = append(
			menu.InlineKeyboard,
			[]telebot.InlineButton{
				buttons.BackToStart,
			},
		)

		err = context.Send(
			&telebot.Photo{
				File:    telebot.FromDisk(paths.ManageGroupImage),
				Caption: texts.ManageGroup,
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
		if err = useCases.SetTemporaryStep(int(context.Sender().ID), steps.ManageGroup); err != nil {
			return err
		}

		return nil
	}
}
