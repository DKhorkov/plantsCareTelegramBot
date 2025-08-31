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

func TestGroupsStorageTestSuite(t *testing.T) {
	suite.Run(t, new(GroupsStorageTestSuite))
}

type GroupsStorageTestSuite struct {
	suite.Suite

	cwd         string
	ctx         context.Context
	dbConnector db.Connector
	connection  *sql.Conn
	storage     *groupsStorage
	logger      *mocklogging.MockLogger
}

func (s *GroupsStorageTestSuite) SetupSuite() {
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

	s.storage = &groupsStorage{
		dbConnector: s.dbConnector,
		logger:      s.logger,
	}
}

func (s *GroupsStorageTestSuite) SetupTest() {
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

func (s *GroupsStorageTestSuite) TearDownTest() {
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

func (s *GroupsStorageTestSuite) TearDownSuite() {
	s.NoError(s.dbConnector.Close())
}

func (s *GroupsStorageTestSuite) createUser(now time.Time, offset int) int {
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

func (s *GroupsStorageTestSuite) createGroupForUser(userID int, title string, now time.Time, offset int) int {
	var groupID int
	err := s.connection.QueryRowContext(
		context.Background(),
		`
			INSERT INTO groups (
				user_id, title, description,
				last_watering_date, next_watering_date, watering_interval,
				created_at, updated_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $7)
			RETURNING id
		`,
		userID,
		title,
		fmt.Sprintf("Описание %d", offset),
		now.AddDate(0, 0, -offset),
		now.AddDate(0, 0, offset),
		7+offset,
		now,
	).Scan(&groupID)
	s.NoError(err)
	return groupID
}

func (s *GroupsStorageTestSuite) TestCreateGroup_Success() {
	now := time.Now().UTC()

	userID := s.createUser(now, 1) // создаём пользователя

	group := entities.Group{
		UserID:           userID,
		Title:            "Офисные растения",
		Description:      "Полив раз в неделю",
		LastWateringDate: now.AddDate(0, 0, -7),
		NextWateringDate: now,
		WateringInterval: 7,
	}

	groupID, err := s.storage.CreateGroup(group)
	s.NoError(err)
	s.Greater(groupID, 0)

	// Проверяем, что группа сохранилась
	var storedGroup entities.Group
	columns := db.GetEntityColumns(&storedGroup)
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT `+selectAllColumns+` FROM groups WHERE id = $1`,
		groupID,
	).Scan(columns...)
	s.NoError(err)

	s.Equal(groupID, storedGroup.ID)
	s.Equal(group.UserID, storedGroup.UserID)
	s.Equal(group.Title, storedGroup.Title)
	s.Equal(group.Description, storedGroup.Description)
	s.WithinDuration(group.LastWateringDate, storedGroup.LastWateringDate, time.Second)
	s.WithinDuration(group.NextWateringDate, storedGroup.NextWateringDate, time.Second)
	s.Equal(group.WateringInterval, storedGroup.WateringInterval)
	s.WithinDuration(now, storedGroup.CreatedAt, 2*time.Second)
	s.WithinDuration(now, storedGroup.UpdatedAt, 2*time.Second)
}

func (s *GroupsStorageTestSuite) TestUpdateGroup_Success() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, "Старое название", now, 1)

	updatedGroup := entities.Group{
		ID:               groupID,
		UserID:           userID,
		Title:            "Новое название",
		Description:      "Новое описание",
		LastWateringDate: now.AddDate(0, 0, -10),
		NextWateringDate: now.AddDate(0, 0, 5),
		WateringInterval: 14,
	}

	err := s.storage.UpdateGroup(updatedGroup)
	s.NoError(err)

	var storedGroup entities.Group
	columns := db.GetEntityColumns(&storedGroup)
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT `+selectAllColumns+` FROM groups WHERE id = $1`,
		groupID,
	).Scan(columns...)
	s.NoError(err)

	s.Equal("Новое название", storedGroup.Title)
	s.Equal("Новое описание", storedGroup.Description)
	s.Equal(14, storedGroup.WateringInterval)
	s.True(storedGroup.UpdatedAt.After(now))
}

func (s *GroupsStorageTestSuite) TestUpdateGroup_NonExistent() {
	nonExistentGroup := entities.Group{
		ID:               999999,
		UserID:           1,
		Title:            "Не существует",
		WateringInterval: 7,
	}

	err := s.storage.UpdateGroup(nonExistentGroup)
	s.NoError(err) // Должен вернуть nil, даже если нет строки (UPDATE 0)
}

func (s *GroupsStorageTestSuite) TestGroupExists_Exists() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	s.createGroupForUser(userID, "Суккуленты", now, 1)

	exists, err := s.storage.GroupExists(entities.Group{
		UserID: userID,
		Title:  "Суккуленты",
	})
	s.NoError(err)
	s.True(exists)
}

func (s *GroupsStorageTestSuite) TestGroupExists_NotExists() {
	exists, err := s.storage.GroupExists(entities.Group{
		UserID: 999,
		Title:  "Не существующая",
	})
	s.NoError(err)
	s.False(exists)
}

func (s *GroupsStorageTestSuite) TestGroupExists_DifferentUserSameTitle() {
	now := time.Now().UTC()
	userID1 := s.createUser(now, 1)
	userID2 := s.createUser(now, 2)
	s.createGroupForUser(userID1, "Одинаковое имя", now, 1)

	exists, err := s.storage.GroupExists(entities.Group{
		UserID: userID2,
		Title:  "Одинаковое имя",
	})
	s.NoError(err)
	s.False(exists)
}

func (s *GroupsStorageTestSuite) TestDeleteGroup_Success() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, "Для удаления", now, 1)

	err := s.storage.DeleteGroup(groupID)
	s.NoError(err)

	var count int
	err = s.connection.QueryRowContext(
		context.Background(),
		`SELECT COUNT(*) FROM groups WHERE id = $1`,
		groupID,
	).Scan(&count)
	s.NoError(err)
	s.Equal(0, count)
}

