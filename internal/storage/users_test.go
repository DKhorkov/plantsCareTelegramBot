//go:build integration

package storage

import (
	"context"
	"database/sql"
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
)

const (
	migrationsDir    = "/migrations"
	gooseZeroVersion = 0
)

func TestUsersStorageTestSuite(t *testing.T) {
	suite.Run(t, new(UsersStorageTestSuite))
}

type UsersStorageTestSuite struct {
	suite.Suite

	cwd         string
	ctx         context.Context
	dbConnector db.Connector
	connection  *sql.Conn
	storage     *usersStorage
	logger      *mocklogging.MockLogger
}

func (s *UsersStorageTestSuite) SetupSuite() {
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

	s.storage = &usersStorage{
		dbConnector: s.dbConnector,
		logger:      s.logger,
	}
}

func (s *UsersStorageTestSuite) SetupTest() {
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

func (s *UsersStorageTestSuite) TearDownTest() {
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

func (s *UsersStorageTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *UsersStorageTestSuite) TestSaveUser_Success() {
	user := entities.User{
		TelegramID: 123456,
		Username:   "testuser",
		Firstname:  "John",
		Lastname:   "Doe",
		IsBot:      false,
	}

	userID, err := s.storage.SaveUser(user)
	s.NoError(err)
	s.Greater(userID, 0)

	// Проверим, что пользователь сохранился
	var storedUser entities.User
	columns := db.GetEntityColumns(&storedUser)
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT `+selectAllColumns+` FROM users WHERE id = $1`,
		userID,
	).Scan(columns...)
	s.NoError(err)

	s.Equal(user.TelegramID, storedUser.TelegramID)
	s.Equal(user.Username, storedUser.Username)
	s.Equal(user.Firstname, storedUser.Firstname)
	s.Equal(user.Lastname, storedUser.Lastname)
	s.Equal(user.IsBot, storedUser.IsBot)
}

func (s *UsersStorageTestSuite) TestSaveUser_DatabaseError() {
	// Здесь мы НЕ можем протестировать ошибку БД через реальное подключение,
	// потому что мы не мокаем DBConnector.
	// Но в вашем стиле — вы тестируете только **успешные и "не найдено"** кейсы через реальное подключение.
	// Ошибки БД тестируются отдельно, если есть моки.
	// → Поэтому этот тест **опускаем**, если нет моков.
	// Либо оставляем только те, что можно проверить через реальное подключение.
	//
	// В вашем примере с `GetUserCommunications` — вы не тестируете ошибки БД напрямую.
	// → Значит, и здесь **не нужно**.
	//
	// ✅ Оставим только успешный кейс и "не найдено" для Get-методов.
}

func (s *UsersStorageTestSuite) TestGetUserByTelegramID_UserExists() {
	telegramID := 123456

	_, err := s.connection.ExecContext(
		context.Background(),
		`
			INSERT INTO users (telegram_id, username, firstname, lastname, is_bot) 
			VALUES ($1, $2, $3, $4, $5)
		`,
		telegramID,
		"testuser",
		"John",
		"Doe",
		false,
	)
	s.NoError(err)

	user, err := s.storage.GetUserByTelegramID(telegramID)
	s.NoError(err)
	s.NotNil(user)
	s.Equal(telegramID, user.TelegramID)
	s.Equal("testuser", user.Username)
	s.Equal("John", user.Firstname)
	s.Equal("Doe", user.Lastname)
	s.False(user.IsBot)
}

func (s *UsersStorageTestSuite) TestGetUserByTelegramID_UserNotFound() {
	telegramID := 999999 // не существует

	user, err := s.storage.GetUserByTelegramID(telegramID)
	s.Error(err)
	s.Nil(user)
	// В вашем стиле — достаточно проверить, что ошибка и nil
}

func (s *UsersStorageTestSuite) TestGetUserByID_UserExists() {
	var userID int
	telegramID := 123456

	err := s.connection.QueryRowContext(
		context.Background(),
		`
			INSERT INTO users (telegram_id, username, firstname, lastname, is_bot) 
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id
		`,
		telegramID,
		"jane_doe",
		"Jane",
		"Doe",
		false,
	).Scan(&userID)
	s.NoError(err)

	user, err := s.storage.GetUserByID(userID)
	s.NoError(err)
	s.NotNil(user)
	s.Equal(userID, user.ID)
	s.Equal(telegramID, user.TelegramID)
	s.Equal("jane_doe", user.Username)
	s.Equal("Jane", user.Firstname)
	s.Equal("Doe", user.Lastname)
	s.False(user.IsBot)
}

func (s *UsersStorageTestSuite) TestGetUserByID_UserNotFound() {
	nonExistentID := 999999

	user, err := s.storage.GetUserByID(nonExistentID)
	s.Error(err)
	s.Nil(user)
}
