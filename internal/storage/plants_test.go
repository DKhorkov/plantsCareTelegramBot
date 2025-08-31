//go:build integration

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DKhorkov/libs/db"
	"github.com/DKhorkov/libs/loadenv"
	"github.com/DKhorkov/libs/logging"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/config"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"os"
	"path"
	"testing"
	"time"
)

func TestPlantsStorageTestSuite(t *testing.T) {
	suite.Run(t, new(PlantsStorageTestSuite))
}

type PlantsStorageTestSuite struct {
	suite.Suite

	cwd         string
	ctx         context.Context
	dbConnector db.Connector
	connection  *sql.Conn
	storage     *plantsStorage
	logger      *mocklogging.MockLogger
}

func (s *PlantsStorageTestSuite) SetupSuite() {
	// Инициализируем переменные окружения для дальнейшего считывания:
	loadenv.Init("../../.env")

	cwd, err := os.Getwd()
	s.NoError(err)

	cfg := config.New()
	logger := logging.New(
		cfg.Logging.Level,
		cfg.Logging.LogFilePath,
	)

	dbConnector, err := db.New(
		db.BuildDsn(cfg.Database),
		cfg.Database.Driver,
		logger,
		db.WithMaxOpenConnections(cfg.Database.Pool.MaxOpenConnections),
		db.WithMaxIdleConnections(cfg.Database.Pool.MaxIdleConnections),
		db.WithMaxConnectionLifetime(cfg.Database.Pool.MaxConnectionLifetime),
		db.WithMaxConnectionIdleTime(cfg.Database.Pool.MaxConnectionIdleTime),
	)
	s.NoError(err, "failed to connect to database")

	s.NoError(goose.SetDialect(cfg.Database.Driver))

	ctrl := gomock.NewController(s.T())
	s.logger = mocklogging.NewMockLogger(ctrl)

	s.cwd = cwd
	s.dbConnector = dbConnector
	s.ctx = context.Background()

	s.storage = &plantsStorage{
		dbConnector: s.dbConnector,
		logger:      s.logger,
	}
}

func (s *PlantsStorageTestSuite) SetupTest() {
	s.NoError(
		goose.Up(
			s.dbConnector.Pool(),
			path.Dir(
				path.Dir(s.cwd),
			)+migrationsDir,
		),
	)

	connection, err := s.dbConnector.Connection(s.ctx)
	s.NoError(err)

	s.connection = connection
}

func (s *PlantsStorageTestSuite) TearDownTest() {
	s.NoError(
		goose.DownTo(
			s.dbConnector.Pool(),
			path.Dir(
				path.Dir(s.cwd),
			)+migrationsDir,
			gooseZeroVersion,
		),
	)

	s.NoError(s.connection.Close())
}

func (s *PlantsStorageTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *PlantsStorageTestSuite) createUser(now time.Time, offset int) int {
	var userID int
	err := s.connection.QueryRowContext(
		context.Background(),
		`
			INSERT INTO users (
				telegram_id, username, firstname, lastname, is_bot, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $6)
			RETURNING id
		`,
		123456789+int64(offset),
		fmt.Sprintf("user%d", offset),
		fmt.Sprintf("First%d", offset),
		fmt.Sprintf("Last%d", offset),
		false,
		now,
	).Scan(&userID)
	s.NoError(err)
	return userID
}

func (s *PlantsStorageTestSuite) createGroupForUser(userID int, now time.Time, offset int) int {
	var groupID int
	err := s.connection.QueryRowContext(
		context.Background(),
		`
			INSERT INTO groups (
				user_id, title, watering_interval, created_at, updated_at
			) VALUES ($1, $2, $3, $4, $4)
			RETURNING id
		`,
		userID,
		fmt.Sprintf("Группа %d", offset),
		7+offset,
		now,
	).Scan(&groupID)
	s.NoError(err)
	return groupID
}

func (s *PlantsStorageTestSuite) createPlantForGroup(groupID, userID int, title string, now time.Time, offset int) int {
	var plantID int
	err := s.connection.QueryRowContext(
		context.Background(),
		`
			INSERT INTO plants (
				group_id, 
				user_id, 
				title, 
				description, 
				photo, 
				created_at, 
				updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $6)
			RETURNING id
		`,
		groupID,
		userID,
		title,
		fmt.Sprintf("Описание растения %d", offset), // description
		[]byte{0xFF, 0xD8, byte(offset % 256)},      // photo: минимальный JPEG-подобный слайс
		now,
	).Scan(&plantID)
	s.NoError(err, "Не удалось создать растение для группы")
	return plantID
}

func (s *PlantsStorageTestSuite) TestCreatePlant_Success() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	plant := entities.Plant{
		GroupID:     groupID,
		UserID:      userID,
		Title:       "Фикус",
		Description: "Зелёный, любит свет",
		Photo:       []byte{0xFF, 0xD8, 0xFF}, // JPEG-заголовок
	}

	plantID, err := s.storage.CreatePlant(plant)
	s.NoError(err)
	s.Greater(plantID, 0)

	var stored entities.Plant
	columns := db.GetEntityColumns(&stored)
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT `+selectAllColumns+` FROM plants WHERE id = $1`,
		plantID,
	).Scan(columns...)
	s.NoError(err)

	s.Equal(plantID, stored.ID)
	s.Equal(groupID, stored.GroupID)
	s.Equal(userID, stored.UserID)
	s.Equal("Фикус", stored.Title)
	s.Equal("Зелёный, любит свет", stored.Description)
	s.Equal([]byte{0xFF, 0xD8, 0xFF}, stored.Photo)
	s.WithinDuration(now, stored.CreatedAt, 2*time.Second)
	s.WithinDuration(now, stored.UpdatedAt, 2*time.Second)
}

func (s *PlantsStorageTestSuite) TestUpdatePlant_Success() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)
	plantID := s.createPlantForGroup(groupID, userID, "Старое имя", now, 1)

	updatedPlant := entities.Plant{
		ID:          plantID,
		GroupID:     groupID,
		UserID:      userID,
		Title:       "Обновлённое имя",
		Description: "Новое описание",
		Photo:       []byte{0x00, 0x01, 0x02},
	}

	err := s.storage.UpdatePlant(updatedPlant)
	s.NoError(err)

	var stored entities.Plant
	columns := db.GetEntityColumns(&stored)
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT `+selectAllColumns+` FROM plants WHERE id = $1`,
		plantID,
	).Scan(columns...)
	s.NoError(err)

	s.Equal("Обновлённое имя", stored.Title)
	s.Equal("Новое описание", stored.Description)
	s.Equal([]byte{0x00, 0x01, 0x02}, stored.Photo)
	s.True(stored.UpdatedAt.After(now))
}

