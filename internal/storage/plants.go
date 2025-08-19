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

func (s *plantsStorage) GetPlant(id int) (*entities.Plant, error) {
	return nil, nil
}
