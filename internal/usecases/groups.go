package usecases

import (
	"fmt"

	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

type groupsUseCases struct {
	storage interfaces.Storage
	logger  logging.Logger
}

func (u *groupsUseCases) GetUserGroups(userID int) ([]entities.Group, error) {
	groups, err := u.storage.GetUserGroups(userID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get groups for user with ID=%d", userID),
			"Error",
			err,
		)
	}

	return groups, err
}

func (u *groupsUseCases) CountUserGroups(userID int) (int, error) {
	count, err := u.storage.CountUserGroups(userID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to count groups for user with ID=%d", userID),
			"Error",
			err,
		)
	}

	return count, err
}
