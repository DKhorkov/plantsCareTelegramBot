package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
)

const (
	plantsTableName   = "plants"
	groupIDColumnName = "group_id"
	photoColumnName   = "photo"
)

type plantsStorage struct {
	dbConnector db.Connector
	logger      logging.Logger
}

func (s *plantsStorage) CreatePlant(plant entities.Plant) (int, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Insert(plantsTableName).
		Columns(
			groupIDColumnName,
			userIDColumnName,
			titleColumnName,
			descriptionColumnName,
			photoColumnName,
		).
		Values(
			plant.GroupID,
			plant.UserID,
			plant.Title,
			plant.Description,
			plant.Photo,
		).
		Suffix(returningIDSuffix).
		PlaceholderFormat(sq.Dollar). // pq postgres driver works only with $ placeholders
		ToSql()
	if err != nil {
		return 0, err
	}

	var plantID int
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(&plantID); err != nil {
		return 0, err
	}

	return plantID, nil
}

func (s *plantsStorage) UpdatePlant(plant entities.Plant) error {
	return nil
}

func (s *plantsStorage) PlantExists(plant entities.Plant) (bool, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return false, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectExists).
		From(plantsTableName).
		Where(
			sq.Eq{
				groupIDColumnName:     plant.GroupID,
				titleColumnName:       plant.Title,
				descriptionColumnName: plant.Description,
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

func (s *plantsStorage) DeletePlant(id int) error {
	return nil
}

func (s *plantsStorage) GetUserPlants(userID int) ([]entities.Plant, error) {
	return []entities.Plant{}, nil
}

func (s *plantsStorage) CountUserPlants(userID int) (int, error) {
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

func (s *plantsStorage) GetGroupPlants(groupID int) ([]entities.Plant, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(plantsTableName).
		Where(sq.Eq{groupIDColumnName: groupID}).
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

	var plants []entities.Plant

	for rows.Next() {
		plant := entities.Plant{}
		columns := db.GetEntityColumns(&plant) // Only pointer to use rows.Scan() successfully

		if err = rows.Scan(columns...); err != nil {
			return nil, err
		}

		plants = append(plants, plant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return plants, nil
}

func (s *plantsStorage) CountGroupPlants(groupID int) (int, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return 0, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectCount).
		From(plantsTableName).
		Where(sq.Eq{groupIDColumnName: groupID}).
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

func (s *plantsStorage) GetPlant(id int) (*entities.Plant, error) {
	ctx := context.Background()

	connection, err := s.dbConnector.Connection(ctx)
	if err != nil {
		return nil, err
	}

	defer db.CloseConnectionContext(ctx, connection, s.logger)

	stmt, params, err := sq.
		Select(selectAllColumns).
		From(plantsTableName).
		Where(sq.Eq{idColumnName: id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	plant := &entities.Plant{}

	columns := db.GetEntityColumns(plant)
	if err = connection.QueryRowContext(ctx, stmt, params...).Scan(columns...); err != nil {
		return nil, err
	}

	return plant, nil
}
