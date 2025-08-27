package entities_test

import (
	"encoding/json"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/entities"
	"testing"
	"time"
)

// TestTemporary_GetGroup тестирует метод GetGroup
func TestTemporary_GetGroup(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond) // убираем наносекунды

	group := &entities.Group{
		ID:               1,
		UserID:           100,
		Title:            "Цветы на подоконнике",
		Description:      "Герань, фиалки, кактусы",
		LastWateringDate: now,
		NextWateringDate: now.AddDate(0, 0, 7),
		WateringInterval: 7,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	validData, _ := json.Marshal(group)

	tests := []struct {
		name        string
		temporary   *entities.Temporary
		expectError bool
		expected    *entities.Group
	}{
		{
			name: "Валидные данные — должен успешно распаковать Group",
			temporary: &entities.Temporary{
				ID:     1,
				UserID: 100,
				Step:   2,
				Data:   validData,
			},
			expectError: false,
			expected:    group,
		},
		{
			name: "Пустые данные — ошибка unmarshal",
			temporary: &entities.Temporary{
				ID:     2,
				UserID: 101,
				Step:   1,
				Data:   []byte{},
			},
			expectError: true,
			expected:    nil,
		},
		{
			name: "Некорректный JSON — ошибка парсинга полей",
			temporary: &entities.Temporary{
				ID:     3,
				UserID: 102,
				Step:   3,
				Data:   []byte(`{"id": "не число", "userId": 1}`),
			},
			expectError: true,
			expected:    nil,
		},
		{
			name: "nil Data — ошибка unmarshal",
			temporary: &entities.Temporary{
				ID:     4,
				UserID: 103,
				Step:   4,
				Data:   nil,
			},
			expectError: true,
			expected:    nil,
		},
		{
			name: "Частичные данные — не хватает обязательных полей",
			temporary: &entities.Temporary{
				ID:     5,
				UserID: 104,
				Step:   5,
				Data:   []byte(`{"id": 1}`), // нет userID, createdAt и т.д.
			},
			expectError: false, // unmarshal пройдёт, но данные будут неполные
			expected: &entities.Group{
				ID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.temporary.GetGroup()

			if tt.expectError {
				if err == nil {
					t.Fatalf("Ожидалась ошибка, но её не было")
				}
				return
			}

			if err != nil {
				t.Fatalf("Не ожидалась ошибка, но получена: %v", err)
			}

			if result.ID != tt.expected.ID {
				t.Errorf("Ожидался ID=%d, получено=%d", tt.expected.ID, result.ID)
			}

			if result.UserID != tt.expected.UserID {
				t.Errorf("Ожидался UserID=%d, получено=%d", tt.expected.UserID, result.UserID)
			}

			if result.Title != tt.expected.Title {
				t.Errorf("Ожидался Title=%s, получено=%s", tt.expected.Title, result.Title)
			}

			if result.Description != tt.expected.Description {
				t.Errorf("Ожидалась Description=%s, получено=%s", tt.expected.Description, result.Description)
			}

			if !result.LastWateringDate.Equal(tt.expected.LastWateringDate) {
				t.Errorf("LastWateringDate не совпадает: ожидалось %v, получено %v", tt.expected.LastWateringDate, result.LastWateringDate)
			}

			if !result.NextWateringDate.Equal(tt.expected.NextWateringDate) {
				t.Errorf("NextWateringDate не совпадает: ожидалось %v, получено %v", tt.expected.NextWateringDate, result.NextWateringDate)
			}

			if result.WateringInterval != tt.expected.WateringInterval {
				t.Errorf("Ожидался интервал=%d, получено=%d", tt.expected.WateringInterval, result.WateringInterval)
			}
		})
	}
}

// TestTemporary_GetPlant тестирует метод GetPlant
func TestTemporary_GetPlant(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)

	plant := &entities.Plant{
		ID:          1,
		GroupID:     1,
		UserID:      100,
		Title:       "Фиалка",
		Description: "Фиолетовая, цветёт весной",
		Photo:       []byte{0xFF, 0xD8, 0xFF, 0xE0}, // JPEG сигнатура
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	validData, _ := json.Marshal(plant)

	tests := []struct {
		name        string
		temporary   *entities.Temporary
		expectError bool
		expected    *entities.Plant
	}{
		{
			name: "Валидные данные — должен успешно распаковать Plant",
			temporary: &entities.Temporary{
				ID:     1,
				UserID: 100,
				Step:   5,
				Data:   validData,
			},
			expectError: false,
			expected:    plant,
		},
		{
			name: "Пустые данные — ошибка unmarshal",
			temporary: &entities.Temporary{
				ID:     2,
				UserID: 101,
				Step:   1,
				Data:   []byte{},
			},
			expectError: true,
			expected:    nil,
		},
		{
			name: "Некорректный JSON — ошибка парсинга",
			temporary: &entities.Temporary{
				ID:     3,
				UserID: 102,
				Step:   3,
				Data:   []byte(`{"id": "abc", "groupId": 1}`),
			},
			expectError: true,
			expected:    nil,
		},
		{
			name: "nil Data — ошибка unmarshal",
			temporary: &entities.Temporary{
				ID:     4,
				UserID: 103,
				Step:   4,
				Data:   nil,
			},
			expectError: true,
			expected:    nil,
		},
		{
			name: "Частичные данные — только ID",
			temporary: &entities.Temporary{
				ID:     5,
				UserID: 104,
				Step:   5,
				Data:   []byte(`{"id": 1}`),
			},
			expectError: false,
			expected: &entities.Plant{
				ID: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.temporary.GetPlant()

			if tt.expectError {
				if err == nil {
					t.Fatalf("Ожидалась ошибка, но её не было")
				}
				return
			}

			if err != nil {
				t.Fatalf("Не ожидалась ошибка, но получена: %v", err)
			}

			if result.ID != tt.expected.ID {
				t.Errorf("Ожидался ID=%d, получено=%d", tt.expected.ID, result.ID)
			}

			if result.GroupID != tt.expected.GroupID {
				t.Errorf("Ожидался GroupID=%d, получено=%d", tt.expected.GroupID, result.GroupID)
			}

			if result.UserID != tt.expected.UserID {
				t.Errorf("Ожидался UserID=%d, получено=%d", tt.expected.UserID, result.UserID)
			}

			if result.Title != tt.expected.Title {
				t.Errorf("Ожидался Title=%s, получено=%s", tt.expected.Title, result.Title)
			}

			if result.Description != tt.expected.Description {
				t.Errorf("Ожидалась Description=%s, получено=%s", tt.expected.Description, result.Description)
			}

			if len(result.Photo) != len(tt.expected.Photo) {
				t.Errorf("Размер фото не совпадает")
			} else if len(result.Photo) > 0 {
				for i, b := range tt.expected.Photo {
					if result.Photo[i] != b {
						t.Errorf("Фото отличается на байте %d", i)
						break
					}
				}
			}
		})
	}
}
