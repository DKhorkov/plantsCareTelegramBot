package handlers

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	mockusecases "github.com/DKhorkov/plantsCareTelegramBot/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
	"time"
)

func TestManageGroupCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger)
		contextData   string
	}

	sender := &telebot.User{ID: 123}
	chat := &telebot.Chat{ID: 456}

	for _, tc := range []testCase{
		{
			name:          "success — group has plants, menu with all buttons sent",
			errorExpected: false,
			contextData:   "10",
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				plants := []entities.Plant{
					{ID: 1, Title: "Phalaenopsis"},
				}

				// Контекст
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Data().Return("10").AnyTimes()

				// Удаляем сообщение
				mockCtx.EXPECT().Delete().Return(nil)

				// Получаем группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Получаем растения
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)

				// Ожидаем отправку сообщения
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				// Вызываем ManageGroup
				mockUsecases.EXPECT().ManageGroup(123, 10).Return(nil)
			},
		},
		{
			name:          "success — group has no plants, menu without see plants button",
			errorExpected: false,
			contextData:   "10",
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Cacti",
					Description:      "Desert beauty",
					LastWateringDate: time.Date(2024, 5, 20, 0, 0, 0, 0, time.UTC),
					WateringInterval: 14,
					NextWateringDate: time.Date(2024, 6, 3, 0, 0, 0, 0, time.UTC),
				}
				var plants []entities.Plant

				// Контекст
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Data().Return("10").AnyTimes()

				// Удаляем сообщение
				mockCtx.EXPECT().Delete().Return(nil)

				// Получаем группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Получаем растения
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)

				// Ожидаем отправку сообщения
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				// Вызываем ManageGroup
				mockUsecases.EXPECT().ManageGroup(123, 10).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
			contextData:   "10",
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(assert.AnError)

				mockLogger.EXPECT().Error(
					"Failed to delete message",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "parse groupID fails",
			errorExpected: true,
			contextData:   "invalid",
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Data().Return("invalid").AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockLogger.EXPECT().Error(
					"Failed to parse groupID",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			contextData:   "10",
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group plants fails",
			errorExpected: true,
			contextData:   "10",
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:    10,
					Title: "Orchids",
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			contextData:   "10",
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				plants := []entities.Plant{{ID: 1, Title: "Phalaenopsis"}}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(assert.AnError)

				mockLogger.EXPECT().Error(
					"Failed to send message",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "manage group fails",
			errorExpected: true,
			contextData:   "10",
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				plants := []entities.Plant{{ID: 1, Title: "Phalaenopsis"}}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().ManageGroup(123, 10).Return(assert.AnError)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)
			mockBot := mockbot.NewMockBot(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			// bot не используется
			handler := ManageGroupCallback(nil, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
