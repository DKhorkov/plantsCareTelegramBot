package handlers

import (
	"github.com/DKhorkov/libs/logging"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	mockusecases "github.com/DKhorkov/plantsCareTelegramBot/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

// Заглушка для телеграм-контекста
func dummyHandler(c telebot.Context) error {
	return nil
}

func TestPrepare(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockUseCases := mockusecases.NewMockUseCases(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Примеры команд
	cmdStart := "/start"
	cmdMenu := "btn_menu"
	cmdPlants := "btn_my_plants"

	// Примеры обработчиков
	handlerStart := func(bot interfaces.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
		assert.Same(t, bot, mockBot)
		assert.Same(t, useCases, mockUseCases)
		assert.Same(t, logger, mockLogger)
		return dummyHandler
	}

	handlerMenu := func(bot interfaces.Bot, useCases interfaces.UseCases, logger logging.Logger) telebot.HandlerFunc {
		assert.Same(t, bot, mockBot)
		assert.Same(t, useCases, mockUseCases)
		assert.Same(t, logger, mockLogger)
		return dummyHandler
	}

	tests := []struct {
		name      string
		handlers  map[any]interfaces.Handler
		setupMock func()
	}{
		{
			name: "registers_all_handlers",
			handlers: map[any]interfaces.Handler{
				cmdStart: handlerStart,
				cmdMenu:  handlerMenu,
			},
			setupMock: func() {
				mockBot.EXPECT().Handle(cmdStart, gomock.Any()).Times(1)
				mockBot.EXPECT().Handle(cmdMenu, gomock.Any()).Times(1)
			},
		},
		{
			name:     "handles_empty_handlers_map",
			handlers: map[any]interfaces.Handler{},
			setupMock: func() {
				// Не должно быть вызовов Handle
			},
		},
		{
			name:     "handles_nil_handlers_map",
			handlers: nil,
			setupMock: func() {
				// Должно безопасно пройти
			},
		},
		{
			name: "registers_multiple_handlers_with_same_dependency_injection",
			handlers: map[any]interfaces.Handler{
				cmdStart:  handlerStart,
				cmdMenu:   handlerMenu,
				cmdPlants: handlerStart, // тот же обработчик — тоже ок
			},
			setupMock: func() {
				mockBot.EXPECT().Handle(cmdStart, gomock.Any()).Times(1)
				mockBot.EXPECT().Handle(cmdMenu, gomock.Any()).Times(1)
				mockBot.EXPECT().Handle(cmdPlants, gomock.Any()).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			}

			Prepare(mockBot, mockUseCases, mockLogger, tt.handlers)
		})
	}
}
