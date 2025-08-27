package usecases

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	customerrors "github.com/DKhorkov/plantsCareTelegramBot/internal/errors"
	mockstorage "github.com/DKhorkov/plantsCareTelegramBot/mocks/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestPlantsUseCases_CountUserPlants(t *testing.T) {
	tests := []struct {
		name       string
		userID     int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantCount  int
		wantErr    bool
	}{
		{
			name:   "Success - count returned",
			userID: 123,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CountUserPlants(123).
					Return(5, nil).
					Times(1)
			},
			wantCount: 5,
			wantErr:   false,
		},
		{
			name:   "Failure - storage error",
			userID: 123,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CountUserPlants(123).
					Return(0, assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to count Plants for User with ID=123",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantCount: 0,
			wantErr:   true,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.CountUserPlants(tt.userID)

			assert.Equal(t, tt.wantCount, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestPlantsUseCases_GetGroupPlants(t *testing.T) {
	now := time.Now()
	expectedPlants := []entities.Plant{
		{
			ID:          1,
			GroupID:     10,
			UserID:      123,
			Title:       "Фикус",
			Description: "Комнатный",
			Photo:       []byte{0xFF, 0xD8},
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	tests := []struct {
		name       string
		groupID    int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantPlants []entities.Plant
		wantErr    bool
	}{
		{
			name:    "Success - plants returned",
			groupID: 10,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroupPlants(10).
					Return(expectedPlants, nil).
					Times(1)
			},
			wantPlants: expectedPlants,
			wantErr:    false,
		},
		{
			name:    "Failure - storage error",
			groupID: 10,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroupPlants(10).
					Return(nil, assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to get Plants for Group with ID=10",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlants: nil,
			wantErr:    true,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.GetGroupPlants(tt.groupID)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.wantPlants), len(got))
				assert.Equal(t, tt.wantPlants[0].ID, got[0].ID)
				assert.Equal(t, tt.wantPlants[0].Title, got[0].Title)
			}
		})
	}
}

func TestPlantsUseCases_CountGroupPlants(t *testing.T) {
	tests := []struct {
		name       string
		groupID    int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantCount  int
		wantErr    bool
	}{
		{
			name:    "Success - count returned",
			groupID: 10,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CountGroupPlants(10).
					Return(3, nil).
					Times(1)
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:    "Failure - storage error",
			groupID: 10,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CountGroupPlants(10).
					Return(0, assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to count Plants for Group with ID=10",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantCount: 0,
			wantErr:   true,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.CountGroupPlants(tt.groupID)

			assert.Equal(t, tt.wantCount, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestPlantsUseCases_CreatePlant(t *testing.T) {
	inputPlant := entities.Plant{
		GroupID:     10,
		UserID:      123,
		Title:       "Кактус",
		Description: "Солнечный подоконник",
		Photo:       []byte{0xFF, 0xD9},
	}

	tests := []struct {
		name       string
		plant      entities.Plant
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantPlant  *entities.Plant
		wantErr    bool
	}{
		{
			name:  "Success - plant created",
			plant: inputPlant,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CreatePlant(inputPlant).
					Return(42, nil).
					Times(1)
			},
			wantPlant: &entities.Plant{
				ID:          42,
				GroupID:     10,
				UserID:      123,
				Title:       "Кактус",
				Description: "Солнечный подоконник",
				Photo:       []byte{0xFF, 0xD9},
			},
			wantErr: false,
		},
		{
			name:  "Failure - storage error",
			plant: inputPlant,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CreatePlant(inputPlant).
					Return(0, assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to create Plant for Group with ID=10",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.CreatePlant(tt.plant)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantPlant.ID, got.ID)
				assert.Equal(t, tt.wantPlant.Title, got.Title)
				assert.Equal(t, tt.wantPlant.GroupID, got.GroupID)
				assert.Equal(t, tt.wantPlant.UserID, got.UserID)
			}
		})
	}
}

func TestPlantsUseCases_GetPlant(t *testing.T) {
	now := time.Now()
	expectedPlant := &entities.Plant{
		ID:          1,
		GroupID:     10,
		UserID:      123,
		Title:       "Фикус",
		Description: "Комнатный",
		Photo:       []byte{0xFF, 0xD8},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name       string
		plantID    int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantPlant  *entities.Plant
		wantErr    bool
	}{
		{
			name:    "Success - plant returned",
			plantID: 1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(expectedPlant, nil).
					Times(1)
			},
			wantPlant: expectedPlant,
			wantErr:   false,
		},
		{
			name:    "Failure - storage error",
			plantID: 1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(nil, assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to get Plant with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.GetPlant(tt.plantID)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPlant.ID, got.ID)
				assert.Equal(t, tt.wantPlant.Title, got.Title)
			}
		})
	}
}

func TestPlantsUseCases_DeletePlant(t *testing.T) {
	tests := []struct {
		name       string
		plantID    int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantErr    bool
	}{
		{
			name:    "Success - plant deleted",
			plantID: 1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					DeletePlant(1).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:    "Failure - storage error",
			plantID: 1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					DeletePlant(1).
					Return(assert.AnError).
					Times(1)
				logger.
					EXPECT().
					Error(
						"Failed to delete Plant with ID=1",
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			err := useCases.DeletePlant(tt.plantID)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestPlantsUseCases_UpdatePlantTitle(t *testing.T) {
	now := time.Now()
	existingPlant := &entities.Plant{
		ID:          1,
		GroupID:     10,
		UserID:      123,
		Title:       "Старый заголовок",
		Description: "Описание",
		Photo:       []byte{0xFF, 0xD8},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	updatedPlant := &entities.Plant{
		ID:          1,
		GroupID:     10,
		UserID:      123,
		Title:       "Новый заголовок",
		Description: "Описание",
		Photo:       []byte{0xFF, 0xD8},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name       string
		plantID    int
		title      string
		setupMocks func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantPlant  *entities.Plant
		wantErr    bool
		wantErrIs  error
	}{
		{
			name:    "Success - title updated",
			plantID: 1,
			title:   "Новый заголовок",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				// Получение растения
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				// Проверка на дубль — не существует
				storage.
					EXPECT().
					PlantExists(gomock.Any()).
					DoAndReturn(func(p entities.Plant) (bool, error) {
						assert.Equal(t, "Новый заголовок", p.Title)
						assert.Equal(t, 10, p.GroupID)
						return false, nil
					}).
					Times(1)

				// Обновление
				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					DoAndReturn(func(p entities.Plant) error {
						assert.Equal(t, "Новый заголовок", p.Title)
						assert.Equal(t, 1, p.ID)
						return nil
					}).
					Times(1)
			},
			wantPlant: updatedPlant,
			wantErr:   false,
		},
		{
			name:    "Failure - plant not found",
			plantID: 999,
			title:   "Новый",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Plant with ID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
		},
		{
			name:    "Failure - PlantExists returns error",
			plantID: 1,
			title:   "Новый",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Any()).
					Return(false, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to check existence for Plant with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
		},
		{
			name:    "Failure - plant already exists (duplicate)",
			plantID: 1,
			title:   "Дубль",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Any()).
					Return(true, nil).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
			wantErrIs: customerrors.ErrPlantAlreadyExists,
		},
		{
			name:    "Failure - UpdatePlant returns error",
			plantID: 1,
			title:   "Новый",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Any()).
					Return(false, nil).
					Times(1)

				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Plant with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.UpdatePlantTitle(tt.plantID, tt.title)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantPlant.ID, got.ID)
				assert.Equal(t, tt.wantPlant.Title, got.Title)
				assert.Equal(t, tt.wantPlant.Description, got.Description)
			}
		})
	}
}

