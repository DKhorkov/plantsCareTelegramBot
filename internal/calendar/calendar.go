package calendar

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DKhorkov/libs/logging"
	"gopkg.in/telebot.v4"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
)

// Calendar represents the main object.
type Calendar struct {
	bot        interfaces.Bot
	logger     logging.Logger
	kb         [][]telebot.InlineButton
	currYear   int
	currMonth  time.Month
	yearsRange [2]int
	language   string
	buttons    map[string]*telebot.InlineButton
}

// NewCalendar builds and returns a Calendar.
func NewCalendar(bot interfaces.Bot, logger logging.Logger, opts ...Option) (*Calendar, error) {
	now := time.Now()
	calendarOptions := options{
		initialYear:  now.Year(),
		initialMonth: now.Month(),
		yearsRange:   [2]int{now.Year(), now.Year()},
		language:     RussianLangAbbr,
	}

	for _, opt := range opts {
		err := opt(&calendarOptions)
		if err != nil {
			return nil, err
		}
	}

	err := calendarOptions.validate()
	if err != nil {
		return nil, err
	}

	unique := utils.GenUniqueParam("") // Уникальное значение для кнопок экземпляра каждого календаря
	btns := map[string]*telebot.InlineButton{
		monthsPerYearButton: {Unique: buttonsPrefix + unique},
		pickedMonthButton:   {Unique: pickedMonthButton + unique},
		ignoreQueryButton:   {Unique: ignoreQueryButton + unique},
		selectedDayButton:   {Unique: selectedDayButton + unique},
		previousMonthButton: {Unique: previousMonthButton + unique},
		nextMonthButton:     {Unique: nextMonthButton + unique},
	}

	if calendarOptions.backButton != nil {
		btns[backButton] = calendarOptions.backButton
	}

	cal := &Calendar{
		bot:        bot,
		logger:     logger,
		kb:         make([][]telebot.InlineButton, 0),
		currYear:   calendarOptions.initialYear,
		currMonth:  calendarOptions.initialMonth,
		yearsRange: calendarOptions.yearsRange,
		language:   calendarOptions.language,
		buttons:    btns,
	}

	// Регистрируем кнопки единожды для календаря:
	for buttonName, handler := range handlers {
		bot.Handle(btns[buttonName], handler(cal))
	}

	return cal, nil
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

// Clears the calendar's keyboard.
func (cal *Calendar) clearKeyboard() {
	cal.kb = make([][]telebot.InlineButton, 0)
}

// Builds a full row width button with a displayed month's name
// The button represents a list of all months when clicked.
func (cal *Calendar) addMonthYearRow() {
	var row []telebot.InlineButton

	btn := telebot.InlineButton{
		Unique: cal.buttons[monthsPerYearButton].Unique,
		Text:   fmt.Sprintf("%s %v", cal.getMonthDisplayName(cal.currMonth), cal.currYear),
	}

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
			Unique: cal.buttons[pickedMonthButton].Unique,
			Text:   monthName,
			Data:   strconv.Itoa(i),
		}

		row = append(row, monthBtn)

		// Arranging the months in 2 columns
		if i%monthsPerRowButtonsCount == 0 {
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
		btn := telebot.InlineButton{
			Unique: utils.GenUniqueParam("weekday_" + strconv.Itoa(i)),
			Text:   wd,
		}

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
	if weekdayNumber == 0 && cal.language == RussianLangAbbr { // russian Sunday exception
		weekdayNumber = 7
	}

	// The difference between English and Russian weekdays order
	// en: Sunday (0), Monday (1), Tuesday (3), ...
	// ru: Monday (1), Tuesday (2), ..., Sunday (7)
	if cal.language != RussianLangAbbr {
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
			Unique: cal.buttons[selectedDayButton].Unique,
			Text:   dayText,
			Data:   dayText,
		}

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

	prev := telebot.InlineButton{
		Unique: cal.buttons[previousMonthButton].Unique,
		Text:   "＜",
	}

	// Hide "prev" button if it rests on the range
	if cal.currYear <= cal.yearsRange[0] && cal.currMonth == 1 {
		prev.Unique = cal.buttons[ignoreQueryButton].Unique
		prev.Text = ""
	}

	next := telebot.InlineButton{
		Unique: cal.buttons[nextMonthButton].Unique,
		Text:   "＞",
	}

	// Hide "next" button if it rests on the range
	if cal.currYear >= cal.yearsRange[1] && cal.currMonth == monthsPerYear {
		next.Unique = cal.buttons[ignoreQueryButton].Unique
		next.Text = ""
	}

	row = append(row, prev, next)
	cal.addRowToKeyboard(&row)
}

// Builds a row of back and menu buttons.
func (cal *Calendar) addBackAndMenuButtonsRow() {
	var row []telebot.InlineButton

	if b, exists := cal.buttons[backButton]; exists {
		row = append(row, *b)
	}

	row = append(row, buttons.Menu)
	cal.addRowToKeyboard(&row)
}

// Returns a formatted date string from the selected date.
func (cal *Calendar) genDateStrFromDay(day int) string {
	return time.Date(cal.currYear, cal.currMonth, day,
		0, 0, 0, 0, time.UTC).Format(dateFormat)
}

// Utility function for passing a row to the calendar's keyboard.
func (cal *Calendar) addRowToKeyboard(row *[]telebot.InlineButton) {
	cal.kb = append(cal.kb, *row)
}

// Inserts an empty button that doesn't process queries
// into the keyboard row.
func (cal *Calendar) addEmptyCell(row *[]telebot.InlineButton) {
	cell := telebot.InlineButton{
		Unique: cal.buttons[ignoreQueryButton].Unique,
		Text:   " ",
	}

	*row = append(*row, cell)
}

// Returns the name of the month in the selected language.
func (cal *Calendar) getMonthDisplayName(month time.Month) string {
	if cal.language == RussianLangAbbr {
		return RussianMonths[month]
	}

	return month.String()
}

// Returns the array of the weekdays names in the selected language.
func (cal *Calendar) getWeekdaysDisplayArray() [AmountOfDaysInWeek]string {
	if cal.language == RussianLangAbbr {
		return RussianWeekdaysAbbrs
	}

	return EnglishWeekdaysAbbrs
}
