package handlers

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	mockusecases "github.com/DKhorkov/plantsCareTelegramBot/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
	"time"
)

func TestManageGroupChangeCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	sender := &telebot.User{ID: 123}
	callback := &telebot.Callback{ID: "callback_123"}

	for _, tc := range []testCase{
		{
			name:          "success — group data sent",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				// Контекст
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				// Удаляем сообщение
				mockCtx.EXPECT().Delete().Return(nil)

				// Получаем временные данные
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				// Десериализуем группу
				// → temp.GetGroup() — не юзкейс, не мокируется

				// Получаем полную группу из БД
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Отправляем фото с меню
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				// Устанавливаем шаг
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupChange).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
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
			name:          "get user temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid group data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   []byte(`invalid json`),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				mockLogger.EXPECT().Error(
					"Failed to get Group from Temporary",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

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
			name:          "set temporary step fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupChange).Return(assert.AnError)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			handler := ManageGroupChangeCallback(nil, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManageGroupChangeTitleCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	sender := &telebot.User{ID: 123}
	chat := &telebot.Chat{ID: 456}
	callback := &telebot.Callback{ID: "callback_123"}

	for _, tc := range []testCase{
		{
			name:          "success — message sent, step and message ID set",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				sentMessage := &telebot.Message{ID: 789}

				// Контекст
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				// Удаляем текущее сообщение
				mockCtx.EXPECT().Delete().Return(nil)

				// Получаем временные данные
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				// Десериализуем группу → temp.GetGroup() — не мокируется

				// Получаем полную группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Отправляем сообщение — ожидаем, что Bot.Send вернёт msg
				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				// Устанавливаем шаг
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupTitle).Return(nil)

				// Сохраняем ID сообщения
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
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
			name:          "get user temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid group data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   []byte(`invalid json`),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				mockLogger.EXPECT().Error(
					"Failed to get Group from Temporary",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil, assert.AnError)

				mockLogger.EXPECT().Error(
					"Failed to send message",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "set temporary step fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupTitle).Return(assert.AnError)
			},
		},
		{
			name:          "set temporary message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupTitle).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(assert.AnError)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			handler := ManageGroupChangeTitleCallback(nil, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManageGroupChangeDescriptionCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	sender := &telebot.User{ID: 123}
	chat := &telebot.Chat{ID: 456}
	callback := &telebot.Callback{ID: "callback_123"}

	for _, tc := range []testCase{
		{
			name:          "success — message sent, step and message ID set",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				sentMessage := &telebot.Message{ID: 789}

				// Контекст
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				// Удаляем текущее сообщение
				mockCtx.EXPECT().Delete().Return(nil)

				// Получаем временные данные
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				// Десериализуем группу → temp.GetGroup() — не мокируется

				// Получаем полную группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Отправляем сообщение
				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				// Устанавливаем шаг
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupDescription).Return(nil)

				// Сохраняем ID сообщения
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
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
			name:          "get user temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid group data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   []byte(`invalid json`),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				mockLogger.EXPECT().Error(
					"Failed to get Group from Temporary",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil, assert.AnError)

				mockLogger.EXPECT().Error(
					"Failed to send message",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "set temporary step fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupDescription).Return(assert.AnError)
			},
		},
		{
			name:          "set temporary message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupDescription).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(assert.AnError)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			handler := ManageGroupChangeDescriptionCallback(nil, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManageGroupChangeLastWateringDateCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	sender := &telebot.User{ID: 123}
	chat := &telebot.Chat{ID: 456}

	for _, tc := range []testCase{
		{
			name:          "success — message sent, step set",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).AnyTimes()

				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				// Контекст
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				// Удаляем сообщение
				mockCtx.EXPECT().Delete().Return(nil)

				// Получаем временные данные
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				// Получаем полную группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Бот передаётся в NewCalendar — ожидаем, что он будет использован
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				// Ожидаем вызов Send
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				// Устанавливаем шаг
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupLastWateringDate).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
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
			name:          "get user temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid group data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   []byte(`invalid json`),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				mockLogger.EXPECT().Error(
					"Failed to get Group from Temporary",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).AnyTimes()

				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

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
			name:          "set temporary step fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockBot.EXPECT().Handle(gomock.Any(), gomock.Any()).AnyTimes()

				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupLastWateringDate).Return(assert.AnError)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			handler := ManageGroupChangeLastWateringDateCallback(mockBot, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManageGroupChangeWateringIntervalCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger)
	}

	sender := &telebot.User{ID: 123}
	chat := &telebot.Chat{ID: 456}

	for _, tc := range []testCase{
		{
			name:          "success — message sent, step set, buttons wrapped correctly",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				// Контекст
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				// Удаляем сообщение
				mockCtx.EXPECT().Delete().Return(nil)

				// Получаем временные данные
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				// Десериализуем группу → temp.GetGroup() — не мокируется

				// Получаем полную группу
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// Ожидаем отправку сообщения
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				// Устанавливаем шаг
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupWateringInterval).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
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
			name:          "get user temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid group data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   []byte(`invalid json`),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				mockLogger.EXPECT().Error(
					"Failed to get Group from Temporary",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

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
			name:          "set temporary step fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Orchids",
					Description:      "White and pink",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangeGroupWateringInterval).Return(assert.AnError)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			// bot не используется, передаём nil
			handler := ManageGroupChangeWateringIntervalCallback(nil, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
