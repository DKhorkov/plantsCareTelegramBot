package handlers

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/libs/pointers"
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

func TestStart(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success first time user",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123, Username: "testuser", FirstName: "Test", LastName: "User", IsBot: false}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().SaveUser(gomock.Any()).Return(1, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, nil)
				mockCtx.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
			},
		},
		{
			name:          "count user groups error",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123, Username: "testuser", FirstName: "Test", LastName: "User", IsBot: false}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().SaveUser(gomock.Any()).Return(1, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, assert.AnError)
			},
		},
		{
			name:          "count user plants error",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123, Username: "testuser", FirstName: "Test", LastName: "User", IsBot: false}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().SaveUser(gomock.Any()).Return(1, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, assert.AnError)
			},
		},
		{
			name:          "ResetTemporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123, Username: "testuser", FirstName: "Test", LastName: "User", IsBot: false}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().SaveUser(gomock.Any()).Return(1, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, nil)
				mockCtx.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(assert.AnError)
			},
		},
		{
			name:          "delete fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Delete().Return(assert.AnError)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
		{
			name:          "save user fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Sender().Return(&telebot.User{ID: 123}).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().SaveUser(gomock.Any()).Return(0, assert.AnError)
			},
		},
		{
			name:          "send message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().SaveUser(gomock.Any()).Return(1, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, nil)
				mockCtx.EXPECT().Send(gomock.Any(), gomock.Any()).Return(assert.AnError)
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).MinTimes(1)
			},
		},
		{
			name:          "has groups and plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)
				mockUsecases.EXPECT().SaveUser(gomock.Any()).Return(1, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(2, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(5, nil)
				mockCtx.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
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

			handler := Start(mockBot, mockUsecases, mockLogger)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAddGroupCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				// Мокируем Sender, Chat, Callback — могут вызываться несколько раз
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()

				// UseCases: GetUserByTelegramID → *User
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)

				// CountUserGroups → не достигнут лимит
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)

				// Успешное удаление сообщения
				mockCtx.EXPECT().Delete().Return(nil)

				mockCtx.EXPECT().Bot().Return(mockBot)

				// Bot.Send → возвращает *Message
				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(&telebot.Message{ID: 789}, nil)

				// SetTemporaryStep и SetTemporaryMessage
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.AddGroupTitle).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, pointers.New(789)).Return(nil)
			},
		},
		{
			name:          "delete fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)

				// Удаление падает
				mockCtx.EXPECT().Delete().Return(assert.AnError)

				// Логируем ошибку
				mockLogger.EXPECT().Error(
					"Failed to delete message",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "get user fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()

				// Ошибка получения пользователя
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return((*entities.User)(nil), assert.AnError)
			},
		},
		{
			name:          "count groups fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, assert.AnError)
			},
		},
		{
			name:          "groups limit reached — respond success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()
				mockCtx.EXPECT().Chat().Return(&telebot.Chat{ID: 456}).AnyTimes()

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(groupsPerUserLimit, nil)

				// Respond успешный
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(nil)
			},
		},
		{
			name:          "groups limit reached — nil callback",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				message := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(&telebot.Chat{ID: 456}).AnyTimes()

				// Callback == nil
				mockCtx.EXPECT().Callback().Return((*telebot.Callback)(nil)).AnyTimes()

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(groupsPerUserLimit, nil)

				// Логируем предупреждение
				mockLogger.EXPECT().Warn(
					"Failed to send Response due to nil callback",
					"Message", message,
					"Sender", sender,
					"Chat", gomock.Any(),
					"Callback", (*telebot.Callback)(nil),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "groups limit reached — respond fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()
				mockCtx.EXPECT().Chat().Return(&telebot.Chat{ID: 456}).AnyTimes()

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(groupsPerUserLimit, nil)

				// Respond возвращает ошибку
				mockCtx.EXPECT().Respond(gomock.AssignableToTypeOf(&telebot.CallbackResponse{})).Return(assert.AnError)

				// Логируем ошибку
				mockLogger.EXPECT().Error(
					"Failed to send Response",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "send fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockCtx.EXPECT().Bot().Return(mockBot)

				// Send падает
				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return((*telebot.Message)(nil), assert.AnError)

				// Логируем ошибку
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
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockCtx.EXPECT().Bot().Return(mockBot)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(&telebot.Message{ID: 789}, nil)

				// SetTemporaryStep падает
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.AddGroupTitle).Return(assert.AnError)
			},
		},
		{
			name:          "set temporary message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				message := &telebot.Message{ID: 789}
				callback := &telebot.Callback{ID: "abc", Message: message}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Callback().Return(callback).AnyTimes()
				mockCtx.EXPECT().Message().Return(message).AnyTimes()

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(&entities.User{ID: 1}, nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockCtx.EXPECT().Delete().Return(nil)

				mockCtx.EXPECT().Bot().Return(mockBot)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(&telebot.Message{ID: 789}, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.AddGroupTitle).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, pointers.New(789)).Return(assert.AnError)
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

			handler := AddGroupCallback(mockBot, mockUsecases, mockLogger)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAddPlantCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				message := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockCtx.EXPECT().Delete().Return(nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(message, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.AddPlantTitle).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, pointers.New(789)).Return(nil)
			},
		},
		{
			name:          "delete fails",
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
			name:          "send fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockCtx.EXPECT().Delete().Return(nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return((*telebot.Message)(nil), assert.AnError)

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
				message := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockCtx.EXPECT().Delete().Return(nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(message, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.AddPlantTitle).Return(assert.AnError)
			},
		},
		{
			name:          "set temporary message fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				message := &telebot.Message{ID: 789}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Bot().Return(mockBot).AnyTimes()

				mockCtx.EXPECT().Delete().Return(nil)

				mockBot.EXPECT().Send(
					chat,
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(message, nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.AddPlantTitle).Return(nil)
				mockUsecases.EXPECT().SetTemporaryMessage(123, pointers.New(789)).Return(assert.AnError)
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

			handler := AddPlantCallback(mockBot, mockUsecases, mockLogger)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManagePlantsCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success with groups with plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}
				groups := []entities.Group{
					{ID: 1, Title: "Garden", UserID: 1},
					{ID: 2, Title: "Balcony", UserID: 1},
				}
				plants1 := []entities.Plant{{ID: 1, GroupID: 1}}
				var plants2 []entities.Plant // без растений

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)

				// Мокаем GetGroupPlants для каждой группы
				mockUsecases.EXPECT().GetGroupPlants(1).Return(plants1, nil)
				mockUsecases.EXPECT().GetGroupPlants(2).Return(plants2, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantsChooseGroup).Return(nil)
			},
		},
		{
			name:          "no groups with plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}
				groups := []entities.Group{
					{ID: 1, Title: "Empty", UserID: 1},
				}
				var plants []entities.Plant

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)
				mockUsecases.EXPECT().GetGroupPlants(1).Return(plants, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantsChooseGroup).Return(nil)
			},
		},
		{
			name:          "delete fails",
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
			name:          "get user fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return((*entities.User)(nil), assert.AnError)
			},
		},
		{
			name:          "get user groups fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "get group plants fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				user := &entities.User{ID: 1}
				groups := []entities.Group{{ID: 1, UserID: 1}}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)
				mockUsecases.EXPECT().GetGroupPlants(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}
				groups := []entities.Group{{ID: 1, UserID: 1}}
				plants := []entities.Plant{{ID: 1, GroupID: 1}}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)
				mockUsecases.EXPECT().GetGroupPlants(1).Return(plants, nil)

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
				user := &entities.User{ID: 1}
				groups := []entities.Group{{ID: 1, UserID: 1}}
				plants := []entities.Plant{{ID: 1, GroupID: 1}}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)
				mockUsecases.EXPECT().GetGroupPlants(1).Return(plants, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManagePlantsChooseGroup).Return(assert.AnError)
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

			handler := ManagePlantsCallback(mockBot, mockUsecases, mockLogger)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestManageGroupsCallback(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success with groups",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}
				groups := []entities.Group{
					{ID: 1, Title: "Garden", UserID: 1},
					{ID: 2, Title: "Balcony", UserID: 1},
				}

				// Мокаем базовые вызовы
				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				// UseCases: GetUserByTelegramID
				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)

				// GetUserGroups
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)

				// context.Send — с фото и меню
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				// SetTemporaryStep
				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroup).Return(nil)
			},
		},
		{
			name:          "success no groups",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}
				var groups []entities.Group

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)

				// Даже без групп — отправляем меню с одной кнопкой BackToStart
				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroup).Return(nil)
			},
		},
		{
			name:          "delete fails",
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
			name:          "get user fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return((*entities.User)(nil), assert.AnError)
			},
		},
		{
			name:          "get user groups fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(nil, assert.AnError)
			},
		},
		{
			name:          "send fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}
				groups := []entities.Group{{ID: 1, UserID: 1}}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)

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
				user := &entities.User{ID: 1}
				groups := []entities.Group{{ID: 1, UserID: 1}}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().GetUserGroups(1).Return(groups, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)

				mockUsecases.EXPECT().SetTemporaryStep(123, steps.ManageGroup).Return(assert.AnError)
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

			handler := ManageGroupsCallback(mockBot, mockUsecases, mockLogger)

			if tt.setupMocks != nil {
				tt.setupMocks(mockBot, mockCtx, mockUsecases, mockLogger)
			}

			err := handler(mockCtx)

			if tt.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
