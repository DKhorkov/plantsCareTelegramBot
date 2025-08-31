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
	"sync"
	"testing"
	"time"
)

func TestNotificationsStorageTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationsStorageTestSuite))
}

type NotificationsStorageTestSuite struct {
	suite.Suite

	cwd         string
	ctx         context.Context
	dbConnector db.Connector
	connection  *sql.Conn
	storage     *notificationsStorage
	logger      *mocklogging.MockLogger
}

func (s *NotificationsStorageTestSuite) SetupSuite() {
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

	s.storage = &notificationsStorage{
		dbConnector: s.dbConnector,
		logger:      s.logger,
	}
}

func (s *NotificationsStorageTestSuite) SetupTest() {
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

func (s *NotificationsStorageTestSuite) TearDownTest() {
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

func (s *NotificationsStorageTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *NotificationsStorageTestSuite) createUser(now time.Time, offset int) int {
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

func (s *NotificationsStorageTestSuite) createGroupForUser(userID int, now time.Time, offset int) int {
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

func (s *NotificationsStorageTestSuite) TestSaveNotification_Success() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	notification := entities.Notification{
		GroupID:   groupID,
		MessageID: 1001,
		Text:      "Напоминание: пора поливать растения!",
		SentAt:    now,
	}

	notificationID, err := s.storage.SaveNotification(notification)
	s.NoError(err)
	s.Greater(notificationID, 0)

	// Проверим, что уведомление сохранилось
	var stored entities.Notification
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT id, group_id, message_id, text, sent_at FROM notifications WHERE id = $1`,
		notificationID,
	).Scan(&stored.ID, &stored.GroupID, &stored.MessageID, &stored.Text, &stored.SentAt)
	s.NoError(err)

	s.Equal(notificationID, stored.ID)
	s.Equal(groupID, stored.GroupID)
	s.Equal(1001, stored.MessageID)
	s.Equal("Напоминание: пора поливать растения!", stored.Text)
	s.WithinDuration(now, stored.SentAt, time.Second)
}

func (s *NotificationsStorageTestSuite) TestSaveNotification_GroupDoesNotExist() {
	now := time.Now().UTC()

	notification := entities.Notification{
		GroupID:   999999, // Такой группы нет
		MessageID: 1001,
		Text:      "Тестовое уведомление",
		SentAt:    now,
	}

	_, err := s.storage.SaveNotification(notification)
	s.Error(err)
	// Ожидаем ошибку от PostgreSQL: "insert or update on table "notifications" violates foreign key constraint"
	s.Contains(err.Error(), "violates foreign key constraint")
}

func (s *NotificationsStorageTestSuite) TestSaveNotification_EmptyText() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	notification := entities.Notification{
		GroupID:   groupID,
		MessageID: 1001,
		Text:      "", // Пустой текст — поле NOT NULL
		SentAt:    now,
	}

	_, err := s.storage.SaveNotification(notification)
	s.NoError(err)
}

func (s *NotificationsStorageTestSuite) TestSaveNotification_MessageID_Zero() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	notification := entities.Notification{
		GroupID:   groupID,
		MessageID: 0, // Допустимо, если Telegram может использовать 0
		Text:      "Сообщение с message_id = 0",
		SentAt:    now,
	}

	notificationID, err := s.storage.SaveNotification(notification)
	s.NoError(err)
	s.Greater(notificationID, 0)

	var stored entities.Notification
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT message_id FROM notifications WHERE id = $1`,
		notificationID,
	).Scan(&stored.MessageID)
	s.NoError(err)
	s.Equal(0, stored.MessageID)
}

func (s *NotificationsStorageTestSuite) TestSaveNotification_SentAt_InFuture() {
	now := time.Now().UTC()
	future := now.Add(1 * time.Hour)
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	notification := entities.Notification{
		GroupID:   groupID,
		MessageID: 1001,
		Text:      "Уведомление из будущего",
		SentAt:    future,
	}

	notificationID, err := s.storage.SaveNotification(notification)
	s.NoError(err)

	var stored entities.Notification
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT sent_at FROM notifications WHERE id = $1`,
		notificationID,
	).Scan(&stored.SentAt)
	s.NoError(err)
	s.WithinDuration(future, stored.SentAt, time.Second)
}

func (s *NotificationsStorageTestSuite) TestSaveNotification_ConcurrentInserts() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	var wg sync.WaitGroup
	errs := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			notification := entities.Notification{
				GroupID:   groupID,
				MessageID: 1000 + i,
				Text:      fmt.Sprintf("Уведомление #%d", i),
				SentAt:    now,
			}
			_, err := s.storage.SaveNotification(notification)
			errs <- err
		}(i)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		s.NoError(err) // Все вставки должны пройти
	}

	// Проверим количество
	var count int
	err := s.connection.QueryRowContext(
		context.Background(),
		`SELECT COUNT(*) FROM notifications WHERE group_id = $1`,
		groupID,
	).Scan(&count)
	s.NoError(err)
	s.Equal(10, count)
}

func (s *NotificationsStorageTestSuite) TestSaveNotification_AfterGroupDeleted() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, now, 1)

	// Сначала удалим группу
	_, err := s.connection.ExecContext(
		context.Background(),
		`DELETE FROM groups WHERE id = $1`,
		groupID,
	)
	s.NoError(err)

	// Теперь попробуем вставить уведомление
	notification := entities.Notification{
		GroupID:   groupID,
		MessageID: 1001,
		Text:      "После удаления группы",
		SentAt:    now,
	}

	_, err = s.storage.SaveNotification(notification)
	s.Error(err)
	s.Contains(err.Error(), "violates foreign key constraint")
}
