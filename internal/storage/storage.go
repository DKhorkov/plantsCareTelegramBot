package storage

import (
	"context"
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	sq "github.com/Masterminds/squirrel"
)

const (
	selectAllColumns       = "*"
	selectCount            = "COUNT(*)"
	usersTableName         = "users"
	plantsTableName        = "plants"
	groupsTableName        = "groups"
	notificationsTableName = "notifications"
	temporaryTableName     = "temporary"
	idColumnName           = "id"
	userIDColumnName       = "user_id"
	stepColumnName         = "step"
	messageIDColumnName    = "message_id"
	dataColumnName         = "data"
	telegramIDColumnName   = "telegram_id"
	usernameColumnName     = "username"
	firstnameColumnName    = "firstname"
	lastnameColumnName     = "lastname"
	isBotColumnName        = "is_bot"
	returningIDSuffix      = "RETURNING id"
	createdAtColumnName    = "created_at"
	updatedAtColumnName    = "updated_at"
	desc                   = "DESC"
	asc                    = "ASC"
)

type Storage struct {
	dbConnector db.Connector
	logger      logging.Logger
}

func New(
	dbConnector db.Connector,
	logger logging.Logger,
) *Storage {
	return &Storage{
		dbConnector: dbConnector,
		logger:      logger,
	}
}

func (s *Storage) SaveUser(user entities.User) (int, error) {
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

func (s *Storage) GetUserByTelegramID(telegramID int) (*entities.User, error) {
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

func (s *Storage) CreateTemporary(temp entities.Temporary) error {
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

func (s *Storage) UpdateTemporary(temp entities.Temporary) error {
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
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	_, err = connection.ExecContext(ctx, stmt, params...)
	return err
}

func (s *Storage) GetTemporaryByUserID(userID int) (*entities.Temporary, error) {
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

func (s *Storage) CreateGroup(group entities.Group) (int, error) {
	return 0, nil
}

func (s *Storage) UpdateGroup(group entities.Group) error {
	return nil
}

func (s *Storage) GroupExists(group entities.Group) (bool, error) {
	return false, nil
}

func (s *Storage) DeleteGroup(id int) error {
	return nil
}

func (s *Storage) GetUserGroups(userID int) ([]entities.Group, error) {
	return []entities.Group{}, nil
}

func (s *Storage) CountUserGroups(userID int) (int, error) {
	ctx := context.Background()
	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectCount).
		From(groupsTableName).
		Where(sq.Eq{userIDColumnName: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Storage) GetGroup(id int) (*entities.Group, error) {
	return nil, nil
}

func (s *Storage) CreatePlant(plant entities.Plant) (int, error) {
	return 0, nil
}

func (s *Storage) UpdatePlant(plant entities.Plant) error {
	return nil
}

func (s *Storage) PlantExists(plant entities.Plant) (bool, error) {
	return false, nil
}

func (s *Storage) DeletePlant(id int) error {
	return nil
}

func (s *Storage) GetUserPlants(userID int) ([]entities.Plant, error) {
	return []entities.Plant{}, nil
}

func (s *Storage) CountUserPlants(userID int) (int, error) {
	ctx := context.Background()
	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectCount).
		From(plantsTableName).
		Where(sq.Eq{userIDColumnName: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, err
	}

	var count int
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *Storage) GetPlant(id int) (*entities.Plant, error) {
	return nil, nil
}
