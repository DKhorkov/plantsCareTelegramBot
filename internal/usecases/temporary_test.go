package usecases

import (
	"encoding/json"
	"errors"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	customerrors "github.com/DKhorkov/plantsCareTelegramBot/internal/errors"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
	mockstorage "github.com/DKhorkov/plantsCareTelegramBot/mocks/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestTemporaryUseCases_GetUserTemporary(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	temp := &entities.Temporary{ID: 1, UserID: 123, Step: 0}

	tests := []struct {
		name       string
		telegramID int
		setupMocks func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantTemp   *entities.Temporary
		wantErr    bool
	}{
		{
			name:       "Success - temporary user returned",
			telegramID: 456,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)
			},
			wantTemp: temp,
			wantErr:  false,
		},
		{
			name:       "Failure - GetUserByTelegramID error",
			telegramID: 456,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=456",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantTemp: nil,
			wantErr:  true,
		},
		{
			name:       "Failure - GetTemporaryByUserID error",
			telegramID: 456,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary for User with ID=123",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantTemp: nil,
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.GetUserTemporary(tt.telegramID)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTemp.ID, got.ID)
				assert.Equal(t, tt.wantTemp.UserID, got.UserID)
			}
		})
	}
}

func TestTemporaryUseCases_SetTemporaryStep(t *testing.T) {
	temp := &entities.Temporary{ID: 1, UserID: 123, Step: 0}

	tests := []struct {
		name       string
		telegramID int
		step       int
		setupMocks func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantErr    bool
	}{
		{
			name:       "Success - step updated",
			telegramID: 456,
			step:       5,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(&entities.User{ID: 123}, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, 5, temp.Step)
						assert.Equal(t, 1, temp.ID)
						return nil
					}).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:       "Failure - GetUserTemporary returns error",
			telegramID: 999,
			step:       5,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - UpdateTemporary returns error",
			telegramID: 456,
			step:       5,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(&entities.User{ID: 123}, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			err := useCases.SetTemporaryStep(tt.telegramID, tt.step)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestTemporaryUseCases_SetTemporaryMessage(t *testing.T) {
	temp := &entities.Temporary{ID: 1, UserID: 123, MessageID: nil}

	tests := []struct {
		name       string
		telegramID int
		messageID  *int
		setupMocks func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantErr    bool
	}{
		{
			name:       "Success - messageID set to new value",
			telegramID: 456,
			messageID:  func() *int { i := 100; return &i }(),
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(&entities.User{ID: 123}, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, 100, *temp.MessageID)
						return nil
					}).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:       "Success - messageID set to nil",
			telegramID: 456,
			messageID:  nil,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(&entities.User{ID: 123}, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Nil(t, temp.MessageID)
						return nil
					}).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:       "Failure - GetUserTemporary error",
			telegramID: 999,
			messageID:  nil,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - UpdateTemporary error",
			telegramID: 456,
			messageID:  func() *int { i := 100; return &i }(),
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(&entities.User{ID: 123}, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			err := useCases.SetTemporaryMessage(tt.telegramID, tt.messageID)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestTemporaryUseCases_AddGroupTitle(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	temp := &entities.Temporary{ID: 1, UserID: 123, Step: 0, MessageID: nil, Data: nil}

	tests := []struct {
		name        string
		telegramID  int
		title       string
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantGroup   *entities.Group
		wantErr     bool
		wantErrIs   error
		checkResult func(*testing.T, *entities.Temporary, *entities.Group)
	}{
		{
			name:       "Success - title added and temporary updated",
			telegramID: 456,
			title:      "Цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				// Получение пользователя
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				// Получение временной записи
				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				// Проверка на существование группы
				storage.
					EXPECT().
					GroupExists(gomock.Any()).
					DoAndReturn(func(g entities.Group) (bool, error) {
						assert.Equal(t, "Цветы", g.Title)
						assert.Equal(t, 123, g.UserID)
						return false, nil
					}).
					Times(1)

				// Обновление временной записи
				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.AddGroupDescription, temp.Step)
						assert.Nil(t, temp.MessageID) // сбрасывается

						// Проверка, что Data содержит корректную сериализованную группу
						var group entities.Group
						assert.NoError(t, json.Unmarshal(temp.Data, &group))
						assert.Equal(t, "Цветы", group.Title)
						assert.Equal(t, 123, group.UserID)
						return nil
					}).
					Times(1)
			},
			wantGroup: &entities.Group{UserID: 123, Title: "Цветы"},
			wantErr:   false,
			checkResult: func(t *testing.T, updatedTemp *entities.Temporary, returnedGroup *entities.Group) {
				assert.Equal(t, "Цветы", returnedGroup.Title)
				assert.Equal(t, 123, returnedGroup.UserID)
			},
		},
		{
			name:       "Failure - GetUserTemporary returns error",
			telegramID: 999,
			title:      "Цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - GroupExists returns error",
			telegramID: 456,
			title:      "Цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					GroupExists(gomock.Any()).
					Return(false, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						gomock.Any(),
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - group already exists",
			telegramID: 456,
			title:      "Цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					GroupExists(gomock.Any()).
					Return(true, nil).
					Times(1)
			},
			wantErr:   true,
			wantErrIs: customerrors.ErrGroupAlreadyExists,
		},
		{
			name:       "Failure - UpdateTemporary returns error",
			telegramID: 456,
			title:      "Цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					GroupExists(gomock.Any()).
					Return(false, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotGroup, err := useCases.AddGroupTitle(tt.telegramID, tt.title)

			if tt.wantErr {
				assert.Nil(t, gotGroup)
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotGroup)
				assert.Equal(t, tt.wantGroup.UserID, gotGroup.UserID)
				assert.Equal(t, tt.wantGroup.Title, gotGroup.Title)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, temp, gotGroup)
			}
		})
	}
}

