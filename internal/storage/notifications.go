package storage

import (
	"context"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
)

const (
	notificationsTableName = "notifications"
	textColumnName         = "text"
	sentAtColumnName       = "sent_at"
)

type notificationsStorage struct {
	dbConnector db.Connector
	logger      logging.Logger
}

func (s *notificationsStorage) SaveNotification(notification entities.Notification) (int, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Insert(notificationsTableName).
		Columns(
			groupIDColumnName,
			messageIDColumnName,
			textColumnName,
			sentAtColumnName,
		).
		Values(
			notification.GroupID,
			notification.MessageID,
			notification.Text,
			notification.SentAt,
		).
		Suffix(returningIDSuffix).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return 0, err
	}

	var notificationID int
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&notificationID); err != nil {
		return 0, err
	}

	return notificationID, nil
}
