package storage

import (
	"context"

	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/logging"

	sq "github.com/Masterminds/squirrel"

	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
)

const (
	plantsTableName  = "plants"
	idColumnName     = "id"
	userIDColumnName = "user_id"
)

type plantsStorage struct {
	dbConnector db.Connector
	logger      logging.Logger
}

func (s *plantsStorage) CreatePlant(plant entities.Plant) (int, error) {
	return 0, nil
}

func (s *plantsStorage) UpdatePlant(plant entities.Plant) error {
	return nil
}

func (s *plantsStorage) PlantExists(plant entities.Plant) (bool, error) {
	return false, nil
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