func TestTemporaryUseCases_AddGroupDescription(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	existingGroup := &entities.Group{UserID: 123, Title: "Цветы"}
	tempData, _ := json.Marshal(existingGroup)
	temp := &entities.Temporary{ID: 1, UserID: 123, Data: tempData, Step: 0, MessageID: &[]int{100}[0]}

	tests := []struct {
		name        string
		telegramID  int
		description string
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantGroup   *entities.Group
		wantErr     bool
		checkResult func(*testing.T, *entities.Temporary, *entities.Group)
	}{
		{
			name:        "Success - description added and temporary updated",
			telegramID:  456,
			description: "Уход за цветами",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				// GetUserTemporary → GetUserByTelegramID + GetTemporaryByUserID
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				// UpdateTemporary — проверяем, что Step и Data обновились
				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.AddGroupLastWateringDate, temp.Step)
						assert.Nil(t, temp.MessageID) // сброшен

						var group entities.Group
						assert.NoError(t, json.Unmarshal(temp.Data, &group))
						assert.Equal(t, "Цветы", group.Title)
						assert.Equal(t, "Уход за цветами", group.Description)
						return nil
					}).
					Times(1)
			},
			wantGroup: &entities.Group{UserID: 123, Title: "Цветы", Description: "Уход за цветами"},
			wantErr:   false,
			checkResult: func(t *testing.T, updatedTemp *entities.Temporary, returnedGroup *entities.Group) {
				assert.Equal(t, "Уход за цветами", returnedGroup.Description)
				assert.Equal(t, steps.AddGroupLastWateringDate, updatedTemp.Step)
				assert.Nil(t, updatedTemp.MessageID)
			},
		},
		{
			name:        "Failure - GetUserTemporary returns error",
			telegramID:  999,
			description: "Описание",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:        "Failure - UpdateTemporary returns error",
			telegramID:  456,
			description: "Описание",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotGroup, err := useCases.AddGroupDescription(tt.telegramID, tt.description)

			if tt.wantErr {
				assert.Nil(t, gotGroup)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotGroup)
				assert.Equal(t, tt.wantGroup.Title, gotGroup.Title)
				assert.Equal(t, tt.wantGroup.Description, gotGroup.Description)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, temp, gotGroup)
			}
		})
	}
}