func TestPlantsUseCases_UpdatePlantDescription(t *testing.T) {
	now := time.Now()
	existingPlant := &entities.Plant{
		ID:          1,
		GroupID:     10,
		UserID:      123,
		Title:       "Фикус",
		Description: "Старое описание",
		Photo:       []byte{0xFF, 0xD8},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	updatedPlant := &entities.Plant{
		ID:          1,
		GroupID:     10,
		UserID:      123,
		Title:       "Фикус",
		Description: "Новое описание",
		Photo:       []byte{0xFF, 0xD8},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name        string
		plantID     int
		description string
		setupMocks  func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantPlant   *entities.Plant
		wantErr     bool
	}{
		{
			name:        "Success - description updated",
			plantID:     1,
			description: "Новое описание",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					DoAndReturn(func(p entities.Plant) error {
						assert.Equal(t, "Новое описание", p.Description)
						return nil
					}).
					Times(1)
			},
			wantPlant: updatedPlant,
			wantErr:   false,
		},
		{
			name:        "Failure - plant not found",
			plantID:     999,
			description: "Новое",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(999).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Plant with ID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
		},
		{
			name:        "Failure - UpdatePlant returns error",
			plantID:     1,
			description: "Новое",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Plant with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.UpdatePlantDescription(tt.plantID, tt.description)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantPlant.ID, got.ID)
				assert.Equal(t, tt.wantPlant.Description, got.Description)
			}
		})
	}
}