func (s *GroupsStorageTestSuite) TestDeleteGroup_NonExistent() {
	err := s.storage.DeleteGroup(999999)
	s.NoError(err) // Удаление несуществующей строки — не ошибка
}

func (s *GroupsStorageTestSuite) TestGetUserGroups_HasGroups() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	s.createGroupForUser(userID, "Группа 1", now, 1)
	s.createGroupForUser(userID, "Группа 2", now, 2)

	groups, err := s.storage.GetUserGroups(userID)
	s.NoError(err)
	s.Len(groups, 2)
	s.Equal("Группа 1", groups[0].Title)
	s.Equal("Группа 2", groups[1].Title)
}

func (s *GroupsStorageTestSuite) TestGetUserGroups_NoGroups() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)

	groups, err := s.storage.GetUserGroups(userID)
	s.NoError(err)
	s.Empty(groups)
}

func (s *GroupsStorageTestSuite) TestCountUserGroups_HasGroups() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	s.createGroupForUser(userID, "Группа A", now, 1)
	s.createGroupForUser(userID, "Группа B", now, 2)
	s.createGroupForUser(userID, "Группа C", now, 3)

	count, err := s.storage.CountUserGroups(userID)
	s.NoError(err)
	s.Equal(3, count)
}

func (s *GroupsStorageTestSuite) TestCountUserGroups_NoGroups() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)

	count, err := s.storage.CountUserGroups(userID)
	s.NoError(err)
	s.Equal(0, count)
}

func (s *GroupsStorageTestSuite) TestGetGroup_Exists() {
	now := time.Now().UTC()
	userID := s.createUser(now, 1)
	groupID := s.createGroupForUser(userID, "Тестовая", now, 1)

	group, err := s.storage.GetGroup(groupID)
	s.NoError(err)
	s.NotNil(group)
	s.Equal(groupID, group.ID)
	s.Equal("Тестовая", group.Title)
}

func (s *GroupsStorageTestSuite) TestGetGroup_NotFound() {
	group, err := s.storage.GetGroup(999999)
	s.Error(err)
	s.Nil(group)
}

func (s *GroupsStorageTestSuite) TestGetGroupsForNotify_HasDueGroups() {
	now := time.Now().UTC()
	userID1 := s.createUser(now, 1)
	userID2 := s.createUser(now, 2)

	// Просроченные группы
	s.createGroupForUser(userID1, "Просрочено вчера", now.AddDate(0, 0, -1), 1)
	s.createGroupForUser(userID2, "Просрочено 2 дня назад", now.AddDate(0, 0, -2), 2)

	// Группа с будущим next_watering_date — не должна попасть
	s.createGroupForUser(s.createUser(now, 3), "Ещё не время", now.AddDate(0, 0, 1), 3)

	groups, err := s.storage.GetGroupsForNotify(10, 0)
	s.NoError(err)
	s.Len(groups, 2)

	titles := []string{groups[0].Title, groups[1].Title}
	s.Contains(titles, "Просрочено вчера")
	s.Contains(titles, "Просрочено 2 дня назад")
}

func (s *GroupsStorageTestSuite) TestGetGroupsForNotify_NoDueGroups() {
	now := time.Now().UTC()
	s.createGroupForUser(s.createUser(now, 1), "Будущее", now.AddDate(0, 0, 1), 1)

	groups, err := s.storage.GetGroupsForNotify(10, 0)
	s.NoError(err)
	s.Empty(groups)
}

func (s *GroupsStorageTestSuite) TestGetGroupsForNotify_WithLimitAndOffset() {
	now := time.Now().UTC()

	// Создаём 5 просроченных групп
	for i := 1; i <= 5; i++ {
		s.createGroupForUser(s.createUser(now, i), fmt.Sprintf("Группа %d", i), now.AddDate(0, 0, -i), i)
	}

	// Лимит 2, оффсет 1 → ожидаем группы 2 и 3 (по ID)
	groups, err := s.storage.GetGroupsForNotify(2, 1)
	s.NoError(err)
	s.Len(groups, 2)

	// Порядок по ID — значит, это 2 и 3
	s.Equal(2, groups[0].ID)
	s.Equal(3, groups[1].ID)
}