func TestTemporaryUseCases_AddGroupLastWateringDate(t *testing.T) {
	now := time.Now().Truncate(24 * time.Hour) // убираем наносекунды для точного сравнения
	user := &entities.User{ID: 123, TelegramID: 456}
	existingGroup := &entities.Group{UserID: 123, Title: "Цветы", Description: "Уход"}
	tempData, _ := json.Marshal(existingGroup)
	temp := &entities.Temporary{ID: 1, UserID: 123, Data: tempData, Step: steps.AddGroupLastWateringDate, MessageID: &[]int{100}[0]}

	tests := []struct {
		name             string
		telegramID       int
		lastWateringDate time.Time
		setupMocks       func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantGroup        *entities.Group
		wantErr          bool
		checkResult      func(*testing.T, *entities.Temporary, *entities.Group)
	}{
		{
			name:             "Success - last watering date added",
			telegramID:       456,
			lastWateringDate: now,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.AddGroupWateringInterval, temp.Step)

						var group entities.Group
						assert.NoError(t, json.Unmarshal(temp.Data, &group))
						assert.WithinDuration(t, now, group.LastWateringDate, time.Second)
						return nil
					}).
					Times(1)
			},
			wantGroup: &entities.Group{
				UserID:           123,
				Title:            "Цветы",
				Description:      "Уход",
				LastWateringDate: now,
			},
			wantErr: false,
			checkResult: func(t *testing.T, updatedTemp *entities.Temporary, returnedGroup *entities.Group) {
				assert.WithinDuration(t, now, returnedGroup.LastWateringDate, time.Second)
				assert.Equal(t, steps.AddGroupWateringInterval, updatedTemp.Step)
				// MessageID не сбрасывается в этом методе → может остаться
			},
		},
		{
			name:             "Failure - GetUserTemporary error",
			telegramID:       999,
			lastWateringDate: now,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:             "Failure - UpdateTemporary returns error",
			telegramID:       456,
			lastWateringDate: now,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotGroup, err := useCases.AddGroupLastWateringDate(tt.telegramID, tt.lastWateringDate)

			if tt.wantErr {
				assert.Nil(t, gotGroup)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotGroup)
				assert.Equal(t, tt.wantGroup.Title, gotGroup.Title)
				assert.WithinDuration(t, tt.wantGroup.LastWateringDate, gotGroup.LastWateringDate, time.Second)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, temp, gotGroup)
			}
		})
	}
}

