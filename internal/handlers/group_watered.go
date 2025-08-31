package handlers

import (
	"errors"
	"strconv"
	"time"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
)

func GroupWateredCallback(_ interfaces.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
	return func(context telebot.Context) error {
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

		now := time.Now()
		wateredDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, group.NextWateringDate.Location())

		_, err = useCases.UpdateGroupLastWateringDate(groupID, wateredDate)
		if err != nil {
			return err
		}

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

		if _, err = context.Bot().EditReplyMarkup(context.Message(), &telebot.ReplyMarkup{}); err != nil {
			logger.Error(
				"Failed to delete ReplyMarkup",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		err = context.Respond(
			&telebot.CallbackResponse{
				CallbackID: context.Callback().ID,
				Text:       texts.GroupWatered,
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
}
