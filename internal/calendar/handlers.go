package calendar

import (
	"strconv"
	"time"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"
)

// Хэндлеры календаря остаются с ним в одном пакете из-за неэкспортируемых полей ради инкапсуляции.

type Handler func(cal *Calendar) telebot.HandlerFunc

func MonthsPerYearCallback(cal *Calendar) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		_, err := cal.bot.EditReplyMarkup(
			ctx.Message(),
			&telebot.ReplyMarkup{
				InlineKeyboard: cal.getMonthPickKeyboard(),
			},
		)
		if err != nil {
			cal.logger.Error(
				"Failed to edit reply markup",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		if err = ctx.Respond(); err != nil {
			cal.logger.Error(
				"Failed to reply to message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		return nil
	}
}

func PickedMonthCallback(cal *Calendar) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		monthNum, err := strconv.Atoi(ctx.Data())
		if err != nil {
			cal.logger.Error(
				"Failed to get month number",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)

			return err
		}

		cal.currMonth = time.Month(monthNum)

		// Show the calendar keyboard with the active selected month back
		_, err = cal.bot.EditReplyMarkup(
			ctx.Message(),
			&telebot.ReplyMarkup{
				InlineKeyboard: cal.GetKeyboard(),
			},
		)
		if err != nil {
			cal.logger.Error(
				"Failed to edit reply markup",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		if err = ctx.Respond(); err != nil {
			cal.logger.Error(
				"Failed to reply to message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		return nil
	}
}

func IgnoreQueryCallback(cal *Calendar) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		err := ctx.Respond()
		if err != nil {
			cal.logger.Error(
				"Failed to reply to message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		return nil
	}
}

func SelectedDayCallback(cal *Calendar) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		dayInt, err := strconv.Atoi(ctx.Data())
		if err != nil {
			return err
		}

		ctx.Message().Payload = cal.genDateStrFromDay(dayInt)

		upd := telebot.Update{Message: ctx.Message()}
		cal.bot.ProcessUpdate(upd)

		if err = ctx.Respond(); err != nil {
			cal.logger.Error(
				"Failed to reply to message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		return nil
	}
}

func PreviousMonthCallback(cal *Calendar) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		// Additional protection against entering the years ranges
		if cal.currMonth > 1 {
			cal.currMonth--
		} else {
			cal.currMonth = monthsPerYear
			if cal.currYear > cal.yearsRange[0] {
				cal.currYear--
			}
		}

		_, err := cal.bot.EditReplyMarkup(
			ctx.Message(),
			&telebot.ReplyMarkup{
				InlineKeyboard: cal.GetKeyboard(),
			},
		)
		if err != nil {
			cal.logger.Error(
				"Failed to edit reply markup",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		if err = ctx.Respond(); err != nil {
			cal.logger.Error(
				"Failed to reply to message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		return nil
	}
}

func NextMonthCallback(cal *Calendar) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		// Additional protection against entering the years ranges
		if cal.currMonth < monthsPerYear {
			cal.currMonth++
		} else {
			if cal.currYear < cal.yearsRange[1] {
				cal.currYear++
			}

			cal.currMonth = 1
		}

		_, err := cal.bot.EditReplyMarkup(
			ctx.Message(),
			&telebot.ReplyMarkup{
				InlineKeyboard: cal.GetKeyboard(),
			},
		)
		if err != nil {
			cal.logger.Error(
				"Failed to edit reply markup",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		if err = ctx.Respond(); err != nil {
			cal.logger.Error(
				"Failed to reply to message",
				"Error", err,
				"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
			)
		}

		return nil
	}
}