func TestTemporaryUseCases_AddGroupWateringInterval(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	user := &entities.User{ID: 123, TelegramID: 456}
	lastWatering := now.AddDate(0, 0, -3) // 3 дня назад
	existingGroup := &entities.Group{
		UserID:           123,
		Title:            "Цветы",
		Description:      "Уход за цветами",
		LastWateringDate: lastWatering,
	}
	tempData, _ := json.Marshal(existingGroup)
	temp := &entities.Temporary{ID: 1, UserID: 123, Data: tempData, Step: steps.AddGroupWateringInterval}

	tests := []struct {
		name             string
		telegramID       int
		wateringInterval int
		setupMocks       func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantGroup        *entities.Group
		wantErr          bool
		checkResult      func(*testing.T, *entities.Group)
	}{
		{
			name:             "Success - interval=7, next watering in future",
			telegramID:       456,
			wateringInterval: 7,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.ConfirmAddGroup, temp.Step)

						var group entities.Group
						assert.NoError(t, json.Unmarshal(temp.Data, &group))
						assert.Equal(t, 7, group.WateringInterval)

						expectedNext := lastWatering.AddDate(0, 0, 7)
						assert.WithinDuration(t, expectedNext, group.NextWateringDate, time.Second)
						return nil
					}).
					Times(1)
			},
			wantGroup: &entities.Group{
				UserID:           123,
				Title:            "Цветы",
				Description:      "Уход за цветами",
				LastWateringDate: lastWatering,
				WateringInterval: 7,
				NextWateringDate: lastWatering.AddDate(0, 0, 7),
			},
			wantErr: false,
			checkResult: func(t *testing.T, got *entities.Group) {
				assert.Equal(t, 7, got.WateringInterval)
				assert.WithinDuration(t, got.NextWateringDate, lastWatering.AddDate(0, 0, 7), time.Second)
			},
		},
		{
			name:             "Success - interval=1, next watering in past → adjusted to today",
			telegramID:       456,
			wateringInterval: 1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				pastWatering := now.AddDate(0, 0, -10) // 10 дней назад
				group := &entities.Group{UserID: 123, LastWateringDate: pastWatering}
				data, _ := json.Marshal(group)

				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(&entities.Temporary{Data: data}, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						var group entities.Group
						assert.NoError(t, json.Unmarshal(temp.Data, &group))
						assert.WithinDuration(t, today, group.NextWateringDate, time.Second)
						return nil
					}).
					Times(1)
			},
			wantGroup: &entities.Group{
				UserID:           123,
				LastWateringDate: now.AddDate(0, 0, -10),
				WateringInterval: 1,
				NextWateringDate: today,
			},
			wantErr: false,
			checkResult: func(t *testing.T, got *entities.Group) {
				assert.WithinDuration(t, today, got.NextWateringDate, time.Second)
			},
		},
		{
			name:             "Failure - GetUserTemporary error",
			telegramID:       999,
			wateringInterval: 7,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:             "Failure - UpdateTemporary returns error",
			telegramID:       456,
			wateringInterval: 7,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotGroup, err := useCases.AddGroupWateringInterval(tt.telegramID, tt.wateringInterval)

			if tt.wantErr {
				assert.Nil(t, gotGroup)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotGroup)
				assert.Equal(t, tt.wantGroup.UserID, gotGroup.UserID)
				assert.Equal(t, tt.wantGroup.WateringInterval, gotGroup.WateringInterval)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, gotGroup)
			}
		})
	}
}

func TestTemporaryUseCases_ResetTemporary(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	temp := &entities.Temporary{ID: 1, UserID: 123, Data: []byte("some data"), MessageID: func() *int { i := 100; return &i }(), Step: 999}

	tests := []struct {
		name        string
		telegramID  int
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantErr     bool
		checkResult func(*testing.T, *entities.Temporary)
	}{
		{
			name:       "Success - temporary reset to defaults",
			telegramID: 456,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Nil(t, temp.Data)
						assert.Nil(t, temp.MessageID)
						assert.Equal(t, steps.Start, temp.Step)
						return nil
					}).
					Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, updatedTemp *entities.Temporary) {
				assert.Nil(t, updatedTemp.Data)
				assert.Nil(t, updatedTemp.MessageID)
				assert.Equal(t, steps.Start, updatedTemp.Step)
			},
		},
		{
			name:       "Failure - GetUserTemporary error",
			telegramID: 999,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - UpdateTemporary returns error",
			telegramID: 456,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			err := useCases.ResetTemporary(tt.telegramID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, temp)
			}
		})
	}
}

