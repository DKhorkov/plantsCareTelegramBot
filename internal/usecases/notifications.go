package usecases

import (
	"fmt"

	"github.com/DKhorkov/libs/logging"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
)

type notificationsUseCases struct {
	storage interfaces.Storage
	logger  logging.Logger
}

func (u *notificationsUseCases) SaveNotification(notification entities.Notification) (*entities.Notification, error) {
	notificationID, err := u.storage.SaveNotification(notification)
	if err != nil {
		u.logger.Error(
			fmt.Sprintf("Failed to save Notification for Group with ID=%d", notification.GroupID),
			"Error", err,
			"Tracing", logging.GetLogTraceback(loggingTraceSkipLevel),
		)

		return nil, err
	}

	notification.ID = notificationID

	return &notification, err
}
