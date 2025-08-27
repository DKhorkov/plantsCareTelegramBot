package utils

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"unicode"
)

func TestGenUniqueParam(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		expected bool // ожидаем корректный результат
	}{
		{name: "Simple base", base: "user", expected: true},
		{name: "Empty base", base: "", expected: true},
		{name: "With underscore", base: "my_param", expected: true},
		{name: "Numeric base", base: "123", expected: true},
		{name: "Special chars base", base: "test@#", expected: true}, // не фильтруем, но можно добавить
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenUniqueParam(tt.base)

			// Проверяем длину: len(base) + 1 (подчёркивание) + hashLength
			expectedMinLength := len(tt.base) + 1 + hashLength
			assert.Equal(t, expectedMinLength, len(result), "Длина должна быть base + 1 + %d", hashLength)

			// Проверяем формат: base + "_" + random
			parts := strings.Split(result, "_")
			assert.GreaterOrEqual(t, len(parts), 2, "Должно быть хотя бы две части, разделённые _")

			prefix := strings.Join(parts[:len(parts)-1], "_") // на случай, если base содержит _
			suffix := parts[len(parts)-1]

			assert.Equal(t, tt.base, prefix, "Префикс должен совпадать с base")
			assert.Equal(t, hashLength, len(suffix), "Суффикс должен быть длины %d", hashLength)

			// Проверяем, что суффикс — только латинские буквы
			for _, r := range suffix {
				assert.True(t, unicode.IsLetter(r), "Суффикс должен содержать только буквы")
				assert.True(t, (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z'), "Только латинские символы")
			}
		})
	}
}

// Проверка уникальности при многократных вызовах
func TestGenUniqueParam_Uniqueness(t *testing.T) {
	trials := 1000
	seen := make(map[string]bool)
	base := "item"

	for i := 0; i < trials; i++ {
		param := GenUniqueParam(base)
		if seen[param] {
			t.Fatalf("Обнаружен дубликат: %s (на итерации %d)", param, i)
		}
		seen[param] = true
	}

	assert.Equal(t, trials, len(seen), "Все сгенерированные параметры должны быть уникальными")
}