func TestTemporaryUseCases_AddPlantTitle(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	temp := &entities.Temporary{ID: 1, UserID: 123, Step: 0, MessageID: func() *int { i := 100; return &i }()}

	tests := []struct {
		name        string
		telegramID  int
		title       string
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantPlant   *entities.Plant
		wantErr     bool
		checkResult func(*testing.T, *entities.Temporary, *entities.Plant)
	}{
		{
			name:       "Success - plant title added and temporary updated",
			telegramID: 456,
			title:      "Кактус",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.AddPlantDescription, temp.Step)
						assert.Nil(t, temp.MessageID) // сброшен

						var plant entities.Plant
						assert.NoError(t, json.Unmarshal(temp.Data, &plant))
						assert.Equal(t, "Кактус", plant.Title)
						assert.Equal(t, 123, plant.UserID)
						return nil
					}).
					Times(1)
			},
			wantPlant: &entities.Plant{UserID: 123, Title: "Кактус"},
			wantErr:   false,
			checkResult: func(t *testing.T, updatedTemp *entities.Temporary, returnedPlant *entities.Plant) {
				assert.Equal(t, "Кактус", returnedPlant.Title)
				assert.Equal(t, 123, returnedPlant.UserID)
				assert.Equal(t, steps.AddPlantDescription, updatedTemp.Step)
				assert.Nil(t, updatedTemp.MessageID)
			},
		},
		{
			name:       "Failure - GetUserTemporary returns error",
			telegramID: 999,
			title:      "Кактус",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - UpdateTemporary returns error",
			telegramID: 456,
			title:      "Кактус",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotPlant, err := useCases.AddPlantTitle(tt.telegramID, tt.title)

			if tt.wantErr {
				assert.Nil(t, gotPlant)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotPlant)
				assert.Equal(t, tt.wantPlant.UserID, gotPlant.UserID)
				assert.Equal(t, tt.wantPlant.Title, gotPlant.Title)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, temp, gotPlant)
			}
		})
	}
}

func TestTemporaryUseCases_AddPlantDescription(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	existingPlant := &entities.Plant{UserID: 123, Title: "Кактус"}
	tempData, _ := json.Marshal(existingPlant)
	temp := &entities.Temporary{ID: 1, UserID: 123, Data: tempData, Step: steps.AddPlantDescription, MessageID: func() *int { i := 100; return &i }()}

	tests := []struct {
		name        string
		telegramID  int
		description string
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantPlant   *entities.Plant
		wantErr     bool
		checkResult func(*testing.T, *entities.Temporary, *entities.Plant)
	}{
		{
			name:        "Success - description added and temporary updated",
			telegramID:  456,
			description: "Маленький зелёный кактус на подоконнике",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.AddPlantGroup, temp.Step)
						assert.Nil(t, temp.MessageID)

						var plant entities.Plant
						assert.NoError(t, json.Unmarshal(temp.Data, &plant))
						assert.Equal(t, "Кактус", plant.Title)
						assert.Equal(t, "Маленький зелёный кактус на подоконнике", plant.Description)
						return nil
					}).
					Times(1)
			},
			wantPlant: &entities.Plant{
				UserID:      123,
				Title:       "Кактус",
				Description: "Маленький зелёный кактус на подоконнике",
			},
			wantErr: false,
			checkResult: func(t *testing.T, updatedTemp *entities.Temporary, returnedPlant *entities.Plant) {
				assert.Equal(t, "Маленький зелёный кактус на подоконнике", returnedPlant.Description)
				assert.Equal(t, steps.AddPlantGroup, updatedTemp.Step)
				assert.Nil(t, updatedTemp.MessageID)
			},
		},
		{
			name:        "Failure - GetUserTemporary error",
			telegramID:  999,
			description: "Описание",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:        "Failure - UpdateTemporary returns error",
			telegramID:  456,
			description: "Описание",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotPlant, err := useCases.AddPlantDescription(tt.telegramID, tt.description)

			if tt.wantErr {
				assert.Nil(t, gotPlant)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotPlant)
				assert.Equal(t, tt.wantPlant.Title, gotPlant.Title)
				assert.Equal(t, tt.wantPlant.Description, gotPlant.Description)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, temp, gotPlant)
			}
		})
	}
}

