package calendar

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"strconv"
	"testing"
	"time"
)

func TestNewCalendar_DefaultOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	now := time.Now()
	expectedYear := now.Year()
	expectedMonth := now.Month()

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(1)
	}

	cal, err := NewCalendar(mockBot, mockLogger)
	require.NoError(t, err)

	require.Equal(t, expectedYear, cal.currYear)
	require.Equal(t, expectedMonth, cal.currMonth)
	require.Equal(t, [2]int{expectedYear, expectedYear}, cal.yearsRange)
	require.Equal(t, RussianLangAbbr, cal.language)

	// Проверим, что кнопки созданы и имеют уникальные префиксы
	for name := range handlers {
		btn := cal.buttons[name]
		require.NotNil(t, btn)
	}
}

func TestNewCalendar_WithOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	backBtn := telebot.InlineButton{Text: "Назад", Unique: "back_123"}

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(1)
	}

	cal, err := NewCalendar(mockBot, mockLogger,
		WithYearsRange([2]int{2020, 2030}),
		WithLanguage(EnglishLangAbbr),
		WithBackButton(backBtn),
	)
	require.NoError(t, err)

	require.Equal(t, 2025, cal.currYear)
	require.Equal(t, time.Now().Month(), cal.currMonth)
	require.Equal(t, [2]int{2020, 2030}, cal.yearsRange)
	require.Equal(t, EnglishLangAbbr, cal.language)
	require.NotNil(t, cal.buttons[backButton])
	require.Equal(t, "back_123", cal.buttons[backButton].Unique)
}

func TestNewCalendar_InvalidOptions(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Некорректный диапазон: начало > конец
	cal, err := NewCalendar(mockBot, mockLogger,
		WithYearsRange([2]int{2025, 2020}),
	)
	require.Error(t, err)
	require.Nil(t, cal)
}

func TestCalendar_getMonthDisplayName(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(2)
	}

	cal, err := NewCalendar(mockBot, mockLogger, WithLanguage(RussianLangAbbr))
	require.NoError(t, err)

	tests := []struct {
		name     string
		month    time.Month
		expected string
	}{
		{"январь", time.January, "Январь"},
		{"март", time.March, "Март"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cal.getMonthDisplayName(tt.month)
			require.Equal(t, tt.expected, result)
		})
	}

	// Английский
	calEn, err := NewCalendar(mockBot, mockLogger, WithLanguage(EnglishLangAbbr))
	require.NoError(t, err)

	require.Equal(t, "January", calEn.getMonthDisplayName(time.January))
}

func TestCalendar_getWeekdaysDisplayArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(2)
	}

	calRu, err := NewCalendar(mockBot, mockLogger, WithLanguage(RussianLangAbbr))
	require.NoError(t, err)

	ruWeekdays := calRu.getWeekdaysDisplayArray()
	require.Equal(t, "Пн", ruWeekdays[0])
	require.Equal(t, "Вс", ruWeekdays[6])

	calEn, err := NewCalendar(mockBot, mockLogger, WithLanguage(EnglishLangAbbr))
	require.NoError(t, err)

	enWeekdays := calEn.getWeekdaysDisplayArray()
	require.Equal(t, "Su", enWeekdays[0])
	require.Equal(t, "Sa", enWeekdays[6])
}

func TestCalendar_genDateStrFromDay(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(1)
	}

	cal, err := NewCalendar(
		mockBot,
		mockLogger,
		WithInitialYear(2025),
		WithInitialMonth(time.February),
	)
	require.NoError(t, err)

	tests := []struct {
		day      int
		expected string
	}{
		{1, "01.02.2025"},
		{10, "10.02.2025"},
		{28, "28.02.2025"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := cal.genDateStrFromDay(tt.day)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestCalendar_addMonthYearRow(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(1)
	}

	cal, err := NewCalendar(mockBot, mockLogger,
		WithInitialYear(2025),
		WithInitialMonth(time.March),
		WithLanguage(RussianLangAbbr),
	)
	require.NoError(t, err)

	cal.clearKeyboard()
	cal.addMonthYearRow()

	require.Len(t, cal.kb, 1)
	btn := cal.kb[0][0]
	require.Equal(t, "Март 2025", btn.Text)
	require.Equal(t, cal.buttons[monthsPerYearButton].Unique, btn.Unique)
}

func TestCalendar_addDaysRows(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(1)
	}

	cal, err := NewCalendar(
		mockBot,
		mockLogger,
		WithInitialMonth(time.January),
		WithInitialYear(2025),
	)
	require.NoError(t, err)

	cal.clearKeyboard()
	cal.addDaysRows()

	// 2025-01-01 — среда → смещение = 2 (Пн, Вт — пустые)
	firstRow := cal.kb[0]
	require.Equal(t, " ", firstRow[0].Text) // Пн
	require.Equal(t, " ", firstRow[1].Text) // Вт
	require.Equal(t, "1", firstRow[2].Text) // Ср
}

