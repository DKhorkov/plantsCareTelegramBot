package usecases

import (
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	customerrors "github.com/DKhorkov/plantsCareTelegramBot/internal/errors"
	mockstorage "github.com/DKhorkov/plantsCareTelegramBot/mocks/storage"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestGroupsUseCases_UpdateGroupLastWateringDate(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	baseGroup := entities.Group{
		ID:               1,
		UserID:           123,
		Title:            "Цветы",
		WateringInterval: 7,
		LastWateringDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		NextWateringDate: time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	tests := []struct {
		name                 string
		id                   int
		lastWateringDate     time.Time
		setupMocks           func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantNextWateringDate time.Time
		wantErr              bool
	}{
		{
			name:             "Success - next watering date in future",
			id:               1,
			lastWateringDate: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantNextWateringDate: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			wantErr:              false,
		},
		{
			name:             "Success - next watering date in past → set to today",
			id:               1,
			lastWateringDate: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantNextWateringDate: today,
			wantErr:              false,
		},
		{
			name:             "Failure - group not found",
			id:               1,
			lastWateringDate: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantNextWateringDate: time.Time{},
			wantErr:              true,
		},
		{
			name:             "Failure - storage update error",
			id:               1,
			lastWateringDate: time.Date(2023, 10, 15, 0, 0, 0, 0, time.UTC),
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantNextWateringDate: time.Time{},
			wantErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.UpdateGroupLastWateringDate(tt.id, tt.lastWateringDate)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantNextWateringDate, got.NextWateringDate)
			}
		})
	}
}

func TestGroupsUseCases_UpdateGroupWateringInterval(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	baseGroup := entities.Group{
		ID:               1,
		UserID:           123,
		Title:            "Цветы",
		WateringInterval: 7,
		LastWateringDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		NextWateringDate: time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC),
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	tests := []struct {
		name                 string
		id                   int
		wateringInterval     int
		setupMocks           func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantNextWateringDate time.Time
		wantErr              bool
	}{
		{
			name:             "Success - next watering date in future",
			id:               1,
			wateringInterval: 10,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantNextWateringDate: today,
			wantErr:              false,
		},
		{
			name:             "Success - new interval causes past date → set to today",
			id:               1,
			wateringInterval: 1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantNextWateringDate: today,
			wantErr:              false,
		},
		{
			name:             "Failure - group not found",
			id:               1,
			wateringInterval: 5,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantNextWateringDate: time.Time{},
			wantErr:              true,
		},
		{
			name:             "Failure - storage update error",
			id:               1,
			wateringInterval: 5,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantNextWateringDate: time.Time{},
			wantErr:              true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.UpdateGroupWateringInterval(tt.id, tt.wateringInterval)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.wantNextWateringDate, got.NextWateringDate)
			}
		})
	}
}