func (s *PlantsStorageTestSuite) TestUpdatePlant_NonExistent() {
	nonExistent := entities.Plant{
		ID:      999999,
		GroupID: 1,
		UserID:  1,
		Title:   "Не существует",
	}

	err := s.storage.UpdatePlant(nonExistent)
	s.NoError(err) // UPDATE 0 строк — не ошибка
}

func (s *PlantsStorageTestSuite) TestPlantExists_Exists() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)
	s.createPlantForGroup(groupID, userID, "Кактус", now, 1)

	exists, err := s.storage.PlantExists(entities.Plant{
		GroupID: groupID,
		Title:   "Кактус",
	})
	s.NoError(err)
	s.True(exists)
}

func (s *PlantsStorageTestSuite) TestPlantExists_NotExists() {
	exists, err := s.storage.PlantExists(entities.Plant{
		GroupID: 999,
		Title:   "Не существующее растение",
	})
	s.NoError(err)
	s.False(exists)
}

func (s *PlantsStorageTestSuite) TestPlantExists_DifferentGroupSameTitle() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID1 := s.createGroupForUser(userID, now, 1)
	groupID2 := s.createGroupForUser(userID, now, 2)
	s.createPlantForGroup(groupID1, userID, "Одинаковое имя", now, 1)

	exists, err := s.storage.PlantExists(entities.Plant{
		GroupID: groupID2,
		Title:   "Одинаковое имя",
	})
	s.NoError(err)
	s.False(exists)
}

func (s *PlantsStorageTestSuite) TestDeletePlant_Success() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)
	plantID := s.createPlantForGroup(groupID, userID, "Для удаления", now, 1)

	err := s.storage.DeletePlant(plantID)
	s.NoError(err)

	var count int
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT COUNT(*) FROM plants WHERE id = $1`,
		plantID,
	).Scan(&count)
	s.NoError(err)
	s.Equal(0, count)
}

func (s *PlantsStorageTestSuite) TestDeletePlant_NonExistent() {
	err := s.storage.DeletePlant(999999)
	s.NoError(err) // Удаление несуществующей — не ошибка
}

func (s *PlantsStorageTestSuite) TestCountUserPlants_HasPlants() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID1 := s.createGroupForUser(userID, now, 1)
	groupID2 := s.createGroupForUser(userID, now, 2)

	s.createPlantForGroup(groupID1, userID, "Растение 1", now, 1)
	s.createPlantForGroup(groupID1, userID, "Растение 2", now, 2)
	s.createPlantForGroup(groupID2, userID, "Растение 3", now, 3)

	count, err := s.storage.CountUserPlants(userID)
	s.NoError(err)
	s.Equal(3, count)
}

func (s *PlantsStorageTestSuite) TestCountUserPlants_NoPlants() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)

	count, err := s.storage.CountUserPlants(userID)
	s.NoError(err)
	s.Equal(0, count)
}

func (s *PlantsStorageTestSuite) TestGetGroupPlants_HasPlants() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	s.createPlantForGroup(groupID, userID, "Спатифиллум", now, 1)
	s.createPlantForGroup(groupID, userID, "Замиокулькас", now, 2)

	plants, err := s.storage.GetGroupPlants(groupID)
	s.NoError(err)
	s.Len(plants, 2)
	s.Equal("Спатифиллум", plants[0].Title)
	s.Equal("Замиокулькас", plants[1].Title)
}

func (s *PlantsStorageTestSuite) TestGetGroupPlants_NoPlants() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	plants, err := s.storage.GetGroupPlants(groupID)
	s.NoError(err)
	s.Empty(plants)
}

func (s *PlantsStorageTestSuite) TestCountGroupPlants_HasPlants() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	s.createPlantForGroup(groupID, userID, "A", now, 1)
	s.createPlantForGroup(groupID, userID, "B", now, 2)
	s.createPlantForGroup(groupID, userID, "C", now, 3)

	count, err := s.storage.CountGroupPlants(groupID)
	s.NoError(err)
	s.Equal(3, count)
}

func (s *PlantsStorageTestSuite) TestCountGroupPlants_NoPlants() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	count, err := s.storage.CountGroupPlants(groupID)
	s.NoError(err)
	s.Equal(0, count)
}

func (s *PlantsStorageTestSuite) TestGetPlant_Exists() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)
	plantID := s.createPlantForGroup(groupID, userID, "Монстера", now, 1)

	plant, err := s.storage.GetPlant(plantID)
	s.NoError(err)
	s.NotNil(plant)
	s.Equal(plantID, plant.ID)
	s.Equal("Монстера", plant.Title)
}

func (s *PlantsStorageTestSuite) TestGetPlant_NotFound() {
	plant, err := s.storage.GetPlant(999999)
	s.Error(err)
	s.Nil(plant)
}