func TestCalendar_addControlButtonsRow(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(1)
	}

	cal, err := NewCalendar(mockBot, mockLogger,
		WithInitialMonth(time.January),
		WithInitialYear(2025),
		WithYearsRange([2]int{2025, 2030}),
	)
	require.NoError(t, err)

	cal.clearKeyboard()
	cal.addControlButtonsRow()

	prevBtn := cal.kb[0][0]
	nextBtn := cal.kb[0][1]

	require.Equal(t, "", prevBtn.Text) // скрыт
	require.Equal(t, "＞", nextBtn.Text)

	require.Equal(t, cal.buttons[ignoreQueryButton].Unique, prevBtn.Unique)
	require.Equal(t, cal.buttons[nextMonthButton].Unique, nextBtn.Unique)
}

func TestCalendar_addBackAndMenuButtonsRow(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
	for range handlers {
		mockBot.
			EXPECT().
			Handle(gomock.Any(), gomock.Any()).
			Times(1)
	}

	backBtn := telebot.InlineButton{Text: "Назад", Unique: "back_123"}
	cal, err := NewCalendar(mockBot, mockLogger, WithBackButton(backBtn))
	require.NoError(t, err)

	cal.clearKeyboard()
	cal.addBackAndMenuButtonsRow()

	row := cal.kb[0]
	require.Len(t, row, 2)
	require.Equal(t, "Назад", row[0].Text)
	require.Equal(t, buttons.Menu.Text, row[1].Text)
}

func TestCalendar_getMonthPickKeyboard(t *testing.T) {
	tests := []struct {
		name         string
		language     string
		hasBackBtn   bool
		expectedRows [][monthsPerRowButtonsCount]string
	}{
		{
			name:       "ru with back",
			language:   RussianLangAbbr,
			hasBackBtn: true,
			expectedRows: [][monthsPerRowButtonsCount]string{
				{"Январь", "Февраль"},
				{"Март", "Апрель"},
				{"Май", "Июнь"},
				{"Июль", "Август"},
				{"Сентябрь", "Октябрь"},
				{"Ноябрь", "Декабрь"},
				{"Назад", buttons.Menu.Text},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
			for range handlers {
				mockBot.
					EXPECT().
					Handle(gomock.Any(), gomock.Any()).
					Times(1)
			}

			var backBtn *telebot.InlineButton
			if tt.hasBackBtn {
				backBtn = &telebot.InlineButton{Text: "Назад", Unique: "back_custom"}
			}

			opts := []Option{WithLanguage(tt.language)}
			if backBtn != nil {
				opts = append(opts, WithBackButton(*backBtn))
			}

			cal, err := NewCalendar(mockBot, mockLogger, opts...)
			require.NoError(t, err)

			kb := cal.getMonthPickKeyboard()

			// Проверяем первую строку
			require.Len(t, kb, len(tt.expectedRows))

			for i, row := range kb {
				for j, btn := range row {
					require.Equal(t, tt.expectedRows[i][j], btn.Text)
				}
			}
		})
	}
}

func TestCalendar_addWeekdaysRow(t *testing.T) {
	tests := []struct {
		name     string
		language string
		expected []string
	}{
		{
			name:     "ru weekdays",
			language: RussianLangAbbr,
			expected: []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"},
		},
		{
			name:     "en weekdays",
			language: EnglishLangAbbr,
			expected: []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
			for range handlers {
				mockBot.
					EXPECT().
					Handle(gomock.Any(), gomock.Any()).
					Times(1)
			}

			cal, err := NewCalendar(mockBot, mockLogger, WithLanguage(tt.language))
			require.NoError(t, err)

			cal.clearKeyboard()
			cal.addWeekdaysRow()

			require.Len(t, cal.kb, 1)
			row := cal.kb[0]
			require.Len(t, row, 7)

			for i, expected := range tt.expected {
				require.Equal(t, expected, row[i].Text)
				require.NotEmpty(t, row[i].Unique)
				require.Contains(t, row[i].Unique, "weekday_"+strconv.Itoa(i))
			}
		})
	}
}