func TestGroupsUseCases_GetUserGroups(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		userID     int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		want       []entities.Group
		wantErr    bool
	}{
		{
			name:   "Success - returns groups",
			userID: 123,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserGroups(123).
					Return([]entities.Group{
						{ID: 1, UserID: 123, Title: "Цветы", CreatedAt: now},
						{ID: 2, UserID: 123, Title: "Суккуленты", CreatedAt: now},
					}, nil).
					Times(1)
			},
			want: []entities.Group{
				{ID: 1, UserID: 123, Title: "Цветы", CreatedAt: now},
				{ID: 2, UserID: 123, Title: "Суккуленты", CreatedAt: now},
			},
			wantErr: false,
		},
		{
			name:   "Failure - storage error",
			userID: 123,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetUserGroups(123).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Groups for User with ID=123",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.GetUserGroups(tt.userID)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGroupsUseCases_CountUserGroups(t *testing.T) {
	tests := []struct {
		name       string
		userID     int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		want       int
		wantErr    bool
	}{
		{
			name:   "Success - returns count",
			userID: 123,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CountUserGroups(123).
					Return(5, nil).
					Times(1)
			},
			want:    5,
			wantErr: false,
		},
		{
			name:   "Failure - storage error",
			userID: 123,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CountUserGroups(123).
					Return(0, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to count Groups for User with ID=123",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.CountUserGroups(tt.userID)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGroupsUseCases_CreateGroup(t *testing.T) {
	inputGroup := entities.Group{
		UserID:           123,
		Title:            "Цветы",
		Description:      "Комнатные растения",
		WateringInterval: 7,
	}

	tests := []struct {
		name       string
		group      entities.Group
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		want       *entities.Group
		wantErr    bool
	}{
		{
			name:  "Success - group created",
			group: inputGroup,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CreateGroup(inputGroup).
					Return(1, nil).
					Times(1)
			},
			want: &entities.Group{
				ID:               1,
				UserID:           123,
				Title:            "Цветы",
				Description:      "Комнатные растения",
				WateringInterval: 7,
			},
			wantErr: false,
		},
		{
			name:  "Failure - storage error",
			group: inputGroup,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					CreateGroup(inputGroup).
					Return(0, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to create Group for User with ID=123",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.CreateGroup(tt.group)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.UserID, got.UserID)
				assert.Equal(t, tt.want.Title, got.Title)
			}
		})
	}
}

