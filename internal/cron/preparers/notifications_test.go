package preparers

import (
	"fmt"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	mockbot "github.com/DKhorkov/plantsCareTelegramBot/mocks/bot"
	mockusecases "github.com/DKhorkov/plantsCareTelegramBot/mocks/usecases"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gopkg.in/telebot.v4"
	"sync"
	"testing"
	"time"
)

func TestNewNotificationsPreparer(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockUsecases := mockusecases.NewMockUseCases(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	preparer := NewNotificationsPreparer(mockBot, mockUsecases, mockLogger, 10, 5)

	assert.NotNil(t, preparer)
	assert.Equal(t, mockBot, preparer.bot)
	assert.Equal(t, mockUsecases, preparer.useCases)
	assert.Equal(t, mockLogger, preparer.logger)
	assert.Equal(t, 10, preparer.limit)
	assert.Equal(t, 5, preparer.offset)
	assert.NotNil(t, preparer.notifiedGroups)
}

func TestNotificationsPreparer_alreadyNotified(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBot := mockbot.NewMockBot(ctrl)
	mockUsecases := mockusecases.NewMockUseCases(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)

	tests := []struct {
		name         string
		groupID      int
		storeValue   interface{}
		storeExists  bool
		setupMocks   func()
		want         bool
		expectLogErr bool
	}{
		{
			name:         "not_notified",
			groupID:      1,
			storeExists:  false,
			want:         false,
			expectLogErr: false,
		},
		{
			name:         "notified_today",
			groupID:      1,
			storeValue:   now,
			storeExists:  true,
			want:         true,
			expectLogErr: false,
		},
		{
			name:         "notified_yesterday",
			groupID:      1,
			storeValue:   yesterday,
			storeExists:  true,
			want:         false,
			expectLogErr: false,
		},
		{
			name:         "invalid_type_in_map",
			groupID:      1,
			storeValue:   "not-a-time",
			storeExists:  true,
			want:         false,
			expectLogErr: true,
			setupMocks: func() {
				mockLogger.EXPECT().Error(
					"Failed to parse date",
					"Tracing", gomock.Any(),
				).Times(1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparer := NewNotificationsPreparer(mockBot, mockUsecases, mockLogger, 10, 0)
			preparer.notifiedGroups = new(sync.Map)

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			if tt.storeExists {
				preparer.notifiedGroups.Store(tt.groupID, tt.storeValue)
			}

			group := entities.Group{ID: tt.groupID}
			result := preparer.alreadyNotified(group)

			assert.Equal(t, tt.want, result)
		})
	}
}

func TestNotificationsPreparer_notify(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockBot := mockbot.NewMockBot(ctrl)
	mockUsecases := mockusecases.NewMockUseCases(ctrl)
	mockLogger := mocklogging.NewMockLogger(ctrl)

	group := entities.Group{
		ID:               1,
		UserID:           100,
		Title:            "Группа 1",
		Description:      "Описание",
		LastWateringDate: time.Date(2025, 9, 2, 0, 0, 0, 0, time.Local),
		WateringInterval: 7,
	}
	user := entities.User{TelegramID: 12345}
	plants := []entities.Plant{{ID: 1, Title: "Фикус"}}
	msg := &telebot.Message{ID: 987, Text: "Напоминание: пора поливать!"}

	tests := []struct {
		name         string
		setupMocks   func()
		expectError  bool
		expectStored bool
	}{
		{
			name: "success_flow",
			setupMocks: func() {
				mockUsecases.EXPECT().GetUserByID(group.UserID).Return(&user, nil).Times(1)
				mockUsecases.EXPECT().GetGroupPlants(group.ID).Return(plants, nil).Times(1)
				mockBot.EXPECT().Send(
					&telebot.Chat{ID: int64(user.TelegramID)},
					gomock.Any(),
					gomock.Any(),
				).Return(msg, nil).Times(1)
				mockUsecases.EXPECT().SaveNotification(gomock.Any()).Return(&entities.Notification{}, nil).Times(1)
			},
			expectError:  false,
			expectStored: true,
		},
		{
			name: "error_get_user",
			setupMocks: func() {
				mockUsecases.EXPECT().GetUserByID(group.UserID).Return(&entities.User{}, fmt.Errorf("user not found")).Times(1)
			},
			expectError:  true,
			expectStored: false,
		},
		{
			name: "error_get_plants",
			setupMocks: func() {
				mockUsecases.EXPECT().GetUserByID(group.UserID).Return(&user, nil).Times(1)
				mockUsecases.EXPECT().GetGroupPlants(group.ID).Return(nil, fmt.Errorf("load error")).Times(1)
			},
			expectError:  true,
			expectStored: false,
		},
		{
			name: "error_send_message",
			setupMocks: func() {
				mockUsecases.EXPECT().GetUserByID(group.UserID).Return(&user, nil).Times(1)
				mockUsecases.EXPECT().GetGroupPlants(group.ID).Return(plants, nil).Times(1)
				mockBot.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("send failed")).Times(1)
				mockLogger.EXPECT().Error("Failed to send message", "Error", gomock.Any()).Times(1)
			},
			expectError:  true,
			expectStored: false,
		},
		{
			name: "error_save_notification",
			setupMocks: func() {
				mockUsecases.EXPECT().GetUserByID(group.UserID).Return(&user, nil).Times(1)
				mockUsecases.EXPECT().GetGroupPlants(group.ID).Return(plants, nil).Times(1)
				mockBot.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(msg, nil).Times(1)
				mockUsecases.EXPECT().SaveNotification(gomock.Any()).Return(&entities.Notification{}, fmt.Errorf("save failed")).Times(1)
			},
			expectError:  true,
			expectStored: true, // Store уже произошёл
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparer := NewNotificationsPreparer(mockBot, mockUsecases, mockLogger, 10, 0)
			preparer.notifiedGroups = new(sync.Map)

			if tt.setupMocks != nil {
				tt.setupMocks()
			}

			err := preparer.notify(group)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			_, stored := preparer.notifiedGroups.Load(group.ID)
			assert.Equal(t, tt.expectStored, stored)
		})
	}
}

func TestNotificationsPreparer_preparePlantsText(t *testing.T) {
	preparer := &NotificationsPreparer{} // не требует зависимостей

	tests := []struct {
		name    string
		plants  []entities.Plant
		want    string
		wantErr bool
	}{
		{
			name:    "empty_plants",
			plants:  nil,
			want:    "В данный сценарий полива пока что не было добавлено ни одно растение!\n",
			wantErr: false,
		},
		{
			name: "one_plant",
			plants: []entities.Plant{
				{Title: "Фикус"},
			},
			want:    "1) Фикус\n",
			wantErr: false,
		},
		{
			name: "two_plants",
			plants: []entities.Plant{
				{Title: "Фикус"},
				{Title: "Папоротник"},
			},
			want:    "1) Фикус\n2) Папоротник\n",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := preparer.preparePlantsText(tt.plants)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
