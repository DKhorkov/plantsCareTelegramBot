package handlers

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/buttons"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	mockusecases "github.com/DKhorkov/plantsCareTelegramBot/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"testing"
)

func TestManagePlantCallback(t *testing.T) {
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
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Data().Return("1")
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// ✅ Мокаем context.Send, а не bot.Send
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().ManagePlant(123, 1).Return(nil)
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
			name:          "invalid plantID in context.Data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Data().Return("invalid-id")
				mockLogger.EXPECT().Error(
					"Failed to parse plantID",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get plant fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Data().Return("1")
				mockUsecases.EXPECT().GetPlant(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				plant := &entities.Plant{ID: 1, GroupID: 10}
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Data().Return("1")
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Data().Return("1")
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				// ✅ context.Send возвращает ошибку
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
			name:          "manage plant usecase fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Data().Return("1")
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().ManagePlant(123, 1).Return(assert.AnError)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tc.setupMocks != nil {
				tc.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			handler := ManagePlantCallback(nil, mockUsecases, mockLogger)

			// Act
			err := handler(mockCtx)

			// Assert
			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBackToManagePlantCallback(t *testing.T) {
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
				sender := telebot.User{ID: 123}
				plant := entities.Plant{
					ID:      1,
					UserID:  123,
					Title:   "Rose",
					GroupID: 10,
					Photo:   []byte{0xFF, 0xD8},
				}
				otherPlant := entities.Plant{
					ID:      2,
					UserID:  123,
					Title:   "Tulip",
					GroupID: 10,
				}
				groupPlants := []entities.Plant{plant, otherPlant}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123}),
				}

				mockCtx.EXPECT().Sender().Return(&sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(&plant, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(groupPlants, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlant).Return(nil)
			},
		},
		{
			name:          "success — multiple rows",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := entities.Plant{ID: 1, Title: "A", GroupID: 10}
				plants := []entities.Plant{
					plant,
					{ID: 2, Title: "B", GroupID: 10},
					{ID: 3, Title: "C", GroupID: 10},
					{ID: 4, Title: "D", GroupID: 10},
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(&plant, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlant).Return(nil)
			},
		},
		{
			name:          "success — no plants in group",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, GroupID: 10}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
				}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return([]entities.Plant{}, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).DoAndReturn(func(photo *telebot.Photo, markup *telebot.ReplyMarkup) error {
					require.Len(t, markup.InlineKeyboard, 1) // Только строка с Back и Menu
					require.Len(t, markup.InlineKeyboard[0], 2)
					require.Equal(t, buttons.BackToManagePlantsChooseGroup.Text, markup.InlineKeyboard[0][0].Text)
					require.Equal(t, buttons.Menu.Text, markup.InlineKeyboard[0][1].Text)
					return nil
				})

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlant).Return(nil)
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
			name:          "temporary has invalid plant data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				temp := &entities.Temporary{UserID: 123, Data: []byte(`{"id": "invalid"}`)}
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockLogger.EXPECT().Error(
					"Failed to get Plant from Temporary",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get plant fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group plants fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				plant := &entities.Plant{ID: 1, GroupID: 10}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				plant := entities.Plant{ID: 1, GroupID: 10}
				plants := []entities.Plant{plant}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(&plant, nil)
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
				plant := entities.Plant{ID: 1, GroupID: 10}
				plants := []entities.Plant{plant}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(&plant, nil)
				mockUsecases.EXPECT().GetGroupPlants(10).Return(plants, nil)
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlant).Return(assert.AnError)
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

			handler := BackToManagePlantCallback(nil, mockUsecases, mockLogger)

			err := handler(mockCtx)

			if tc.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
