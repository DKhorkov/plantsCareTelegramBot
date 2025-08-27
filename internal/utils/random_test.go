package utils_test

import (
	"github.com/DKhorkov/plantsCareTelegramBot/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode"
)

func TestRandSequence(t *testing.T) {
	tests := []struct {
		name   string
		length int
		valid  bool // ожидаем ли корректную строку
	}{
		{name: "Positive length 1", length: 1, valid: true},
		{name: "Positive length 5", length: 5, valid: true},
		{name: "Positive length 10", length: 10, valid: true},
		{name: "Zero length", length: 0, valid: true}, // допустимо: пустая строка
		// Отрицательные значения не обрабатываются — функция просто создаст срез длины 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.RandSequence(tt.length)

			// Проверяем длину
			assert.Equal(t, tt.length, len(result), "Длина строки должна совпадать с запрошенной")

			// Проверяем содержимое, только если строка не пустая
			if tt.valid && len(result) > 0 {
				for _, r := range result {
					assert.True(t, unicode.IsLetter(r), "Каждый символ должен быть буквой")
					assert.True(t, (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'), "Только латинские буквы")
				}
			}
		})
	}
}

func TestRandSequence_Uniqueness(t *testing.T) {
	// Проверим, что вызовы дают разные строки
	n := 8
	seen := make(map[string]bool)
	trials := 100

	for i := 0; i < trials; i++ {
		s := utils.RandSequence(n)
		if seen[s] {
			t.Fatalf("Обнаружена повторяющаяся строка: %s (на итерации %d)", s, i)
		}
		seen[s] = true
	}

	assert.Equal(t, trials, len(seen), "Все строки должны быть уникальными при 100 вызовах")
}
