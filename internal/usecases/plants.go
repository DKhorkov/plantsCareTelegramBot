package usecases

import (
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

type plantsUseCases struct {
	storage interfaces.Storage
	logger  logging.Logger
}

func (u *plantsUseCases) GetUserPlants(userID int) ([]entities.Plant, error) {
	plants, err := u.storage.GetUserPlants(userID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get plants for user with ID=%d", userID),
			"Error",
			err,
		)
	}

	return plants, err
}

func (u *plantsUseCases) CountUserPlants(userID int) (int, error) {
	count, err := u.storage.CountUserPlants(userID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to count plants for user with ID=%d", userID),
			"Error",
			err,
		)
	}

	return count, err
}
