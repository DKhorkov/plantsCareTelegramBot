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
			fmt.Sprintf("Failed to get Plants for User with ID=%d", userID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return plants, err
}

func (u *plantsUseCases) CountUserPlants(userID int) (int, error) {
	count, err := u.storage.CountUserPlants(userID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to count Plants for User with ID=%d", userID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return count, err
}

func (u *plantsUseCases) GetGroupPlants(groupID int) ([]entities.Plant, error) {
	plants, err := u.storage.GetGroupPlants(groupID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get Plants for Group with ID=%d", groupID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return plants, err
}

func (u *plantsUseCases) CountGroupPlants(groupID int) (int, error) {
	count, err := u.storage.CountGroupPlants(groupID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to count Plants for Group with ID=%d", groupID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return count, err
}

func (u *plantsUseCases) CreatePlant(plant entities.Plant) (*entities.Plant, error) {
	plantID, err := u.storage.CreatePlant(plant)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to create Plant for Group with ID=%d", plant.GroupID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	plant.ID = plantID

	return &plant, err
}

func (u *plantsUseCases) PlantExists(plant entities.Plant) (bool, error) {
	exists, err := u.storage.PlantExists(plant)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to check Plant existence for Group with ID=%d", plant.GroupID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return exists, err
}
