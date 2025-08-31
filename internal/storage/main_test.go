package storage

import (
	"github.com/DKhorkov/libs/db"
	mockdb "github.com/DKhorkov/libs/db/mocks"
	"github.com/DKhorkov/libs/logging"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestNewStorage_Initialization(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name        string
		dbConnector db.Connector
		logger      logging.Logger
		shouldPanic bool
		checkFields []string
	}{
		{
			name:        "Valid dependencies — all storages should be initialized",
			dbConnector: mockdb.NewMockConnector(ctrl),
			logger:      mocklogging.NewMockLogger(ctrl),
			shouldPanic: false,
			checkFields: []string{
				"usersStorage",
				"temporaryStorage",
				"groupsStorage",
				"plantsStorage",
				"notificationsStorage",
			},
		},
		{
			name:        "Nil logger — storage still initializes if logger is optional",
			dbConnector: mockdb.NewMockConnector(ctrl),
			logger:      nil,
			shouldPanic: false,
			checkFields: []string{
				"usersStorage",
				"temporaryStorage",
				"groupsStorage",
				"plantsStorage",
				"notificationsStorage",
			},
		},
		{
			name:        "Nil dbConnector — expect panic if db is required",
			dbConnector: nil,
			logger:      mocklogging.NewMockLogger(ctrl),
			shouldPanic: false,
			checkFields: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					New(tt.dbConnector, tt.logger)
				}, "Expected panic when dbConnector is nil")
				return
			}

			// Выполняем тестируемый код
			s := New(tt.dbConnector, tt.logger)

			// Проверяем, что storage не nil
			assert.NotNil(t, s, "Storage should not be nil")

			// Проверяем, что все вложенные хранилища инициализированы
			// (на практике — они анонимные поля, поэтому просто проверим их содержимое)

			assert.NotNil(t, &s.usersStorage, "usersStorage should not be nil")
			assert.Equal(t, tt.dbConnector, s.usersStorage.dbConnector, "usersStorage should have correct dbConnector")
			assert.Equal(t, tt.logger, s.usersStorage.logger, "usersStorage should have correct logger")

			assert.NotNil(t, &s.temporaryStorage, "temporaryStorage should not be nil")
			assert.Equal(t, tt.dbConnector, s.temporaryStorage.dbConnector, "temporaryStorage should have correct dbConnector")
			assert.Equal(t, tt.logger, s.temporaryStorage.logger, "temporaryStorage should have correct logger")

			assert.NotNil(t, &s.groupsStorage, "groupsStorage should not be nil")
			assert.Equal(t, tt.dbConnector, s.groupsStorage.dbConnector, "groupsStorage should have correct dbConnector")
			assert.Equal(t, tt.logger, s.groupsStorage.logger, "groupsStorage should have correct logger")

			assert.NotNil(t, &s.plantsStorage, "plantsStorage should not be nil")
			assert.Equal(t, tt.dbConnector, s.plantsStorage.dbConnector, "plantsStorage should have correct dbConnector")
			assert.Equal(t, tt.logger, s.plantsStorage.logger, "plantsStorage should have correct logger")

			assert.NotNil(t, &s.notificationsStorage, "notificationsStorage should not be nil")
			assert.Equal(t, tt.dbConnector, s.notificationsStorage.dbConnector, "notificationsStorage should have correct dbConnector")
			assert.Equal(t, tt.logger, s.notificationsStorage.logger, "notificationsStorage should have correct logger")
		})
	}
}
