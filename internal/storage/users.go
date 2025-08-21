package storage

import (
	"context"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
)

const (
	selectAllColumns     = "*"
	usersTableName       = "users"
	telegramIDColumnName = "telegram_id"
	usernameColumnName   = "username"
	firstnameColumnName  = "firstname"
	lastnameColumnName   = "lastname"
	isBotColumnName      = "is_bot"
	returningIDSuffix    = "RETURNING id"
)

type usersStorage struct {
	dbConnector db.Connector
	logger      logging.Logger
}

func (s *usersStorage) SaveUser(user entities.User) (int, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Insert(usersTableName).
		Columns(
			telegramIDColumnName,
			usernameColumnName,
			firstnameColumnName,
			lastnameColumnName,
			isBotColumnName,
		).
		Values(
			user.TelegramID,
			user.Username,
			user.Firstname,
			user.Lastname,
			user.IsBot,
		).
		Suffix(returningIDSuffix).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return 0, err
	}

	var userID int
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&userID); err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *usersStorage) GetUserByTelegramID(telegramID int) (*entities.User, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(usersTableName).
		Where(sq.Eq{telegramIDColumnName: telegramID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &entities.User{}

	columns := db.GetEntityColumns(user)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *usersStorage) GetUserByID(id int) (*entities.User, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(usersTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	user := &entities.User{}

	columns := db.GetEntityColumns(user)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return user, nil
}
