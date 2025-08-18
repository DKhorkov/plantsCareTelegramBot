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

func (u *groupsUseCases) CreateGroup(group entities.Group) (*entities.Group, error) {
	groupID, err := u.storage.CreateGroup(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to create group for user with ID=%d", group.UserID),
			"Error",
			err,
		)
	}

	group.ID = groupID

	return &group, err
}

func (u *groupsUseCases) GroupExists(group entities.Group) (bool, error) {
	exists, err := u.storage.GroupExists(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to check Group existence for user with ID=%d", group.UserID),
			"Error",
			err,
		)
	}

	return exists, err
}

func (u *groupsUseCases) GetGroup(id int) (*entities.Group, error) {
	group, err := u.storage.GetGroup(id)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get group with ID=%d", id),
			"Error",
			err,
		)
	}

	return group, err
}
