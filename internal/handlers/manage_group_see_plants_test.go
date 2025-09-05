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

func TestManageGroupSeePlantsCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	for _, tc := range []testCase{
		{
			name:          "success — one row of plants",
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
				plant1 := &entities.Plant{ID: 1, Title: "Rose", GroupID: 10}
				plant2 := &entities.Plant{ID: 2, Title: "Tulip", GroupID: 10}
				plants := []entities.Plant{*plant1, *plant2}
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

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupSeePlants).Return(nil)
			},
		},
		{
			name:          "success — multiple rows",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				group := &entities.Group{ID: 10, Title: "Garden"}
				plants := []entities.Plant{
					{ID: 1, Title: "A", GroupID: 10},
					{ID: 2, Title: "B", GroupID: 10},
					{ID: 3, Title: "C", GroupID: 10},
					{ID: 4, Title: "D", GroupID: 10},
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
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupSeePlants).Return(nil)
			},
		},
		{
			name:          "success — no plants in group",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				group := &entities.Group{ID: 10, Title: "Empty", Description: "No plants"}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Group{ID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return([]entities.Plant{}, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupSeePlants).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
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
				temp := &entities.Temporary{UserID: 123, Data: []byte(`{"id": "invalid"}`)}
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
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Group{ID: 10})}
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group plants fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				group := &entities.Group{ID: 10}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Group{ID: 10})}
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
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				group := &entities.Group{ID: 10, Title: "Garden"}
				plants := []entities.Plant{{ID: 1, Title: "Rose", GroupID: 10}}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Group{ID: 10})}
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
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				group := &entities.Group{ID: 10, Title: "Garden"}
				plants := []entities.Plant{{ID: 1, Title: "Rose", GroupID: 10}}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Group{ID: 10})}
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroupSeePlants).Return(assert.AnError)
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

			handler := ManageGroupSeePlantsCallback(nil, mockUsecases, mockLogger)

			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
