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
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
)

func ManageGroupCallback(
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

		groupID, err := strconv.Atoi(context.Data())
		if err != nil {
			logger.Error(
				"Failed to parse groupID",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		group, err := useCases.GetGroup(groupID)
		if err != nil {
			return err
		}

		plants, err := useCases.GetGroupPlants(group.ID)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{
			ResizeKeyboard: true,
			InlineKeyboard: [][]telebot.InlineButton{},
		}

		if len(plants) > 0 {
			menu.InlineKeyboard = append(
				menu.InlineKeyboard,
				[]telebot.InlineButton{
					buttons.ManageGroupSeePlants,
				},
			)
		}

		menu.InlineKeyboard = append(
			menu.InlineKeyboard,
			[]telebot.InlineButton{
				buttons.ManageGroupChange,
			},
			[]telebot.InlineButton{
				buttons.ManageGroupRemoval,
			},
			[]telebot.InlineButton{
				buttons.BackToManageGroup,
				buttons.Menu,
			},
		)

		err = context.Send(
			&telebot.Photo{
				File: telebot.FromDisk(paths.ManageGroupActionImage),
				Caption: fmt.Sprintf(
					texts.ManageGroupAction,
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

		if err = useCases.ManageGroup(int(context.Sender().ID), groupID); err != nil {
			return err
		}

		return nil
	}
}
