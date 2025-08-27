package middlewares_test

import (
	"github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/middlewares"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
)

func TestLogging_Middleware(t *testing.T) {
	tests := []struct {
		name         string
		setupContext func() telebot.Context
		setupMocks   func(logger *mocks.MockLogger)
	}{
		{
			name: "Message context - should log message info",
			setupContext: func() telebot.Context {
				bot := &telebot.Bot{}
				msg := &telebot.Message{
					Sender: &telebot.User{ID: 12345},
					Text:   "Привет, бот!",
				}

				return telebot.NewContext(bot, telebot.Update{Message: msg})
			},
			setupMocks: func(logger *mocks.MockLogger) {
				logger.
					EXPECT().
					Info(
						"Received new message",
						"From", int64(12345),
						"Message", "Привет, бот!",
					).
					Times(1)
			},
		},
		{
			name: "Callback context - should log callback info",
			setupContext: func() telebot.Context {
				bot := &telebot.Bot{}
				callback := &telebot.Callback{
					Unique: "menu:open",
					Data:   "open_profile",
					Message: &telebot.Message{
						Sender: &telebot.User{ID: 67890},
					},
				}
				return telebot.NewContext(bot, telebot.Update{Callback: callback})
			},
			setupMocks: func(logger *mocks.MockLogger) {
				logger.
					EXPECT().
					Info(
						"Received new callback",
						"From", int64(67890),
						"Unique", "menu:open",
						"Data", "open_profile",
					).
					Times(1)
			},
		},
	}

	ctrl := gomock.NewController(t)
	mockLogger := mocks.NewMockLogger(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сбрасываем ожидания перед каждым тестом
			ctrl = gomock.NewController(t)
			mockLogger = mocks.NewMockLogger(ctrl)

			// Настройка контекста
			var ctx telebot.Context
			if tt.setupContext != nil {
				ctx = tt.setupContext()
			}

			// Настройка моков
			if tt.setupMocks != nil {
				tt.setupMocks(mockLogger)
			}

			// Создаём middleware
			middlewareFunc := middlewares.Logging(mockLogger)

			// Dummy next handler
			next := func(c telebot.Context) error {
				return nil
			}

			// Выполняем middleware
			handler := middlewareFunc(next)
			err := handler(ctx)

			// Проверки
			assert.NoError(t, err)
		})
	}
}
