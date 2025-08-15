package calendar

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
)

const (
	monthsPerYear = 12
)

// Calendar represents the main object.
type Calendar struct {
	bot        *telebot.Bot
	logger     logging.Logger
	opt        *Options
	kb         [][]telebot.InlineButton
	currYear   int
	currMonth  time.Month
	backButton telebot.InlineButton
}

// NewCalendar builds and returns a Calendar.
func NewCalendar(bot *telebot.Bot, logger logging.Logger, opt Options) *Calendar {
	if opt.YearRange == [2]int{0, 0} {
		opt.YearRange = [2]int{MinYearLimit, MaxYearLimit}
	}

	if opt.InitialYear == 0 {
		opt.InitialYear = time.Now().Year()
	}

	if opt.InitialMonth == 0 {
		opt.InitialMonth = time.Now().Month()
	}

	err := opt.validate()
	if err != nil {
		panic(err)
	}

	return &Calendar{
		bot:       bot,
		logger:    logger,
		kb:        make([][]telebot.InlineButton, 0),
		opt:       &opt,
		currYear:  opt.InitialYear,
		currMonth: opt.InitialMonth,
	}
}

// Options represents a struct for passing optional
// properties for customizing a calendar keyboard.
type Options struct {
	// The year that will be initially active in the calendar.
	// Default value - today's year
	InitialYear int

	// The month that will be initially active in the calendar
	// Default value - today's month
	InitialMonth time.Month

	// The range of displayed years
	// Default value - {1970, 292277026596} (time.Unix years range)
	YearRange [2]int

	// The language of all designations.
	// If equals "ru" the designations would be Russian,
	// otherwise - English
	Language string
}

// GetKeyboard builds the calendar inline-keyboard.
func (cal *Calendar) GetKeyboard() [][]telebot.InlineButton {
	cal.clearKeyboard()

	cal.addMonthYearRow()
	cal.addWeekdaysRow()
	cal.addDaysRows()
	cal.addControlButtonsRow()
	cal.addBackAndMenuButtonsRow()

	return cal.kb
}

// SetBackButton sets back button for the calendar inline-keyboard logic.
func (cal *Calendar) SetBackButton(button telebot.InlineButton) {
	cal.backButton = button
}

// Clears the calendar's keyboard.
func (cal *Calendar) clearKeyboard() {
	cal.kb = make([][]telebot.InlineButton, 0)
}

// Builds a full row width button with a displayed month's name
// The button represents a list of all months when clicked.
func (cal *Calendar) addMonthYearRow() {
	var row []telebot.InlineButton

	btn := telebot.InlineButton{
		Unique: utils.GenUniqueParam("month_year_btn"),
		Text:   fmt.Sprintf("%s %v", cal.getMonthDisplayName(cal.currMonth), cal.currYear),
	}

	cal.bot.Handle(&btn, func(ctx telebot.Context) error {
		_, err := cal.bot.EditReplyMarkup(
			ctx.Message(),
			&telebot.ReplyMarkup{
				InlineKeyboard: cal.getMonthPickKeyboard(),
			},
		)
		if err != nil {
			cal.logger.Error("Failed to edit reply markup", "Error", err)
		}

		if err = ctx.Respond(); err != nil {
			cal.logger.Error("Failed to reply to message", "Error", err)
		}

		return nil
	})

	row = append(row, btn)
	cal.addRowToKeyboard(&row)
}

// Builds a keyboard with a list of months to pick.
func (cal *Calendar) getMonthPickKeyboard() [][]telebot.InlineButton {
	cal.clearKeyboard()

	var row []telebot.InlineButton

	// Generating a list of months
	for i := 1; i <= monthsPerYear; i++ {
		monthName := cal.getMonthDisplayName(time.Month(i))
		monthBtn := telebot.InlineButton{
			Unique: utils.GenUniqueParam("month_pick_" + strconv.Itoa(i)),
			Text:   monthName, Data: strconv.Itoa(i),
		}

		cal.bot.Handle(&monthBtn, func(ctx telebot.Context) error {
			monthNum, err := strconv.Atoi(ctx.Data())
			if err != nil {
				log.Fatal(err)
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
				cal.logger.Error("Failed to edit reply markup", "Error", err)
			}

			if err = ctx.Respond(); err != nil {
				cal.logger.Error("Failed to reply to message", "Error", err)
			}

			return nil
		})

		row = append(row, monthBtn)

		// Arranging the months in 2 columns
		if i%2 == 0 {
			cal.addRowToKeyboard(&row)
			row = []telebot.InlineButton{} // empty row
		}
	}

	cal.addBackAndMenuButtonsRow()

	return cal.kb
}

// Builds a row of non-clickable buttons
// that display weekdays names.
func (cal *Calendar) addWeekdaysRow() {
	var row []telebot.InlineButton

	for i, wd := range cal.getWeekdaysDisplayArray() {
		btn := telebot.InlineButton{Unique: utils.GenUniqueParam("weekday_" + strconv.Itoa(i)), Text: wd}
		cal.bot.Handle(&btn, cal.ignoreQuery())
		row = append(row, btn)
	}

	cal.addRowToKeyboard(&row)
}

