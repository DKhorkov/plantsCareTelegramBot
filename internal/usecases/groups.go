package usecases

import (
	"fmt"
	"time"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	customerrors "github.com/DKhorkov/plantsCareTelegramBot/internal/errors"
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

		return nil, err
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

		return 0, err
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

		return nil, err
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

func (u *groupsUseCases) GetGroup(id int) (*entities.Group, error) {
	group, err := u.storage.GetGroup(id)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get Group with ID=%d", id),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)

		return nil, err
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

		return nil, err
	}

	return groups, err
}

func (u *groupsUseCases) DeleteGroup(id int) error {
	err := u.storage.DeleteGroup(id)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to delete Group with ID=%d", id),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)
	}

	return err
}

func (u *groupsUseCases) UpdateGroupTitle(id int, title string) (*entities.Group, error) {
	group, err := u.GetGroup(id)
	if err != nil {
		return nil, err
	}

	group.Title = title

	exists, err := u.storage.GroupExists(*group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to check existence for Group with ID=%d", group.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)

		return nil, err
	}

	if exists {
		return nil, customerrors.ErrGroupAlreadyExists
	}

	if err = u.storage.UpdateGroup(*group); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Group with ID=%d", group.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)

		return nil, err
	}

	return group, err
}

func (u *groupsUseCases) UpdateGroupDescription(id int, description string) (*entities.Group, error) {
	group, err := u.GetGroup(id)
	if err != nil {
		return nil, err
	}

	group.Description = description
	if err = u.storage.UpdateGroup(*group); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Group with ID=%d", group.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)

		return nil, err
	}

	return group, err
}

func (u *groupsUseCases) UpdateGroupLastWateringDate(id int, lastWateringDate time.Time) (*entities.Group, error) {
	group, err := u.GetGroup(id)
	if err != nil {
		return nil, err
	}

	nextWateringDate := lastWateringDate.AddDate(0, 0, group.WateringInterval)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, nextWateringDate.Location())

	if nextWateringDate.Before(today) {
		nextWateringDate = today
	}

	group.LastWateringDate = lastWateringDate

	group.NextWateringDate = nextWateringDate

	if err = u.storage.UpdateGroup(*group); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Group with ID=%d", group.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)

		return nil, err
	}

	return group, err
}

func (u *groupsUseCases) UpdateGroupWateringInterval(id, wateringInterval int) (*entities.Group, error) {
	group, err := u.GetGroup(id)
	if err != nil {
		return nil, err
	}

	nextWateringDate := group.LastWateringDate.AddDate(0, 0, wateringInterval)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, nextWateringDate.Location())

	if nextWateringDate.Before(today) {
		nextWateringDate = today
	}

	group.WateringInterval = wateringInterval

	group.NextWateringDate = nextWateringDate

	if err = u.storage.UpdateGroup(*group); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Group with ID=%d", group.ID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)

		return nil, err
	}

	return group, err
}
