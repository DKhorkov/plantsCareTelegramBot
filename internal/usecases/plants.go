package usecases

import (
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	customerrors "github.com/DKhorkov/plantsCareTelegramBot/internal/errors"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

type plantsUseCases struct {
	storage interfaces.Storage
	logger  logging.Logger
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

func (u *plantsUseCases) GetPlant(id int) (*entities.Plant, error) {
	plant, err := u.storage.GetPlant(id)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get Plant with ID=%d", id),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return plant, err
}

func (u *plantsUseCases) DeletePlant(id int) error {
	err := u.storage.DeletePlant(id)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to delete Plant with ID=%d", id),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return err
}

func (u *plantsUseCases) UpdatePlantTitle(id int, title string) (*entities.Plant, error) {
	plant, err := u.GetPlant(id)
	if err != nil {
		return nil, err
	}

	plant.Title = title

	exists, err := u.storage.PlantExists(*plant)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to check existence for Plant with ID=%d", plant.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	if exists {
		return nil, customerrors.ErrPlantAlreadyExists
	}

	if err = u.storage.UpdatePlant(*plant); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Plant with ID=%d", plant.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return plant, err
}

func (u *plantsUseCases) UpdatePlantDescription(id int, description string) (*entities.Plant, error) {
	plant, err := u.GetPlant(id)
	if err != nil {
		return nil, err
	}

	plant.Description = description
	if err = u.storage.UpdatePlant(*plant); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Plant with ID=%d", plant.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return plant, err
}

func (u *plantsUseCases) UpdatePlantGroup(id, groupID int) (*entities.Plant, error) {
	plant, err := u.GetPlant(id)
	if err != nil {
		return nil, err
	}

	plant.GroupID = groupID

	exists, err := u.storage.PlantExists(*plant)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to check existence for Plant with ID=%d", plant.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	if exists {
		return nil, customerrors.ErrPlantAlreadyExists
	}

	if err = u.storage.UpdatePlant(*plant); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Plant with ID=%d", plant.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return plant, err
}

func (u *plantsUseCases) UpdatePlantPhoto(id int, photo []byte) (*entities.Plant, error) {
	plant, err := u.GetPlant(id)
	if err != nil {
		return nil, err
	}

	plant.Photo = photo
	if err = u.storage.UpdatePlant(*plant); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Plant with ID=%d", plant.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return plant, err
}