func TestTemporaryUseCases_AddPlantGroup(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	existingPlant := &entities.Plant{UserID: 123, Title: "Кактус"}
	tempData, _ := json.Marshal(existingPlant)
	temp := &entities.Temporary{ID: 1, UserID: 123, Data: tempData, Step: steps.AddPlantGroup}

	tests := []struct {
		name        string
		telegramID  int
		groupID     int
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantPlant   *entities.Plant
		wantErr     bool
		checkResult func(*testing.T, *entities.Plant)
	}{
		{
			name:       "Success - group assigned and plant not exists",
			telegramID: 456,
			groupID:    777,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Eq(entities.Plant{UserID: 123, Title: "Кактус", GroupID: 777})).
					Return(false, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.AddPlantPhotoQuestion, temp.Step)

						var plant entities.Plant
						assert.NoError(t, json.Unmarshal(temp.Data, &plant))
						assert.Equal(t, 777, plant.GroupID)
						assert.Equal(t, "Кактус", plant.Title)
						return nil
					}).
					Times(1)
			},
			wantPlant: &entities.Plant{UserID: 123, Title: "Кактус", GroupID: 777},
			wantErr:   false,
			checkResult: func(t *testing.T, got *entities.Plant) {
				assert.Equal(t, 777, got.GroupID)
				assert.Equal(t, "Кактус", got.Title)
			},
		},
		{
			name:       "Failure - GetUserTemporary error",
			telegramID: 999,
			groupID:    777,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - PlantExists returns error",
			telegramID: 456,
			groupID:    777,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Eq(entities.Plant{UserID: 123, Title: "Кактус", GroupID: 777})).
					Return(false, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						gomock.Any(),
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - plant already exists",
			telegramID: 456,
			groupID:    777,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Eq(entities.Plant{UserID: 123, Title: "Кактус", GroupID: 777})).
					Return(true, nil).
					Times(1)

				// Не должно быть UpdateTemporary
				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Times(0)
			},
			wantErr: true,
			checkResult: func(t *testing.T, got *entities.Plant) {
				assert.Nil(t, got)
			},
		},
		{
			name:       "Failure - UpdateTemporary returns error",
			telegramID: 456,
			groupID:    777,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Eq(entities.Plant{UserID: 123, Title: "Кактус", GroupID: 777})).
					Return(false, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotPlant, err := useCases.AddPlantGroup(tt.telegramID, tt.groupID)

			if tt.wantErr {
				assert.Nil(t, gotPlant)
				assert.Error(t, err)
				if tt.name == "Failure - plant already exists" {
					assert.True(t, errors.Is(err, customerrors.ErrPlantAlreadyExists))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotPlant)
				assert.Equal(t, tt.wantPlant.UserID, gotPlant.UserID)
				assert.Equal(t, tt.wantPlant.Title, gotPlant.Title)
				assert.Equal(t, tt.wantPlant.GroupID, gotPlant.GroupID)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, gotPlant)
			}
		})
	}
}

func TestTemporaryUseCases_AddPlantPhoto(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	existingPlant := &entities.Plant{UserID: 123, Title: "Кактус", GroupID: 777}
	tempData, _ := json.Marshal(existingPlant)
	temp := &entities.Temporary{ID: 1, UserID: 123, Data: tempData, Step: steps.AddPlantPhotoQuestion}

	photo := []byte("photo_data")

	tests := []struct {
		name        string
		telegramID  int
		photo       []byte
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantPlant   *entities.Plant
		wantErr     bool
		checkResult func(*testing.T, *entities.Plant)
	}{
		{
			name:       "Success - photo added and temporary updated",
			telegramID: 456,
			photo:      photo,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.ConfirmAddPlant, temp.Step)
						assert.Nil(t, temp.MessageID)

						var plant entities.Plant
						assert.NoError(t, json.Unmarshal(temp.Data, &plant))
						assert.Equal(t, photo, plant.Photo)
						return nil
					}).
					Times(1)
			},
			wantPlant: &entities.Plant{
				UserID:  123,
				Title:   "Кактус",
				GroupID: 777,
				Photo:   photo,
			},
			wantErr: false,
			checkResult: func(t *testing.T, got *entities.Plant) {
				assert.Equal(t, photo, got.Photo)
				assert.Equal(t, steps.ConfirmAddPlant, temp.Step)
				assert.Nil(t, temp.MessageID)
			},
		},
		{
			name:       "Failure - GetUserTemporary error",
			telegramID: 999,
			photo:      photo,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - UpdateTemporary returns error",
			telegramID: 456,
			photo:      photo,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			gotPlant, err := useCases.AddPlantPhoto(tt.telegramID, tt.photo)

			if tt.wantErr {
				assert.Nil(t, gotPlant)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotPlant)
				assert.Equal(t, tt.wantPlant.Title, gotPlant.Title)
				assert.Equal(t, tt.wantPlant.Description, gotPlant.Description)
			}

		})
	}
}

