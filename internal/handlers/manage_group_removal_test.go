package handlers

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
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

func TestManageGroupRemovalCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	for _, tc := range []testCase{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				group := &entities.Group{
					ID:               10,
					Title:            "Garden",
					Description:      "My outdoor garden",
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
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupRemoval).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
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
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid group data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				temp := &entities.Temporary{
					UserID: 123,
					Data:   []byte(`{"id": "invalid"}`),
				}
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
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

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
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupRemoval).Return(assert.AnError)
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

			handler := ManageGroupRemovalCallback(nil, mockUsecases, mockLogger)

			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBackToManageGroupActionCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	sender := &telebot.User{ID: 123}

	for _, tc := range []testCase{
		{
			name:          "success — group has plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Garden",
					Description:      "My outdoor garden",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				plant := &entities.Plant{ID: 1, Title: "Rose", GroupID: 10}
				plants := []entities.Plant{*plant}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				// Arrange: моки
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)

				// Act: проверка Send
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupAction).Return(nil)
			},
		},
		{
			name:          "success — no plants in group",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{
					ID:               10,
					Title:            "Empty",
					Description:      "No plants",
					LastWateringDate: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
					WateringInterval: 7,
					NextWateringDate: time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				}
				var plants []entities.Plant // пусто
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupAction).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes() // ✅ Обязательно
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
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes() // ✅
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
					Data:   []byte(`{"id": "invalid"}`),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes() // ✅
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
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes() // ✅
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group plants fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{ID: 10}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes() // ✅
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{ID: 10, Title: "Garden"}
				plants := []entities.Plant{{ID: 1, Title: "Rose", GroupID: 10}}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes() // ✅
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
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
			name:          "set temporary step fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				group := &entities.Group{ID: 10, Title: "Garden"}
				plants := []entities.Plant{{ID: 1, Title: "Rose", GroupID: 10}}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes() // ✅
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupAction).Return(assert.AnError)
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

			handler := BackToManageGroupActionCallback(nil, mockUsecases, mockLogger)

			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfirmGroupRemovalCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	sender := &telebot.User{ID: 123}
	callback := &telebot.Callback{ID: "callback_123"}

	for _, tc := range []testCase{
		{
			name:          "success — user has groups and plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				user := &entities.User{ID: 500, TelegramID: 123}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(&telebot.CallbackResponse{
					CallbackID: "callback_123",
					Text:       texts.GroupDeleted,
				}).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(500).Return(2, nil)
				mockUsecases.EXPECT().CountUserPlants(500).Return(5, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
			},
		},
		{
			name:          "success — user has groups but no plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				user := &entities.User{ID: 500, TelegramID: 123}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(500).Return(1, nil)
				mockUsecases.EXPECT().CountUserPlants(500).Return(0, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
			},
		},
		{
			name:          "success — user has plants but no other groups",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				user := &entities.User{ID: 500, TelegramID: 123}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(500).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(500).Return(3, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
			},
		},
		{
			name:          "success — user has no groups and no plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				user := &entities.User{ID: 500, TelegramID: 123}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(500).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(500).Return(0, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
			},
		},
		{
			name:          "send error",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				user := &entities.User{ID: 500, TelegramID: 123}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(500).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(500).Return(0, nil)

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
			name:          "count user plants error",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				user := &entities.User{ID: 500, TelegramID: 123}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(500).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(500).Return(0, assert.AnError)
			},
		},
		{
			name:          "get user temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid group data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   []byte(`{"id": "invalid"}`),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockLogger.EXPECT().Error(
					"Failed to get Group from Temporary",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "delete group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(assert.AnError)
			},
		},
		{
			name:          "respond callback fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(assert.AnError)
				mockLogger.EXPECT().Error(
					"Failed to send Response",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(assert.AnError)
				mockLogger.EXPECT().Error(
					"Failed to delete message",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get user by telegram id fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "count user groups fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				user := &entities.User{ID: 500, TelegramID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(500).Return(0, assert.AnError)
			},
		},
		{
			name:          "reset temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}
				user := &entities.User{ID: 500, TelegramID: 123}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeleteGroup(10).Return(nil)
				mockCtx.EXPECT().Respond(gomock.Any()).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(500).Return(1, nil)
				mockUsecases.EXPECT().CountUserPlants(500).Return(1, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().ResetTemporary(123).Return(assert.AnError)
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

			handler := ConfirmGroupRemovalCallback(nil, mockUsecases, mockLogger)

			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
