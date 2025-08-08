package usecases

import (
	"fmt"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

type UseCases struct {
	storage interfaces.Storage
	logger  logging.Logger
}

func New(
	storage interfaces.Storage,
	logger logging.Logger,
) *UseCases {
	return &UseCases{
		storage: storage,
		logger:  logger,
	}
}

func (u *UseCases) SaveUser(user entities.User, messageID int) (int, error) {
	// Затенение, чтобы не переписывать реальный объект юзера, который нужно будет сохранить, если такого не существует:
	if user, err := u.storage.GetUserByTelegramID(user.TelegramID); err == nil {
		return user.ID, nil
	}

	userID, err := u.storage.SaveUser(user)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to save user with telegramID=%d", user.TelegramID),
			"Error",
			err,
		)

		return 0, err
	}

	temp := entities.Temporary{
		UserID:    userID,
		Step:      startStep,
		MessageID: messageID,
	}

	// TODO при проблемах логики следует сделать в рамках транзакции с сохранением пользователя
	if err = u.storage.CreateTemporary(temp); err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to save temporary info for user with ID=%d", userID),
			"Temp",
			temp,
			"Error",
			err,
		)

		return 0, err
	}

	return userID, nil
}

func (u *UseCases) GetUserGroups(userID int) ([]entities.Group, error) {
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

func (u *UseCases) CountUserGroups(userID int) (int, error) {
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

func (u *UseCases) GetUserPlants(userID int) ([]entities.Plant, error) {
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

func (u *UseCases) CountUserPlants(userID int) (int, error) {
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
