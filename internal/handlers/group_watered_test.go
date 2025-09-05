package handlers

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	mockusecases "github.com/DKhorkov/plantsCareTelegramBot/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
	"time"
)

func TestGroupWateredCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger)
		contextData   string
		callback      *telebot.Callback
	}

	sender := &telebot.User{ID: 123}
	chat := &telebot.Chat{ID: 456}
	message := &telebot.Message{ID: 789}
	callbackID := "callback_123"

	for _, tc := range []testCase{
		{
			name:          "success — group watered, markup removed, response sent",
			errorExpected: false,
			contextData:   "10",
			callback: &telebot.Callback{
				ID:      callbackID,
				Sender:  sender,
				Message: message,
			},
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}

				// Общие ожидания
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()
				mockCtx.EXPECT().Callback().Return(&telebot.Callback{
					ID:      callbackID,
					Sender:  sender,
					Message: message,
				}).AnyTimes()

				// Получаем группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Обновляем дату полива
				mockUsecases.EXPECT().UpdateGroupLastWateringDate(
					10,
					gomock.AssignableToTypeOf(time.Time{}),
				).Return(group, nil)

				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()
				mockBot.EXPECT().EditReplyMarkup(message, &telebot.ReplyMarkup{}).Return(message, nil)

				// Отправляем ответ
				mockCtx.EXPECT().Respond(&telebot.CallbackResponse{
					CallbackID: callbackID,
					Text:       texts.GroupWatered,
				}).Return(nil)
			},
		},
		{
			name:          "parse groupID fails",
			errorExpected: true,
			contextData:   "invalid",
			callback: &telebot.Callback{
				ID:      callbackID,
				Sender:  sender,
				Message: message,
			},
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				// Общие ожидания
				mockCtx.EXPECT().Data().Return("invalid").AnyTimes()
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()

				// Логгируем ошибку
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
			callback: &telebot.Callback{
				ID:      callbackID,
				Sender:  sender,
				Message: message,
			},
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				// Общие ожидания
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(&telebot.Callback{
					ID:      callbackID,
					Sender:  sender,
					Message: message,
				}).AnyTimes()

				// Ошибка получения группы
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "update last watering date fails",
			errorExpected: true,
			contextData:   "10",
			callback: &telebot.Callback{
				ID:      callbackID,
				Sender:  sender,
				Message: message,
			},
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}

				// Общие ожидания
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(&telebot.Callback{
					ID:      callbackID,
					Sender:  sender,
					Message: message,
				}).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()

				// Получаем группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Ошибка обновления даты
				mockUsecases.EXPECT().UpdateGroupLastWateringDate(
					10,
					gomock.AssignableToTypeOf(time.Time{}),
				).Return(nil, assert.AnError)
			},
		},
		{
			name:          "callback is nil",
			errorExpected: true,
			contextData:   "10",
			callback:      nil,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				// Общие ожидания
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()
				mockCtx.EXPECT().Callback().Return((*telebot.Callback)(nil)).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				// Получаем группу
				mockUsecases.EXPECT().GetGroup(10).Return(&entities.Group{ID: 10}, nil)

				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}

				// Обновляем дату полива
				mockUsecases.EXPECT().UpdateGroupLastWateringDate(
					10,
					gomock.AssignableToTypeOf(time.Time{}),
				).Return(group, nil)

				// Логгируем предупреждение
				mockLogger.EXPECT().Warn(
					"Failed to send Response due to nil callback",
					"Message", message,
					"Sender", sender,
					"Chat", chat,
					"Callback", (*telebot.Callback)(nil),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "edit reply markup fails",
			errorExpected: true,
			contextData:   "10",
			callback: &telebot.Callback{
				ID:      callbackID,
				Sender:  sender,
				Message: message,
			},
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}

				// Общие ожидания
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()
				mockCtx.EXPECT().Callback().Return(&telebot.Callback{
					ID:      callbackID,
					Sender:  sender,
					Message: message,
				}).AnyTimes()

				// Получаем группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Обновляем дату
				mockUsecases.EXPECT().UpdateGroupLastWateringDate(
					10,
					gomock.AssignableToTypeOf(time.Time{}),
				).Return(group, nil)

				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()
				mockBot.EXPECT().EditReplyMarkup(message, &telebot.ReplyMarkup{}).Return(message, assert.AnError)

				// Логгируем ошибку
				mockLogger.EXPECT().Error(
					"Failed to delete ReplyMarkup",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "respond fails",
			errorExpected: true,
			contextData:   "10",
			callback: &telebot.Callback{
				ID:      callbackID,
				Sender:  sender,
				Message: message,
			},
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}

				// Общие ожидания
				mockCtx.EXPECT().Data().Return("10").AnyTimes()
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()
				mockCtx.EXPECT().Callback().Return(&telebot.Callback{
					ID:      callbackID,
					Sender:  sender,
					Message: message,
				}).AnyTimes()

				// Получаем группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Обновляем дату
				mockUsecases.EXPECT().UpdateGroupLastWateringDate(
					10,
					gomock.AssignableToTypeOf(time.Time{}),
				).Return(group, nil)

				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()
				mockBot.EXPECT().EditReplyMarkup(message, &telebot.ReplyMarkup{}).Return(message, nil)

				// Ошибка ответа
				mockCtx.EXPECT().Respond(&telebot.CallbackResponse{
					CallbackID: callbackID,
					Text:       texts.GroupWatered,
				}).Return(assert.AnError)

				// Логгируем ошибку
				mockLogger.EXPECT().Error(
					"Failed to send Response",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
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

			handler := GroupWateredCallback(nil, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
