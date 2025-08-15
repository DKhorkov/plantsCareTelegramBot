package usecases

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
)

type temporaryUseCases struct {
	storage interfaces.Storage
	logger  logging.Logger
}

func (u *temporaryUseCases) GetUserTemporary(telegramID int) (*entities.Temporary, error) {
	user, err := u.storage.GetUserByTelegramID(telegramID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get user with telegramID=%d", telegramID),
			"Error",
			err,
		)

		return nil, err
	}

	temp, err := u.storage.GetTemporaryByUserID(user.ID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get Temporary data for user with ID=%d", user.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return temp, nil
}

func (u *temporaryUseCases) SetTemporaryStep(telegramID, step int) error {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return err
	}

	temp.Step = step
	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update temporary data with ID=%d", temp.ID),
			"Error",
			err,
		)

		return err
	}

	return nil
}

func (u *temporaryUseCases) SetTemporaryMessage(telegramID, messageID int) error {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return err
	}

	temp.MessageID = &messageID
	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update temporary data with ID=%d", temp.ID),
			"Error",
			err,
		)

		return err
	}

	return nil
}

func (u *temporaryUseCases) AddGroupTitle(telegramID int, title string) (*entities.Group, error) {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return nil, err
	}

	group := &entities.Group{
		UserID: temp.UserID,
		Title:  title,
	}

	data, err := json.Marshal(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for user with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.MessageID = nil // not to delete already deleted message

	temp.Step = steps.AddGroupDescriptionStep

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update temporary data with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return group, nil
}

func (u *temporaryUseCases) AddGroupDescription(telegramID int, description string) (*entities.Group, error) {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return nil, err
	}

	group := &entities.Group{}
	if err = json.Unmarshal(temp.Data, group); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to unmarshal data for user with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	group.Description = description

	data, err := json.Marshal(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for user with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.MessageID = nil // not to delete already deleted message

	temp.Step = steps.AddGroupLastWateringDateStep

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update temporary data with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return group, nil
}

func (u *temporaryUseCases) AddGroupLastWateringDate(
	telegramID int,
	lastWateringDate time.Time,
) (*entities.Group, error) {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return nil, err
	}

	group := &entities.Group{}
	if err = json.Unmarshal(temp.Data, group); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to unmarshal data for user with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	group.LastWateringDate = lastWateringDate

	data, err := json.Marshal(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for user with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.AddGroupWateringIntervalStep

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update temporary data with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return group, nil
}
