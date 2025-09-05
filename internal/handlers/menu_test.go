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
)

func TestBackToMenu(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success with groups and plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(2, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(3, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)
			},
		},
		{
			name:          "success with groups only",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(2, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)
			},
		},
		{
			name:          "success with plants only",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(3, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)
			},
		},
		{
			name:          "success no groups no plants",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, nil)

				mockCtx.EXPECT().Send(
					gomock.AssignableToTypeOf(&telebot.Photo{}),
					gomock.AssignableToTypeOf(&telebot.ReplyMarkup{}),
				).Return(nil)
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
			name:          "reset temporary fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(assert.AnError)
			},
		},
		{
			name:          "count user groups fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, assert.AnError)
			},
		},
		{
			name:          "count user plants fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, assert.AnError)
			},
		},
		{
			name:          "send fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				sender := &telebot.User{ID: 123}
				chat := &telebot.Chat{ID: 456}
				user := &entities.User{ID: 1}

				mockCtx.EXPECT().Sender().Return(sender).AnyTimes()
				mockCtx.EXPECT().Chat().Return(chat).AnyTimes()
				mockCtx.EXPECT().Delete().Return(nil)

				mockUsecases.EXPECT().GetUserByTelegramID(123).Return(user, nil)
				mockUsecases.EXPECT().ResetTemporary(123).Return(nil)
				mockUsecases.EXPECT().CountUserGroups(1).Return(0, nil)
				mockUsecases.EXPECT().CountUserPlants(1).Return(0, nil)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockBot := mockbot.NewMockBot(ctrl)
			mockCtx := mockbot.NewMockContext(ctrl)
			mockUsecases := mockusecases.NewMockUseCases(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			handler := BackToMenu(mockBot, mockUsecases, mockLogger)

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
