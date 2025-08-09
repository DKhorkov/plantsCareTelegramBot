package usecases

import (
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

type UseCases struct {
	usersUseCases
	groupsUseCases
	plantsUseCases
	temporaryUseCases
}

func New(
	storage interfaces.Storage,
	logger logging.Logger,
) *UseCases {
	return &UseCases{
		usersUseCases: usersUseCases{
			storage: storage,
			logger:  logger,
		},
		groupsUseCases: groupsUseCases{
			storage: storage,
			logger:  logger,
		},
		plantsUseCases: plantsUseCases{
			storage: storage,
			logger:  logger,
		},
		temporaryUseCases: temporaryUseCases{
			storage: storage,
			logger:  logger,
		},
	}
}