func TestGroupsUseCases_UpdateGroup(t *testing.T) {
	updatedGroup := entities.Group{
		ID:               1,
		UserID:           123,
		Title:            "Обновлённые цветы",
		Description:      "Обновлённое описание",
		WateringInterval: 10,
		LastWateringDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		NextWateringDate: time.Date(2023, 10, 11, 0, 0, 0, 0, time.UTC),
	}

	tests := []struct {
		name       string
		group      entities.Group
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantErr    bool
	}{
		{
			name:  "Success - group updated",
			group: updatedGroup,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					UpdateGroup(updatedGroup).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name:  "Failure - storage error",
			group: updatedGroup,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					UpdateGroup(updatedGroup).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Group for with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			err := useCases.UpdateGroup(tt.group)

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGroupsUseCases_GetGroup(t *testing.T) {
	now := time.Now()
	group := entities.Group{
		ID:               1,
		UserID:           123,
		Title:            "Цветы",
		Description:      "Комнатные растения",
		WateringInterval: 7,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	tests := []struct {
		name       string
		id         int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		want       *entities.Group
		wantErr    bool
	}{
		{
			name: "Success - group found",
			id:   1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&group, nil).
					Times(1)
			},
			want:    &group,
			wantErr: false,
		},
		{
			name: "Failure - storage error",
			id:   1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.GetGroup(tt.id)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGroupsUseCases_GetGroupsForNotify(t *testing.T) {
	now := time.Now()
	groups := []entities.Group{
		{ID: 1, UserID: 123, Title: "Цветы", NextWateringDate: now},
	}

	tests := []struct {
		name       string
		limit      int
		offset     int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		want       []entities.Group
		wantErr    bool
	}{
		{
			name:   "Success - returns groups for notify",
			limit:  10,
			offset: 0,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroupsForNotify(10, 0).
					Return(groups, nil).
					Times(1)
			},
			want:    groups,
			wantErr: false,
		},
		{
			name:   "Failure - storage error",
			limit:  10,
			offset: 0,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroupsForNotify(10, 0).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Groups for Notify with limit=10 and offset=0",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.GetGroupsForNotify(tt.limit, tt.offset)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestGroupsUseCases_DeleteGroup(t *testing.T) {
	tests := []struct {
		name       string
		id         int
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		wantErr    bool
	}{
		{
			name: "Success - group deleted",
			id:   1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					DeleteGroup(1).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "Failure - storage error",
			id:   1,
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					DeleteGroup(1).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to delete Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			err := useCases.DeleteGroup(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGroupsUseCases_UpdateGroupTitle(t *testing.T) {
	now := time.Now()
	baseGroup := entities.Group{
		ID:               1,
		UserID:           123,
		Title:            "Цветы",
		Description:      "Комнатные",
		WateringInterval: 7,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	tests := []struct {
		name       string
		id         int
		title      string
		setupMocks func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		want       *entities.Group
		wantErr    bool
		wantErrIs  error
	}{
		{
			name:  "Success - title updated",
			id:    1,
			title: "Новые цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					GroupExists(gomock.Any()).
					DoAndReturn(func(group entities.Group) (bool, error) {
						return false, nil // группа с таким названием не существует
					}).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					DoAndReturn(func(group entities.Group) error {
						assert.Equal(t, "Новые цветы", group.Title)
						return nil
					}).
					Times(1)
			},
			want: &entities.Group{
				ID:          1,
				UserID:      123,
				Title:       "Новые цветы",
				Description: "Комнатные",
			},
			wantErr: false,
		},
		{
			name:  "Failure - group not found",
			id:    1,
			title: "Новые цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "Failure - group with title already exists",
			id:    1,
			title: "Суккуленты",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					GroupExists(gomock.Any()).
					DoAndReturn(func(group entities.Group) (bool, error) {
						return true, nil // такая группа уже есть
					}).
					Times(1)
			},
			want:      nil,
			wantErr:   true,
			wantErrIs: customerrors.ErrGroupAlreadyExists,
		},
		{
			name:  "Failure - storage error on GroupExists",
			id:    1,
			title: "Новые цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					GroupExists(gomock.Any()).
					Return(false, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to check existence for Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "Failure - storage error on UpdateGroup",
			id:    1,
			title: "Новые цветы",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					GroupExists(gomock.Any()).
					Return(false, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.UpdateGroupTitle(tt.id, tt.title)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Title, got.Title)
				assert.Equal(t, tt.want.Description, got.Description)
			}
		})
	}
}

func TestGroupsUseCases_UpdateGroupDescription(t *testing.T) {
	now := time.Now()
	baseGroup := entities.Group{
		ID:               1,
		UserID:           123,
		Title:            "Цветы",
		Description:      "Комнатные",
		WateringInterval: 7,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	tests := []struct {
		name        string
		id          int
		description string
		setupMocks  func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger)
		want        *entities.Group
		wantErr     bool
	}{
		{
			name:        "Success - description updated",
			id:          1,
			description: "Обновлённое описание",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					DoAndReturn(func(group entities.Group) error {
						assert.Equal(t, "Обновлённое описание", group.Description)
						return nil
					}).
					Times(1)
			},
			want: &entities.Group{
				ID:          1,
				UserID:      123,
				Title:       "Цветы",
				Description: "Обновлённое описание",
			},
			wantErr: false,
		},
		{
			name:        "Failure - group not found",
			id:          1,
			description: "Описание",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(nil, assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to get Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:        "Failure - storage error on UpdateGroup",
			id:          1,
			description: "Описание",
			setupMocks: func(storage *mockstorage.MockStorage, logger *mocklogging.MockLogger) {
				storage.
					EXPECT().
					GetGroup(1).
					Return(&baseGroup, nil).
					Times(1)

				storage.
					EXPECT().
					UpdateGroup(gomock.Any()).
					Return(assert.AnError).
					Times(1)

				logger.
					EXPECT().
					Error(
						"Failed to update Group with ID=1",
						"Error", assert.AnError,
						"Tracing", gomock.Any(),
					).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := mockstorage.NewMockStorage(ctrl)
			mockLogger := mocklogging.NewMockLogger(ctrl)

			if tt.setupMocks != nil {
				tt.setupMocks(mockStorage, mockLogger)
			}

			useCases := &groupsUseCases{
				storage: mockStorage,
				logger:  mockLogger,
			}

			got, err := useCases.UpdateGroupDescription(tt.id, tt.description)

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.want.ID, got.ID)
				assert.Equal(t, tt.want.Description, got.Description)
			}
		})
	}
}
