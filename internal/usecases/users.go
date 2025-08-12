package usecases

import (
	"fmt"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/steps"
)

type usersUseCases struct {
	storage interfaces.Storage
	logger  logging.Logger
}

func (u *usersUseCases) SaveUser(user entities.User) (int, error) {
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
		UserID: userID,
		Step:   steps.StartStep,
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

func (u *usersUseCases) GetUserByTelegramID(telegramID int) (*entities.User, error) {
	user, err := u.storage.GetUserByTelegramID(telegramID)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to get user with telegramID=%d", telegramID),
			"Error",
			err,
		)

		return nil, err
	}

	return user, nil
}
