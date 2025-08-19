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
			fmt.Sprintf("Failed to get Temporary User with telegramID=%d", telegramID),
			"Error",
			err,
		)

		return nil, err
	}

	temp, err := u.storage.GetTemporaryByUserID(user.ID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get Temporary for User with ID=%d", user.ID),
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
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
			"Error",
			err,
		)

		return err
	}

	return nil
}

func (u *temporaryUseCases) SetTemporaryMessage(telegramID int, messageID *int) error {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return err
	}

	temp.MessageID = messageID
	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
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
			fmt.Sprintf("Failed to marshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.AddGroupDescription
	temp.MessageID = nil // not to delete already deleted message

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
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
			fmt.Sprintf("Failed to unmarshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	group.Description = description

	data, err := json.Marshal(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.AddGroupLastWateringDate
	temp.MessageID = nil // not to delete already deleted message

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
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
			fmt.Sprintf("Failed to unmarshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	group.LastWateringDate = lastWateringDate

	data, err := json.Marshal(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.AddGroupWateringInterval

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return group, nil
}

func (u *temporaryUseCases) AddGroupWateringInterval(telegramID, wateringInterval int) (*entities.Group, error) {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return nil, err
	}

	group := &entities.Group{}
	if err = json.Unmarshal(temp.Data, group); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to unmarshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	nextWateringDate := group.LastWateringDate.AddDate(0, 0, wateringInterval)

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, nextWateringDate.Location())

	if nextWateringDate.Before(today) {
		nextWateringDate = today
	}

	group.WateringInterval = wateringInterval
	group.NextWateringDate = nextWateringDate

	data, err := json.Marshal(group)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.ConfirmAddGroup

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return group, nil
}

func (u *temporaryUseCases) ResetTemporary(telegramID int) error {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return err
	}

	temp.Data = nil
	temp.MessageID = nil
	temp.Step = steps.Start

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
			"Error",
			err,
		)

		return err
	}

	return nil
}

func (u *temporaryUseCases) AddPlantTitle(telegramID int, title string) (*entities.Plant, error) {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return nil, err
	}

	plant := &entities.Plant{
		UserID: temp.UserID,
		Title:  title,
	}

	data, err := json.Marshal(plant)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.AddPlantDescription
	temp.MessageID = nil // not to delete already deleted message

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return plant, nil
}

func (u *temporaryUseCases) AddPlantDescription(telegramID int, description string) (*entities.Plant, error) {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return nil, err
	}

	plant := &entities.Plant{}
	if err = json.Unmarshal(temp.Data, plant); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to unmarshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	plant.Description = description

	data, err := json.Marshal(plant)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.AddPlantGroup
	temp.MessageID = nil // not to delete already deleted message

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return plant, nil
}

func (u *temporaryUseCases) AddPlantGroup(telegramID, groupID int) (*entities.Plant, error) {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return nil, err
	}

	plant := &entities.Plant{}
	if err = json.Unmarshal(temp.Data, plant); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to unmarshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	plant.GroupID = groupID

	data, err := json.Marshal(plant)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.AddPlantPhotoQuestion

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return plant, nil
}

func (u *temporaryUseCases) AddPlantPhoto(telegramID int, photo []byte) (*entities.Plant, error) {
	temp, err := u.GetUserTemporary(telegramID)
	if err != nil {
		return nil, err
	}

	plant := &entities.Plant{}
	if err = json.Unmarshal(temp.Data, plant); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to unmarshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	plant.Photo = photo

	data, err := json.Marshal(plant)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to marshal data for User with ID=%d", temp.UserID),
			"Error",
			err,
		)

		return nil, err
	}

	temp.Data = data
	temp.Step = steps.ConfirmAddPlant
	temp.MessageID = nil // not to delete already deleted message

	if err = u.storage.UpdateTemporary(*temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to update Temporary with ID=%d", temp.ID),
			"Error",
			err,
		)

		return nil, err
	}

	return plant, nil
}
