package storage

import (
	"context"
	"time"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
)

const (
	temporaryTableName  = "temporary"
	stepColumnName      = "step"
	messageIDColumnName = "message_id"
	dataColumnName      = "data"
	updatedAtColumnName = "updated_at"
)

type temporaryStorage struct {
	dbConnector db.Connector
	logger      logging.Logger
}

func (s *temporaryStorage) CreateTemporary(temp entities.Temporary) error {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Insert(temporaryTableName).
		Columns(
			userIDColumnName,
			stepColumnName,
			messageIDColumnName,
			dataColumnName,
		).
		Values(
			temp.UserID,
			temp.Step,
			temp.MessageID,
			temp.Data,
		).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	_, err = connection.ExecContext(ctx, stmt, params...)

	return err
}

func (s *temporaryStorage) UpdateTemporary(temp entities.Temporary) error {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Update(temporaryTableName).
		Where(sq.Eq{idColumnName: temp.ID}).
		Set(stepColumnName, temp.Step).
		Set(messageIDColumnName, temp.MessageID).
		Set(dataColumnName, temp.Data).
		Set(updatedAtColumnName, time.Now()).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	_, err = connection.ExecContext(ctx, stmt, params...)

	return err
}

func (s *temporaryStorage) GetTemporaryByUserID(userID int) (*entities.Temporary, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(temporaryTableName).
		Where(sq.Eq{userIDColumnName: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	temporary := &entities.Temporary{}

	columns := db.GetEntityColumns(temporary)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return temporary, nil
}
