package calendar

import (
	"errors"
	"time"

	"gopkg.in/telebot.v4"

	vd "github.com/go-ozzo/ozzo-validation"
)

type Option func(opts *options) error

// options represents a struct for passing optional
// properties for customizing a calendar keyboard.
type options struct {
	// The year that will be initially active in the calendar.
	// Default value - today's year
	initialYear int

	// The month that will be initially active in the calendar
	// Default value - today's month
	initialMonth time.Month

	// The range of displayed years
	// Default value - {1970, 292277026596} (time.Unix years range)
	yearsRange [2]int

	// The language of all designations.
	// If equals "ru" the designations would be Russian,
	// otherwise - English
	language string

	// Кнопка для возврата на предыдущий этап
	backButton *telebot.InlineButton
}

func (opts *options) validate() error {
	return vd.ValidateStruct(opts,
		vd.Field(&opts.yearsRange, vd.Required, vd.By(func(v any) error {
			rng, ok := v.([2]int)
			if !ok {
				return errors.New("invalid range")
			}

			if rng[0] < MinYearLimit || rng[1] > MaxYearLimit {
				return errors.New("yearsRange exceeds the acceptable limits of time.Unix")
			}

			return nil
		})),
		vd.Field(&opts.initialYear, vd.Required,
			vd.Min(opts.yearsRange[0]),
			vd.Max(opts.yearsRange[1]),
		),
		vd.Field(&opts.initialMonth, vd.Required, vd.Min(1), vd.Max(monthsPerYear)),
	)
}

func WithInitialYear(year int) Option {
	return func(options *options) error {
		if year == 0 {
			options.initialYear = time.Now().Year()

			return nil
		}

		options.initialYear = year

		return nil
	}
}

func WithInitialMonth(month time.Month) Option {
	return func(options *options) error {
		if month == 0 {
			options.initialMonth = time.Now().Month()

			return nil
		}

		options.initialMonth = month

		return nil
	}
}

func WithYearsRange(years [2]int) Option {
	return func(options *options) error {
		if years == [2]int{0, 0} {
			options.yearsRange = [2]int{MinYearLimit, MaxYearLimit}

			return nil
		}

		options.yearsRange = years

		return nil
	}
}

func WithLanguage(language string) Option {
	return func(options *options) error {
		options.language = language

		return nil
	}
}

func WithBackButton(button telebot.InlineButton) Option {
	return func(options *options) error {
		options.backButton = &button

		return nil
	}
}
