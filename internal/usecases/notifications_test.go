package usecases

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	mockstorage "github.com/DKhorkov/plantsCareTelegramBot/mocks/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestNotificationsUseCases_SaveNotification(t *testing.T) {
	inputNotification := entities.Notification{
		GroupID:   1,
		MessageID: 100,
		Text:      "Не забудьте полить цветы!",
		SentAt:    time.Time{}, // будет заполнено позже
	}

	tests := []struct {
		name         string
		notification entities.Notification
		setupMocks   func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		want         *entities.Notification
		wantErr      bool
	}{
		{
			name:         "Success - notification saved",
			notification: inputNotification,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					SaveNotification(inputNotification).
					Return(42, nil).
					Times(1)
			},
			want: &entities.Notification{
				ID:        42,
				GroupID:   1,
				MessageID: 100,
				Text:      "Не забудьте полить цветы!",
				SentAt:    time.Time{},
			},
			wantErr: false,
		},
		{
			name:         "Failure - storage error",
			notification: inputNotification,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					SaveNotification(inputNotification).
					Return(0, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to save Notification for Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &notificationsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.SaveNotification(tt.notification)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.GroupID, got.GroupID)
				assert.Equal(t, tt.want.MessageID, got.MessageID)
				assert.Equal(t, tt.want.Text, got.Text)
				assert.Equal(t, tt.want.SentAt, got.SentAt) // zero time — ок, если не устанавливается в use case
			}
		})
	}
}
