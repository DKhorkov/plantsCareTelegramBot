package calendar

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
	"time"
)

func TestMonthsPerYearCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mocklogging.MockLogger)
	}{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:          "edit fails",
			errorExpected: false, // ошибка логируется, но не возвращается
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, assert.AnError)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
		{
			name:          "respond fails",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(assert.AnError)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			// Регистрация обработчиков
			for range handlers {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).Times(1)
			}

			cal, err := NewCalendar(mockBot, mockLogger, WithInitialYear(2025), WithInitialMonth(time.March))
			require.NoError(t, err)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockLogger)
			}

			handler := MonthsPerYearCallback(cal)
			err = handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPickedMonthCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mocklogging.MockLogger)
	}{
		{
			name:          "valid month",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Data().Return("4").AnyTimes()
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:          "invalid data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Data().Return("abc").AnyTimes()
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
		{
			name:          "edit fails",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Data().Return("6").AnyTimes()
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, assert.AnError)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
		{
			name:          "respond fails",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Data().Return("5").AnyTimes()
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(assert.AnError)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			for range handlers {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).Times(1)
			}

			cal, err := NewCalendar(mockBot, mockLogger, WithInitialYear(2025), WithInitialMonth(time.January))
			require.NoError(t, err)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockLogger)
			}

			handler := PickedMonthCallback(cal)
			err = handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestIgnoreQueryCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mocklogging.MockLogger)
	}{
		{
			name:          "respond success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Respond().Return(nil)
			},
		},
		{
			name:          "respond fails",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Respond().Return(assert.AnError)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			for range handlers {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).Times(1)
			}

			cal, err := NewCalendar(mockBot, mockLogger)
			require.NoError(t, err)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockLogger)
			}

			handler := IgnoreQueryCallback(cal)
			err = handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestSelectedDayCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mocklogging.MockLogger)
	}{
		{
			name:          "valid day",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Data().Return("5").AnyTimes()
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().ProcessUpdate(gomock.Any()).Do(func(update telebot.Update) {
					require.Equal(t, "05.02.2025", update.Message.Payload)
				})
			},
		},
		{
			name:          "invalid data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Data().Return("xyz").AnyTimes()
			},
		},
		{
			name:          "respond fails",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Data().Return("15").AnyTimes()
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(assert.AnError)
				mockBot.EXPECT().ProcessUpdate(gomock.Any()).Do(func(update telebot.Update) {
					require.Equal(t, "15.02.2025", update.Message.Payload)
				})
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			for range handlers {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).Times(1)
			}

			cal, err := NewCalendar(mockBot, mockLogger, WithInitialYear(2025), WithInitialMonth(time.February))
			require.NoError(t, err)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockLogger)
			}

			handler := SelectedDayCallback(cal)
			err = handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestPreviousMonthCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		initialYear   int
		initialMonth  time.Month
		expectedYear  int
		expectedMonth time.Month
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mocklogging.MockLogger)
	}{
		{
			name:          "from March 2025 to February 2025",
			errorExpected: false,
			initialYear:   2025,
			initialMonth:  time.March,
			expectedYear:  2025,
			expectedMonth: time.February,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:          "from January 2025 to December 2024",
			errorExpected: false,
			initialYear:   2025,
			initialMonth:  time.January,
			expectedYear:  2024,
			expectedMonth: time.December,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:          "edit fails",
			errorExpected: false,
			initialYear:   2025,
			initialMonth:  time.April,
			expectedYear:  2025,
			expectedMonth: time.March,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, assert.AnError)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
		{
			name:          "respond fails",
			errorExpected: false,
			initialYear:   2025,
			initialMonth:  time.April,
			expectedYear:  2025,
			expectedMonth: time.March,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(assert.AnError)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			// Регистрация обработчиков
			for range handlers {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).Times(1)
			}

			cal, err := NewCalendar(mockBot, mockLogger, WithInitialYear(tt.initialYear), WithInitialMonth(tt.initialMonth))
			cal.yearsRange = [2]int{2020, 2030} // Устанавливаем диапазон лет
			require.NoError(t, err)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockLogger)
			}

			handler := PreviousMonthCallback(cal)
			err = handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			// Проверяем результат
			require.Equal(t, tt.expectedMonth, cal.currMonth)
			require.Equal(t, tt.expectedYear, cal.currYear)
		})
	}
}

func TestNextMonthCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		initialYear   int
		initialMonth  time.Month
		expectedYear  int
		expectedMonth time.Month
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mocklogging.MockLogger)
	}{
		{
			name:          "from March 2025 to April 2025",
			errorExpected: false,
			initialYear:   2025,
			initialMonth:  time.March,
			expectedYear:  2025,
			expectedMonth: time.April,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:          "from December 2025 to January 2026",
			errorExpected: false,
			initialYear:   2025,
			initialMonth:  time.December,
			expectedYear:  2026,
			expectedMonth: time.January,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
			},
		},
		{
			name:          "edit fails",
			errorExpected: false,
			initialYear:   2025,
			initialMonth:  time.February,
			expectedYear:  2025,
			expectedMonth: time.March,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(nil)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, assert.AnError)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
		{
			name:          "respond fails",
			errorExpected: false,
			initialYear:   2025,
			initialMonth:  time.February,
			expectedYear:  2025,
			expectedMonth: time.March,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockLogger *mocklogging.MockLogger) {
				msg := &telebot.Message{Chat: &telebot.Chat{ID: 123}}
				mockCtx.EXPECT().Message().Return(msg).AnyTimes()
				mockCtx.EXPECT().Respond().Return(assert.AnError)
				mockBot.EXPECT().EditReplyMarkup(msg, gomock.Any()).Return(nil, nil)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			for range handlers {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).Times(1)
			}

			cal, err := NewCalendar(mockBot, mockLogger, WithInitialYear(tt.initialYear), WithInitialMonth(tt.initialMonth))
			cal.yearsRange = [2]int{2020, 2030}
			require.NoError(t, err)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockLogger)
			}

			handler := NextMonthCallback(cal)
			err = handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedMonth, cal.currMonth)
			require.Equal(t, tt.expectedYear, cal.currYear)
		})
	}
}
