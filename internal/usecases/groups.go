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
			fmt.Sprintf("Failed to get Groups for User with ID=%d", userID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return groups, err
}

func (u *groupsUseCases) CountUserGroups(userID int) (int, error) {
	count, err := u.storage.CountUserGroups(userID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to count Groups for User with ID=%d", userID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return count, err
}

func (u *groupsUseCases) CreateGroup(group entities.Group) (*entities.Group, error) {
	groupID, err := u.storage.CreateGroup(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to create Group for User with ID=%d", group.UserID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	group.ID = groupID

	return &group, err
}

func (u *groupsUseCases) UpdateGroup(group entities.Group) error {
	err := u.storage.UpdateGroup(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Group for with ID=%d", group.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return err
}

func (u *groupsUseCases) GroupExists(group entities.Group) (bool, error) {
	exists, err := u.storage.GroupExists(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to check Group existence for User with ID=%d", group.UserID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return exists, err
}

func (u *groupsUseCases) GetGroup(id int) (*entities.Group, error) {
	group, err := u.storage.GetGroup(id)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get Group with ID=%d", id),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return group, err
}

func (u *groupsUseCases) GetGroupsForNotify(limit, offset int) ([]entities.Group, error) {
	groups, err := u.storage.GetGroupsForNotify(limit, offset)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get Groups for Notify with limit=%d and offset=%d", limit, offset),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return groups, err
}
