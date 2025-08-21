package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
)

const (
	selectCount                = "COUNT(*)"
	groupsTableName            = "groups"
	idColumnName               = "id"
	userIDColumnName           = "user_id"
	titleColumnName            = "title"
	descriptionColumnName      = "description"
	lastWateringDateColumnName = "last_watering_date"
	nextWateringDateColumnName = "next_watering_date"
	wateringIntervalColumnName = "watering_interval"
	selectExists               = "1"
)

type groupsStorage struct {
	dbConnector db.Connector
	logger      logging.Logger
}

func (s *groupsStorage) CreateGroup(group entities.Group) (int, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Insert(groupsTableName).
		Columns(
			userIDColumnName,
			titleColumnName,
			descriptionColumnName,
			lastWateringDateColumnName,
			nextWateringDateColumnName,
			wateringIntervalColumnName,
		).
		Values(
			group.UserID,
			group.Title,
			group.Description,
			group.LastWateringDate,
			group.NextWateringDate,
			group.WateringInterval,
		).
		Suffix(returningIDSuffix).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return 0, err
	}

	var groupID int
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&groupID); err != nil {
		return 0, err
	}

	return groupID, nil
}

func (s *groupsStorage) UpdateGroup(group entities.Group) error {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Update(groupsTableName).
		Where(sq.Eq{idColumnName: group.ID}).
		Set(userIDColumnName, group.UserID).
		Set(titleColumnName, group.Title).
		Set(lastWateringDateColumnName, group.LastWateringDate).
		Set(wateringIntervalColumnName, group.WateringInterval).
		Set(nextWateringDateColumnName, group.NextWateringDate).
		Set(updatedAtColumnName, time.Now()).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return err
	}

	_, err = connection.ExecContext(ctx, stmt, params...)

	return err
}

func (s *groupsStorage) GroupExists(group entities.Group) (bool, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return false, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectExists).
		From(groupsTableName).
		Where(
			sq.Eq{
				userIDColumnName: group.UserID,
				titleColumnName:  group.Title,
				// descriptionColumnName: group.Description,
			},
		).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, err
	}

	var exists bool
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // Запись не найдена
		}

		return false, err
	}

	return exists, nil
}

func (s *groupsStorage) DeleteGroup(id int) error {
	return nil
}

func (s *groupsStorage) GetUserGroups(userID int) ([]entities.Group, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(groupsTableName).
		Where(sq.Eq{userIDColumnName: userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := connection.QueryContext(
		ctx,
		stmt,
		params...,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			logging.LogErrorContext(
				ctx,
				s.logger,
				"error during closing SQL rows",
				err,
			)
		}
	}()

	var groups []entities.Group

	for rows.Next() {
		group := entities.Group{}
		columns := db.GetEntityColumns(&group) // Only pointer to use rows.Scan() successfully

		if err = rows.Scan(columns...); err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *groupsStorage) CountUserGroups(userID int) (int, error) {
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

func (s *groupsStorage) GetGroup(id int) (*entities.Group, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(groupsTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	group := &entities.Group{}

	columns := db.GetEntityColumns(group)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupsStorage) GetGroupsForNotify(limit, offset int) ([]entities.Group, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(groupsTableName).
		Where(
			sq.Expr(
				nextWateringDateColumnName + " < CURRENT_TIMESTAMP",
			),
		).
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := connection.QueryContext(
		ctx,
		stmt,
		params...,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			logging.LogErrorContext(
				ctx,
				s.logger,
				"error during closing SQL rows",
				err,
			)
		}
	}()

	var groups []entities.Group

	for rows.Next() {
		group := entities.Group{}
		columns := db.GetEntityColumns(&group) // Only pointer to use rows.Scan() successfully

		if err = rows.Scan(columns...); err != nil {
			return nil, err
		}

		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}