func TestCalendar_GetKeyboard(t *testing.T) {
	tests := []struct {
		name           string
		year           int
		month          time.Month
		language       string
		yearsRange     [2]int
		hasBackButton  bool
		controlButtons struct{ prev, next bool }
	}{
		{
			name:           "ru 2025-02 (feb) middle range, has back",
			year:           2025,
			month:          time.February,
			language:       RussianLangAbbr,
			yearsRange:     [2]int{2020, 2030},
			hasBackButton:  true,
			controlButtons: struct{ prev, next bool }{true, true},
		},
		{
			name:           "en 2020-01 (jan) start of range, no back",
			year:           2020,
			month:          time.January,
			language:       EnglishLangAbbr,
			yearsRange:     [2]int{2020, 2030},
			hasBackButton:  false,
			controlButtons: struct{ prev, next bool }{false, true},
		},
		{
			name:           "ru 2030-12 (dec) end of range, has back",
			year:           2030,
			month:          time.December,
			language:       RussianLangAbbr,
			yearsRange:     [2]int{2020, 2030},
			hasBackButton:  true,
			controlButtons: struct{ prev, next bool }{true, false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			// Ожидаем, что bot.Handle вызывается для всех кнопок из handlers
			for range handlers {
				mockBot.
					EXPECT().
					Handle(gomock.Any(), gomock.Any()).
					Times(1)
			}

			// Подготавливаем опции
			opts := []Option{
				WithInitialYear(tt.year),
				WithInitialMonth(tt.month),
				WithYearsRange(tt.yearsRange),
				WithLanguage(tt.language),
			}

			var backBtn *telebot.InlineButton
			if tt.hasBackButton {
				backBtn = &telebot.InlineButton{Text: "Назад", Unique: "back_custom"}
				opts = append(opts, WithBackButton(*backBtn))
			}

			// Создаём календарь
			cal, err := NewCalendar(mockBot, mockLogger, opts...)
			require.NoError(t, err)

			// Вызываем тестируемый метод
			kb := cal.GetKeyboard()

			// 1. Проверяем строку: Месяц и год
			monthName := cal.getMonthDisplayName(tt.month)
			require.GreaterOrEqual(t, len(kb), 5) // 5+ строк: месяц, дни недели, дни, контроли, меню

			monthYearRow := kb[0]
			require.Len(t, monthYearRow, 1)
			require.Equal(t, monthName+" "+strconv.Itoa(tt.year), monthYearRow[0].Text)
			require.Equal(t, cal.buttons[monthsPerYearButton].Unique, monthYearRow[0].Unique)

			// 2. Проверяем строку дней недели
			weekdaysRow := kb[1]
			require.Len(t, weekdaysRow, 7)

			expectedWeekdays := cal.getWeekdaysDisplayArray()
			for i, wd := range expectedWeekdays {
				require.Equal(t, wd, weekdaysRow[i].Text)
				require.Contains(t, weekdaysRow[i].Unique, "weekday_"+strconv.Itoa(i))
			}

			// 3. Проверяем первую строку с днями
			daysRow := kb[2]
			require.GreaterOrEqual(t, len(daysRow), 7)

			// 4. Проверяем контрольные кнопки (предыдущий/следующий месяц)
			controlRow := kb[len(kb)-2]
			require.Len(t, controlRow, 2)

			if tt.controlButtons.prev {
				require.Equal(t, "＜", controlRow[0].Text)
				require.Equal(t, cal.buttons[previousMonthButton].Unique, controlRow[0].Unique)
			} else {
				require.Equal(t, "", controlRow[0].Text)
				require.Equal(t, cal.buttons[ignoreQueryButton].Unique, controlRow[0].Unique)
			}

			if tt.controlButtons.next {
				require.Equal(t, "＞", controlRow[1].Text)
				require.Equal(t, cal.buttons[nextMonthButton].Unique, controlRow[1].Unique)
			} else {
				require.Equal(t, "", controlRow[1].Text)
				require.Equal(t, cal.buttons[ignoreQueryButton].Unique, controlRow[1].Unique)
			}

			// 5. Проверяем строку с «Назад» и «Меню»
			menuRow := kb[len(kb)-1]
			if tt.hasBackButton {
				require.Len(t, menuRow, 2)
				require.Equal(t, "Назад", menuRow[0].Text)
				require.Equal(t, buttons.Menu.Text, menuRow[1].Text)
			} else {
				require.Len(t, menuRow, 1)
				require.Equal(t, buttons.Menu.Text, menuRow[0].Text)
			}
		})
	}
}
