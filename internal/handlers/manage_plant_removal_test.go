package handlers

import (
	"encoding/json"
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
)

func mustMarshal(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	return data
}

func TestConfirmPlantRemovalCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		tempData      []byte
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success",
			errorExpected: false,
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				callback := &telebot.Callback{ID: "callback-123"}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()

				// GetUserTemporary → возвращает Temporary с Data
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)

				// Удаляем растение по ID из десериализованного объекта
				mockUsecases.EXPECT().DeletePlant(1).Return(nil)

				// Подтверждаем callback
				mockCtx.EXPECT().Respond(&telebot.CallbackResponse{
					CallbackID: "callback-123",
					Text:       texts.PlantDeleted,
				}).Return(nil)

				// Удаляем сообщение
				mockCtx.EXPECT().Delete().Return(nil)

				// Получаем пользователя
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)

				// Считаем группы и растения
				mockUsecases.EXPECT().CountUserGroups(1).Return(2, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(3, nil)

				// Отправляем фото с меню
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				// Сбрасываем временные данные
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
			},
		},
		{
			name:          "get user temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "temporary has invalid plant data",
			errorExpected: true,
			tempData:      []byte(`{"id": "not-a-number"}`), // invalid JSON
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				temp := &entities.Temporary{UserID: 123, Data: []byte(`{"id": "not-a-number"}`)}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockLogger.EXPECT().Error(
					"Failed to get Plant from Temporary",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "delete plant fails",
			errorExpected: true,
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeletePlant(1).Return(assert.AnError)
			},
		},
		{
			name:          "respond fails",
			errorExpected: true,
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				callback := &telebot.Callback{ID: "callback-123"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeletePlant(1).Return(nil)
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(assert.AnError)
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
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				callback := &telebot.Callback{ID: "callback-123"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeletePlant(1).Return(nil)
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(nil)
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
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				callback := &telebot.Callback{ID: "callback-123"}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeletePlant(1).Return(nil)
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(nil, assert.AnError)
			},
		},
		{
			name:          "count user groups fails",
			errorExpected: true,
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				callback := &telebot.Callback{ID: "callback-123"}
				user := &entities.User{ID: 1}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeletePlant(1).Return(nil)
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, assert.AnError)
			},
		},
		{
			name:          "count user plants fails",
			errorExpected: true,
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				callback := &telebot.Callback{ID: "callback-123"}
				user := &entities.User{ID: 1}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeletePlant(1).Return(nil)
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(2, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				callback := &telebot.Callback{ID: "callback-123"}
				user := &entities.User{ID: 1}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeletePlant(1).Return(nil)
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(2, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(3, nil)
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
			name:          "reset temporary fails",
			errorExpected: true,
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				callback := &telebot.Callback{ID: "callback-123"}
				user := &entities.User{ID: 1}
				temp := &entities.Temporary{UserID: 123, Data: mustMarshal(t, &entities.Plant{ID: 1})}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockUsecases.EXPECT().GetUserTemporary(123).Return(temp, nil)
				mockUsecases.EXPECT().DeletePlant(1).Return(nil)
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(nil)
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(2, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(3, nil)
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(assert.AnError)
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

			// Устанавливаем моки
			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			// Создаём хэндлер
			handler := ConfirmPlantRemovalCallback(mockBot, mockUsecases, mockLogger)

			// Вызываем хэндлер
			err := handler(mockCtx)

			// Проверяем ошибку
			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManagePlantRemovalCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
		tempData      []byte
	}{
		{
			name:          "success",
			errorExpected: false,
			tempData:      mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, Title: "Rose", Description: "Beautiful", GroupID: 10, Photo: []byte{0xFF, 0xD8}}
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

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantRemoval).Return(nil)
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
			tempData:      []byte(`{"id": "invalid"}`),
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
			tempData:      mustMarshal(t, &entities.Plant{ID: 1}),
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
			tempData:      mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
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
			tempData:      mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, Title: "Rose", Description: "Beautiful", GroupID: 10, Photo: []byte{0xFF, 0xD8}}
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
			tempData:      mustMarshal(t, &entities.Plant{ID: 1, GroupID: 10}),
			setupMocks: func(mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				plant := &entities.Plant{ID: 1, Title: "Rose", Description: "Beautiful", GroupID: 10, Photo: []byte{0xFF, 0xD8}}
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

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantRemoval).Return(assert.AnError)
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

			handler := ManagePlantRemovalCallback(mockBot, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestBackToManagePlantActionCallback(t *testing.T) {
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
				plant := &entities.Plant{ID: 1, Title: "Rose", Description: "Beautiful", GroupID: 10, Photo: []byte{0xFF, 0xD8}}
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

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantAction).Return(nil)
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
				plant := &entities.Plant{ID: 1, Title: "Rose", Description: "Beautiful", GroupID: 10, Photo: []byte{0xFF, 0xD8}}
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
				plant := &entities.Plant{ID: 1, Title: "Rose", Description: "Beautiful", GroupID: 10, Photo: []byte{0xFF, 0xD8}}
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

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantAction).Return(assert.AnError)
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

			handler := BackToManagePlantActionCallback(mockBot, mockUsecases, mockLogger)
			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