func TestPlantsUseCases_UpdatePlantGroup(t *testing.T) {
	now := time.Now()
	existingPlant := &entities.Plant{
		ID:          1,
		GroupID:     10,
		UserID:      123,
		Title:       "Фикус",
		Description: "Комнатный",
		Photo:       []byte{0xFF, 0xD8},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name       string
		plantID    int
		groupID    int
		setupMocks func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantPlant  *entities.Plant
		wantErr    bool
		wantErrIs  error
	}{
		{
			name:    "Success - group updated",
			plantID: 1,
			groupID: 20,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Any()).
					DoAndReturn(func(p entities.Plant) (bool, error) {
						assert.Equal(t, 20, p.GroupID)
						return false, nil
					}).
					Times(1)

				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					DoAndReturn(func(p entities.Plant) error {
						assert.Equal(t, 20, p.GroupID)
						return nil
					}).
					Times(1)
			},
			wantPlant: &entities.Plant{
				ID:      1,
				GroupID: 20,
			},
			wantErr: false,
		},
		{
			name:    "Failure - update plant error",
			plantID: 1,
			groupID: 20,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Any()).
					DoAndReturn(func(p entities.Plant) (bool, error) {
						assert.Equal(t, 20, p.GroupID)
						return false, nil
					}).
					Times(1)

				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Plant with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
		{
			name:    "Failure - get plant error",
			plantID: 1,
			groupID: 20,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Plant with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
		},
		{
			name:    "Failure - plant already exists in new group",
			plantID: 1,
			groupID: 20,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Any()).
					Return(true, nil).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
			wantErrIs: customerrors.ErrPlantAlreadyExists,
		},
		{
			name:    "Failure - PlantExists error",
			plantID: 1,
			groupID: 20,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					PlantExists(gomock.Any()).
					Return(false, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to check existence for Plant with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPlant: nil,
			wantErr:   true,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.UpdatePlantGroup(tt.plantID, tt.groupID)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantPlant.ID, got.ID)
				assert.Equal(t, tt.wantPlant.GroupID, got.GroupID)
			}
		})
	}
}

func TestPlantsUseCases_UpdatePlantPhoto(t *testing.T) {
	now := time.Now()
	existingPlant := &entities.Plant{
		ID:          1,
		GroupID:     10,
		UserID:      123,
		Title:       "Фикус",
		Description: "Комнатное растение",
		Photo:       []byte{0xFF, 0xD8}, // старое фото
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	newPhoto := []byte{0xFF, 0xD9} // новое фото

	tests := []struct {
		name       string
		plantID    int
		photo      []byte
		setupMocks func(*mockstorage.MockStorage, *mocklogging.MockLogger)
		wantPhoto  []byte
		wantErr    bool
	}{
		{
			name:    "Success - photo updated successfully",
			plantID: 1,
			photo:   newPhoto,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				// Получение растения
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				// Проверка, что в UpdatePlant передаётся объект с новым фото
				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					DoAndReturn(func(p entities.Plant) error {
						assert.Equal(t, 1, p.ID)
						assert.Equal(t, newPhoto, p.Photo)
						assert.Equal(t, "Фикус", p.Title)
						return nil
					}).
					Times(1)
			},
			wantPhoto: newPhoto,
			wantErr:   false,
		},
		{
			name:    "Failure - plant not found",
			plantID: 999,
			photo:   newPhoto,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				// Ошибка при получении
				storage.
					EXPECT().
					GetPlant(999).
					Return(nil, assert.AnError).
					Times(1)

				// Логирование ошибки
				logger.
					EXPECT().
					Error(
						"Failed to get Plant with ID=999",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPhoto: nil,
			wantErr:   true,
		},
		{
			name:    "Failure - UpdatePlant returns error",
			plantID: 1,
			photo:   newPhoto,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Plant with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantPhoto: nil,
			wantErr:   true,
		},
		{
			name:    "Success - photo set to nil (clear photo)",
			plantID: 1,
			photo:   nil,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetPlant(1).
					Return(existingPlant, nil).
					Times(1)

				storage.
					EXPECT().
					UpdatePlant(gomock.Any()).
					DoAndReturn(func(p entities.Plant) error {
						assert.Nil(t, p.Photo)
						return nil
					}).
					Times(1)
			},
			wantPhoto: nil,
			wantErr:   false,
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

			useCases := &plantsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.UpdatePlantPhoto(tt.plantID, tt.photo)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantPhoto, got.Photo)
				assert.Equal(t, 1, got.ID)
			}
		})
	}
}
