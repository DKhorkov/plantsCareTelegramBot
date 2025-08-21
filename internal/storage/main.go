package storage

import (
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
)

type Storage struct {
	usersStorage
	temporaryStorage
	groupsStorage
	plantsStorage
	notificationsStorage
}

func New(
	dbConnector db.Connector,
	logger logging.Logger,
) *Storage {
	return &Storage{
		usersStorage: usersStorage{
			dbConnector: dbConnector,
			logger:      logger,
		},
		temporaryStorage: temporaryStorage{
			dbConnector: dbConnector,
			logger:      logger,
		},
		groupsStorage: groupsStorage{
			dbConnector: dbConnector,
			logger:      logger,
		},
		plantsStorage: plantsStorage{
			dbConnector: dbConnector,
			logger:      logger,
		},
		notificationsStorage: notificationsStorage{
			dbConnector: dbConnector,
			logger:      logger,
		},
	}
}
