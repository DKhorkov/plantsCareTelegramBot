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
)

func TestManagePlantChangeCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
				}

				// Ожидания
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantChange).Return(nil)
			},
		},
		{
			name:          "delete message fails",
			errorExpected: true,
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
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
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid plant data",
			errorExpected: true,
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				temp := &entities.Temporary{UserID: 123, Data: []byte(`{"id": "invalid"}`)}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
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
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, GroupID: 10}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
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
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantChange).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockCtx, mockUsecases, mockLogger)
			}

			handler := ManagePlantChangeCallback(mockBot, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManagePlantChangeTitleCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	tests := []testCase{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
				}
				sentMessage := &telebot.Message{ID: 789}

				// Ожидания
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantTitle).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(nil)
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
			name:          "temporary has invalid plant data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				temp := &entities.Temporary{UserID: 123, Data: []byte(`{"id": "invalid"}`)}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
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
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, GroupID: 10}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
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
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantTitle).Return(assert.AnError)
			},
		},
		{
			name:          "set temporary message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantTitle).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			handler := ManagePlantChangeTitleCallback(mockBot, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManagePlantChangeDescriptionCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	tests := []testCase{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
				}
				sentMessage := &telebot.Message{ID: 789}

				// Ожидания
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantDescription).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(nil)
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
			name:          "temporary has invalid plant data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				temp := &entities.Temporary{UserID: 123, Data: []byte(`{"id": "invalid"}`)}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
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
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, GroupID: 10}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
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
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantDescription).Return(assert.AnError)
			},
		},
		{
			name:          "set temporary message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantDescription).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			handler := ManagePlantChangeDescriptionCallback(mockBot, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManagePlantChangeGroupCallback(t *testing.T) {
	type testCase struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}

	tests := []testCase{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				currentGroup := &entities.Group{ID: 10, Title: "Garden"}
				otherGroups := []entities.Group{
					{ID: 11, Title: "Indoor"},
					{ID: 12, Title: "Balcony"},
					{ID: 13, Title: "Terrace"},
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123}),
				}

				// Ожидания
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(currentGroup, nil)
				mockUsecases.EXPECT().GetUserGroups(123).Return(otherGroups, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantGroup).Return(nil)
			},
		},
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				currentGroup := &entities.Group{ID: 10, Title: "Garden"}
				otherGroups := []entities.Group{
					*currentGroup,
					{ID: 11, Title: "Indoor"},
					{ID: 12, Title: "Balcony"},
					{ID: 13, Title: "Terrace"},
				}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123}),
				}

				// Ожидания
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(currentGroup, nil)
				mockUsecases.EXPECT().GetUserGroups(123).Return(otherGroups, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantGroup).Return(nil)
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
			name:          "temporary has invalid plant data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				temp := &entities.Temporary{UserID: 123, Data: []byte(`{"id": "invalid"}`)}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
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
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get current group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, GroupID: 10}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get user groups fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, GroupID: 10, UserID: 123}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123})}
				currentGroup := &entities.Group{ID: 10, Title: "Garden"}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(currentGroup, nil)
				mockUsecases.EXPECT().GetUserGroups(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				currentGroup := &entities.Group{ID: 10, Title: "Garden"}
				otherGroups := []entities.Group{
					{ID: 11, Title: "Indoor"},
				}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(currentGroup, nil)
				mockUsecases.EXPECT().GetUserGroups(123).Return(otherGroups, nil)

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
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				currentGroup := &entities.Group{ID: 10, Title: "Garden"}
				otherGroups := []entities.Group{
					{ID: 11, Title: "Indoor"},
				}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(currentGroup, nil)
				mockUsecases.EXPECT().GetUserGroups(123).Return(otherGroups, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantGroup).Return(assert.AnError)
			},
		},
		{
			name:          "no other groups — only back and menu",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				currentGroup := &entities.Group{ID: 10, Title: "Garden"}
				var otherGroups []entities.Group // Нет других групп
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(currentGroup, nil)
				mockUsecases.EXPECT().GetUserGroups(123).Return(otherGroups, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantGroup).Return(nil)
			},
		},
		{
			name:          "one group — fits in one row",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				currentGroup := &entities.Group{ID: 10, Title: "Garden"}
				otherGroups := []entities.Group{
					{ID: 11, Title: "Indoor"},
				}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(currentGroup, nil)
				mockUsecases.EXPECT().GetUserGroups(123).Return(otherGroups, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantGroup).Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			handler := ManagePlantChangeGroupCallback(mockBot, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManagePlantChangePhotoCallback(t *testing.T) {
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
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{
					UserID: 123,
					Data:   mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123}),
				}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantPhoto).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(nil)
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
			name:          "temporary has invalid plant data",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				temp := &entities.Temporary{UserID: 123, Data: []byte(`{"id": "invalid"}`)}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
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
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, GroupID: 10}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123})}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
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
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123})}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantPhoto).Return(assert.AnError)
			},
		},
		{
			name:          "set temporary message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				plant := &entities.Plant{
					ID:          1,
					UserID:      123,
					Title:       "Rose",
					Description: "Beautiful flower",
					GroupID:     10,
					Photo:       []byte{0xFF, 0xD8},
				}
				group := &entities.Group{ID: 10, Title: "Garden"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10, UserID: 123})}
				sentMessage := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().GetPlant(1).Return(plant, nil)
				mockUsecases.EXPECT().GetGroup(10).Return(group, nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(sentMessage, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ChangePlantPhoto).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, &sentMessage.ID).Return(assert.AnError)
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

			handler := ManagePlantChangePhotoCallback(nil, mockUsecases, mockLogger)

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
