package handlers

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/texts"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	mockusecases "github.com/DKhorkov/plantsCareTelegramBot/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestHelp(t *testing.T) {
	tests := []struct {
		name          string
		errorExpected bool
		setupMocks    func(*mockbot.MockBot, *mockbot.MockContext, *mockusecases.MockUseCases, *mocklogging.MockLogger)
	}{
		{
			name:          "success",
			errorExpected: false,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Delete().Return(nil)

				mockCtx.EXPECT().Send(texts.OnHelp).Return(nil)
			},
		},
		{
			name:          "delete fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Delete().Return(assert.AnError)

				mockLogger.EXPECT().Error(
					"Failed to delete /help message",
					"Error", gomock.Any(),
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
		{
			name:          "send fails",
			errorExpected: true,
			setupMocks: func(mockBot *mockbot.MockBot, mockCtx *mockbot.MockContext, mockUsecases *mockusecases.MockUseCases, mockLogger *mocklogging.MockLogger) {
				mockCtx.EXPECT().Delete().Return(nil)

				mockCtx.EXPECT().Send(texts.OnHelp).Return(assert.AnError)

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

			handler := Help(mockBot, mockUsecases, mockLogger)

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