// Builds a table of clickable cells (buttons) - active month's days.
func (cal *Calendar) addDaysRows() {
	beginningOfMonth := time.Date(cal.currYear, cal.currMonth, 1, 0, 0, 0, 0, time.UTC)
	amountOfDaysInMonth := beginningOfMonth.AddDate(0, 1, -1).Day()

	var row []telebot.InlineButton

	// Calculating the number of empty buttons that need to be inserted forward
	weekdayNumber := int(beginningOfMonth.Weekday())
	if weekdayNumber == 0 && cal.opt.Language == RussianLangAbbr { // russian Sunday exception
		weekdayNumber = 7
	}

	// The difference between English and Russian weekdays order
	// en: Sunday (0), Monday (1), Tuesday (3), ...
	// ru: Monday (1), Tuesday (2), ..., Sunday (7)
	if cal.opt.Language != RussianLangAbbr {
		weekdayNumber++
	}

	// Inserting empty buttons forward
	for i := 1; i < weekdayNumber; i++ {
		cal.addEmptyCell(&row)
	}

	// Inserting month's days' buttons
	for i := 1; i <= amountOfDaysInMonth; i++ {
		dayText := strconv.Itoa(i)
		cell := telebot.InlineButton{
			Unique: utils.GenUniqueParam("day_" + strconv.Itoa(i)),
			Text:   dayText, Data: dayText,
		}

		cal.bot.Handle(&cell, func(ctx telebot.Context) error {
			dayInt, err := strconv.Atoi(ctx.Data())
			if err != nil {
				return err
			}

			ctx.Message().Payload = cal.genDateStrFromDay(dayInt)

			upd := telebot.Update{Message: ctx.Message()}
			cal.bot.ProcessUpdate(upd)

			if err = ctx.Respond(); err != nil {
				cal.logger.Error("Failed to reply to message", "Error", err)
			}

			return nil
		})

		row = append(row, cell)

		if len(row)%AmountOfDaysInWeek == 0 {
			cal.addRowToKeyboard(&row)
			row = []telebot.InlineButton{} // empty row
		}
	}

	// Inseting empty buttons at the end
	if len(row) > 0 {
		for i := len(row); i < AmountOfDaysInWeek; i++ {
			cal.addEmptyCell(&row)
		}

		cal.addRowToKeyboard(&row)
	}
}

// Builds a row of  control buttons for swiping the calendar.
func (cal *Calendar) addControlButtonsRow() {
	var row []telebot.InlineButton

	prev := telebot.InlineButton{Unique: utils.GenUniqueParam("prev_month"), Text: "＜"}

	// Hide "prev" button if it rests on the range
	if cal.currYear <= cal.opt.YearRange[0] && cal.currMonth == 1 {
		prev.Text = ""
	} else {
		cal.bot.Handle(&prev, func(ctx telebot.Context) error {
			// Additional protection against entering the years ranges
			if cal.currMonth > 1 {
				cal.currMonth--
			} else {
				cal.currMonth = monthsPerYear
				if cal.currYear > cal.opt.YearRange[0] {
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
				cal.logger.Error("Failed to edit reply markup", "Error", err)
			}

			if err = ctx.Respond(); err != nil {
				cal.logger.Error("Failed to reply to message", "Error", err)
			}

			return nil
		})
	}

	next := telebot.InlineButton{Unique: utils.GenUniqueParam("next_month"), Text: "＞"}

	// Hide "next" button if it rests on the range
	if cal.currYear >= cal.opt.YearRange[1] && cal.currMonth == monthsPerYear {
		next.Text = ""
	} else {
		cal.bot.Handle(&next, func(ctx telebot.Context) error {
			// Additional protection against entering the years ranges
			if cal.currMonth < monthsPerYear {
				cal.currMonth++
			} else {
				if cal.currYear < cal.opt.YearRange[1] {
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
				cal.logger.Error("Failed to edit reply markup", "Error", err)
			}

			if err = ctx.Respond(); err != nil {
				cal.logger.Error("Failed to reply to message", "Error", err)
			}

			return nil
		})
	}

	row = append(row, prev, next)
	cal.addRowToKeyboard(&row)
}

// Builds a row of back and menu buttons.
func (cal *Calendar) addBackAndMenuButtonsRow() {
	var row []telebot.InlineButton

	row = append(row, cal.backButton, buttons.MenuButton)
	cal.addRowToKeyboard(&row)
}

// Returns a formatted date string from the selected date.
func (cal *Calendar) genDateStrFromDay(day int) string {
	return time.Date(cal.currYear, cal.currMonth, day,
		0, 0, 0, 0, time.UTC).Format("02.01.2006")
}

// Utility function for passing a row to the calendar's keyboard.
func (cal *Calendar) addRowToKeyboard(row *[]telebot.InlineButton) {
	cal.kb = append(cal.kb, *row)
}

// Inserts an empty button that doesn't process queries
// into the keyboard row.
func (cal *Calendar) addEmptyCell(row *[]telebot.InlineButton) {
	cell := telebot.InlineButton{Unique: utils.GenUniqueParam("empty_cell"), Text: " "}
	cal.bot.Handle(&cell, cal.ignoreQuery())
	*row = append(*row, cell)
}

// Query stub.
func (cal *Calendar) ignoreQuery() func(ctx telebot.Context) error {
	return func(ctx telebot.Context) error {
		err := ctx.Respond()
		if err != nil {
			cal.logger.Error("Failed to reply to message", "Error", err)
		}

		return nil
	}
}
