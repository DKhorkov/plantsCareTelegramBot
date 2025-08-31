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

func TestTemporaryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(TemporaryStorageTestSuite))
}

type TemporaryStorageTestSuite struct {
	suite.Suite

	cwd         string
	ctx         context.Context
	dbConnector db.Connector
	connection  *sql.Conn
	storage     *temporaryStorage
	logger      *mocklogging.MockLogger
}

func (s *TemporaryStorageTestSuite) SetupSuite() {
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

	s.storage = &temporaryStorage{
		dbConnector: s.dbConnector,
		logger:      s.logger,
	}
}

func (s *TemporaryStorageTestSuite) SetupTest() {
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

func (s *TemporaryStorageTestSuite) TearDownTest() {
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

func (s *TemporaryStorageTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *TemporaryStorageTestSuite) createUser(now time.Time, offset int) int {
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

func (s *TemporaryStorageTestSuite) TestCreateTemporary_Success_WithMessageID() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)

	messageID := 1001
	temp := entities.Temporary{
		UserID:    userID,
		Step:      5,
		MessageID: &messageID,
		Data:      []byte(`{"state": "in_progress", "plant": "Фикус"}`),
	}

	err := s.storage.CreateTemporary(temp)
	s.NoError(err)

	// Проверим, что запись появилась
	var stored entities.Temporary
	columns := db.GetEntityColumns(&stored)
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT `+selectAllColumns+` FROM temporary WHERE user_id = $1`,
		userID,
	).Scan(columns...)
	s.NoError(err)

	s.Equal(userID, stored.UserID)
	s.Equal(5, stored.Step)
	s.NotNil(stored.MessageID)
	s.Equal(1001, *stored.MessageID)
	s.Equal(`{"state": "in_progress", "plant": "Фикус"}`, string(stored.Data))
}

func (s *TemporaryStorageTestSuite) TestCreateTemporary_Success_WithoutMessageID() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)

	temp := entities.Temporary{
		UserID:    userID,
		Step:      3,
		MessageID: nil, // message_id = NULL
		Data:      []byte(`{"action": "rename"}`),
	}

	err := s.storage.CreateTemporary(temp)
	s.NoError(err)

	var stored entities.Temporary
	columns := db.GetEntityColumns(&stored)
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT `+selectAllColumns+` FROM temporary WHERE user_id = $1`,
		userID,
	).Scan(columns...)
	s.NoError(err)

	s.Equal(userID, stored.UserID)
	s.Equal(3, stored.Step)
	s.Nil(stored.MessageID) // Должно быть nil
	s.Equal(`{"action": "rename"}`, string(stored.Data))
}

func (s *TemporaryStorageTestSuite) TestUpdateTemporary_Success() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)

	// Сначала вставим запись вручную
	var tempID int
	err := s.connection.QueryRowContext(
		context.Background(),
		`
			INSERT INTO temporary (user_id, step, message_id, data)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`,
		userID,
		1,
		int64(999),
		[]byte(`{"initial": true}`),
	).Scan(&tempID)
	s.NoError(err)

	// Теперь обновим
	newMessageID := 2002
	updatedTemp := entities.Temporary{
		ID:        tempID,
		UserID:    userID,
		Step:      7,
		MessageID: &newMessageID,
		Data:      []byte(`{"updated": true, "step": 7}`),
	}

	err = s.storage.UpdateTemporary(updatedTemp)
	s.NoError(err)

	// Проверим
	var stored entities.Temporary
	columns := db.GetEntityColumns(&stored)
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT `+selectAllColumns+` FROM temporary WHERE id = $1`,
		tempID,
	).Scan(columns...)
	s.NoError(err)

	s.Equal(7, stored.Step)
	s.NotNil(stored.MessageID)
	s.Equal(2002, *stored.MessageID)
	s.Equal(`{"updated": true, "step": 7}`, string(stored.Data))
}

func (s *TemporaryStorageTestSuite) TestUpdateTemporary_NonExistent() {
	nonExistent := entities.Temporary{
		ID:     999999,
		UserID: 1,
		Step:   1,
		Data:   []byte("test"),
	}

	err := s.storage.UpdateTemporary(nonExistent)
	s.NoError(err) // UPDATE 0 строк — не ошибка
}

func (s *TemporaryStorageTestSuite) TestGetTemporaryByUserID_Exists() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)

	var tempID int
	err := s.connection.QueryRowContext(
		context.Background(),
		`
			INSERT INTO temporary (user_id, step, message_id, data)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`,
		userID,
		4,
		int64(5005),
		[]byte(`{"session": "active"}`),
	).Scan(&tempID)
	s.NoError(err)

	temp, err := s.storage.GetTemporaryByUserID(userID)
	s.NoError(err)
	s.NotNil(temp)
	s.Equal(tempID, temp.ID)
	s.Equal(4, temp.Step)
	s.NotNil(temp.MessageID)
	s.Equal(5005, *temp.MessageID)
	s.Equal(`{"session": "active"}`, string(temp.Data))
}

func (s *TemporaryStorageTestSuite) TestGetTemporaryByUserID_UniqueTest() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)

	// Вставим две записи
	_, err := s.connection.ExecContext(
		context.Background(),
		`
			INSERT INTO temporary (user_id, step, data)
			VALUES ($1, 1, 'data1'), ($1, 2, 'data2')
		`,
		userID,
	)
	s.Error(err)
}
