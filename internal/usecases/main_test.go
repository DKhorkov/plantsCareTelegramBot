package usecases

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	mockstorage "github.com/DKhorkov/plantsCareTelegramBot/mocks/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestNew_UseCasesInitialization(t *testing.T) {
	// Настройка моков
	ctrl := gomock.NewController(t)

	mockStorage := mockstorage.NewMockStorage(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	// Табличный тест
	tests := []struct {
		name     string
		storage  *mockstorage.MockStorage
		logger   *mocklogging.MockLogger
		validate func(*testing.T, *UseCases)
	}{
		{
			name:    "All dependencies provided - should initialize all use cases",
			storage: mockStorage,
			logger:  mockLogger,
			validate: func(t *testing.T, uc *UseCases) {
				assert.NotNil(t, uc, "UseCases should not be nil")

				// Проверяем, что все вложенные use cases инициализированы
				assert.NotNil(t, uc.usersUseCases.storage, "usersUseCases.storage should be set")
				assert.NotNil(t, uc.usersUseCases.logger, "usersUseCases.logger should be set")

				assert.NotNil(t, uc.groupsUseCases.storage, "groupsUseCases.storage should be set")
				assert.NotNil(t, uc.groupsUseCases.logger, "groupsUseCases.logger should be set")

				assert.NotNil(t, uc.plantsUseCases.storage, "plantsUseCases.storage should be set")
				assert.NotNil(t, uc.plantsUseCases.logger, "plantsUseCases.logger should be set")

				assert.NotNil(t, uc.temporaryUseCases.storage, "temporaryUseCases.storage should be set")
				assert.NotNil(t, uc.temporaryUseCases.logger, "temporaryUseCases.logger should be set")

				assert.NotNil(t, uc.notificationsUseCases.storage, "notificationsUseCases.storage should be set")
				assert.NotNil(t, uc.notificationsUseCases.logger, "notificationsUseCases.logger should be set")

				// Проверяем, что зависимости переданы те же
				assert.Same(t, mockStorage, uc.usersUseCases.storage, "Storage should be the same instance")
				assert.Same(t, mockLogger, uc.usersUseCases.logger, "Logger should be the same instance")
			},
		},
		{
			name:    "Nil storage - should still initialize (no validation in New)",
			storage: nil,
			logger:  mockLogger,
			validate: func(t *testing.T, uc *UseCases) {
				assert.NotNil(t, uc, "UseCases should still be created even with nil storage")
				assert.Nil(t, uc.usersUseCases.storage, "storage should be nil as passed")
				assert.NotNil(t, uc.usersUseCases.logger, "logger should be set")
			},
		},
		{
			name:    "Nil logger - should still initialize",
			storage: mockStorage,
			logger:  nil,
			validate: func(t *testing.T, uc *UseCases) {
				assert.NotNil(t, uc, "UseCases should still be created")
				assert.Nil(t, uc.usersUseCases.logger, "logger should be nil as passed")
				assert.NotNil(t, uc.usersUseCases.storage, "storage should be set")
			},
		},
	}

	// Запуск всех тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := New(tt.storage, tt.logger)
			tt.validate(t, uc)
		})
	}
}
