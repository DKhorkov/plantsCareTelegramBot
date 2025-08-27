package usecases

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	mockstorage "github.com/DKhorkov/plantsCareTelegramBot/mocks/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestUsersUseCases_SaveUser(t *testing.T) {
	now := time.Now()
	existingUser := entities.User{
		ID:         123,
		TelegramID: 5551234,
		Username:   "existing",
		Firstname:  "Иван",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	newUser := entities.User{
		TelegramID: 5559999,
		Username:   "newuser",
		Firstname:  "Петр",
	}

	temp := entities.Temporary{
		UserID: 456,
		Step:   steps.Start,
	}

	tests := []struct {
		name       string
		inputUser  entities.User
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantID     int
		wantErr    bool
	}{
		{
			name:      "Success - user already exists, return existing ID",
			inputUser: newUser,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				// Пользователь уже существует
				storage.
					EXPECT().
					GetUserByTelegramID(5559999).
					Return(&existingUser, nil).
					Times(1)
				// SaveUser и CreateTemporary не вызываются
			},
			wantID:  123,
			wantErr: false,
		},
		{
			name:      "Success - new user saved with temporary record",
			inputUser: newUser,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				// Пользователь не найден
				storage.
					EXPECT().
					GetUserByTelegramID(5559999).
					Return(nil, assert.AnError).
					Times(1)
				// Сохранение нового
				storage.
					EXPECT().
					SaveUser(newUser).
					Return(456, nil).
					Times(1)
				// Создание временной записи
				storage.
					EXPECT().
					CreateTemporary(temp).
					Return(nil).
					Times(1)
			},
			wantID:  456,
			wantErr: false,
		},
		{
			name:      "Failure - error on SaveUser",
			inputUser: newUser,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(5559999).
					Return(nil, assert.AnError).
					Times(1)
				storage.
					EXPECT().
					SaveUser(newUser).
					Return(0, assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to save User with telegramID=5559999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantID:  0,
			wantErr: true,
		},
		{
			name:      "Failure - error on CreateTemporary",
			inputUser: newUser,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(5559999).
					Return(nil, assert.AnError).
					Times(1)
				storage.
					EXPECT().
					SaveUser(newUser).
					Return(456, nil).
					Times(1)
				storage.
					EXPECT().
					CreateTemporary(temp).
					Return(assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to save Temporary for User with ID=456",
						"Temp", temp,
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantID:  0,
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

			useCases := &usersUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotID, err := useCases.SaveUser(tt.inputUser)

			assert.Equal(t, tt.wantID, gotID)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUsersUseCases_GetUserByTelegramID(t *testing.T) {
	now := time.Now()
	expectedUser := &entities.User{
		ID:         123,
		TelegramID: 5551234,
		Username:   "testuser",
		Firstname:  "Иван",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name       string
		telegramID int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantUser   *entities.User
		wantErr    bool
	}{
		{
			name:       "Success - user found by TelegramID",
			telegramID: 5551234,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(5551234).
					Return(expectedUser, nil).
					Times(1)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
		{
			name:       "Failure - storage error",
			telegramID: 5551234,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(5551234).
					Return(nil, assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to get User with telegramID=5551234",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantUser: nil,
			wantErr:  true,
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

			useCases := &usersUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.GetUserByTelegramID(tt.telegramID)

			assert.Equal(t, tt.wantUser, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestUsersUseCases_GetUserByID(t *testing.T) {
	now := time.Now()
	expectedUser := &entities.User{
		ID:         123,
		TelegramID: 5551234,
		Username:   "testuser",
		Firstname:  "Иван",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name       string
		userID     int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantUser   *entities.User
		wantErr    bool
	}{
		{
			name:   "Success - user found by ID",
			userID: 123,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByID(123).
					Return(expectedUser, nil).
					Times(1)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
		{
			name:   "Failure - storage error",
			userID: 123,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByID(123).
					Return(nil, assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to get User with ID=123",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantUser: nil,
			wantErr:  true,
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

			useCases := &usersUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.GetUserByID(tt.userID)

			assert.Equal(t, tt.wantUser, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