func TestTemporaryUseCases_ManagePlant(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	temp := &entities.Temporary{ID: 1, UserID: 123, Step: 0}

	tests := []struct {
		name        string
		telegramID  int
		plantID     int
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantErr     bool
		checkResult func(*testing.T, *entities.Temporary)
	}{
		{
			name:       "Success - plant ID stored and step updated",
			telegramID: 456,
			plantID:    888,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.ManagePlantAction, temp.Step)

						var plant entities.Plant
						assert.NoError(t, json.Unmarshal(temp.Data, &plant))
						assert.Equal(t, 888, plant.ID)
						return nil
					}).
					Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, updatedTemp *entities.Temporary) {
				var plant entities.Plant
				assert.NoError(t, json.Unmarshal(updatedTemp.Data, &plant))
				assert.Equal(t, 888, plant.ID)
				assert.Equal(t, steps.ManagePlantAction, updatedTemp.Step)
			},
		},
		{
			name:       "Failure - GetUserTemporary returns error",
			telegramID: 999,
			plantID:    888,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - UpdateTemporary returns error",
			telegramID: 456,
			plantID:    888,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			err := useCases.ManagePlant(tt.telegramID, tt.plantID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, temp)
			}
		})
	}
}

func TestTemporaryUseCases_ManageGroup(t *testing.T) {
	user := &entities.User{ID: 123, TelegramID: 456}
	temp := &entities.Temporary{ID: 1, UserID: 123, Step: 0}

	tests := []struct {
		name        string
		telegramID  int
		groupID     int
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantErr     bool
		checkResult func(*testing.T, *entities.Temporary)
	}{
		{
			name:       "Success - group ID stored and step updated",
			telegramID: 456,
			groupID:    777,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					DoAndReturn(func(temp entities.Temporary) error {
						assert.Equal(t, steps.ManageGroupAction, temp.Step)

						var group entities.Group
						assert.NoError(t, json.Unmarshal(temp.Data, &group))
						assert.Equal(t, 777, group.ID)
						return nil
					}).
					Times(1)
			},
			wantErr: false,
			checkResult: func(t *testing.T, updatedTemp *entities.Temporary) {
				var group entities.Group
				assert.NoError(t, json.Unmarshal(updatedTemp.Data, &group))
				assert.Equal(t, 777, group.ID)
				assert.Equal(t, steps.ManageGroupAction, updatedTemp.Step)
			},
		},
		{
			name:       "Failure - GetUserTemporary returns error",
			telegramID: 999,
			groupID:    777,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Temporary User with telegramID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:       "Failure - UpdateTemporary returns error",
			telegramID: 456,
			groupID:    777,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserByTelegramID(456).
					Return(user, nil).
					Times(1)

				storage.
					EXPECT().
					GetTemporaryByUserID(123).
					Return(temp, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateTemporary(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Temporary with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
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

			useCases := &temporaryUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			err := useCases.ManageGroup(tt.telegramID, tt.groupID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, temp)
			}
		})
	}
}
