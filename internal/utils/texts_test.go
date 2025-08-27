package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWateringInterval(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{name: "Один день", input: 1, expected: "1 день"},
		{name: "Два дня", input: 2, expected: "2 дня"},
		{name: "Три дня", input: 3, expected: "3 дня"},
		{name: "Четыре дня", input: 4, expected: "4 дня"},
		{name: "Пять дней", input: 5, expected: "5 дней"},
		{name: "Двадцать дней", input: 20, expected: "20 дней"},
		{name: "Двадцать один день", input: 21, expected: "21 день"},
		{name: "Двадцать два дня", input: 22, expected: "22 дня"},
		{name: "Двадцать три дня", input: 23, expected: "23 дня"},
		{name: "Двадцать четыре дня", input: 24, expected: "24 дня"},
		{name: "Двадцать пять дней", input: 25, expected: "25 дней"},
		{name: "Одиннадцать дней", input: 11, expected: "11 дней"},
		{name: "Двенадцать дней", input: 12, expected: "12 дней"},
		{name: "Тринадцать дней", input: 13, expected: "13 дней"},
		{name: "Четырнадцать дней", input: 14, expected: "14 дней"},
		{name: "Пятьдесят один день", input: 51, expected: "51 день"},
		{name: "Сто один день", input: 101, expected: "101 день"},
		{name: "Сто два дня", input: 102, expected: "102 дня"},
		{name: "Сто три дня", input: 103, expected: "103 дня"},
		{name: "Сто четыре дня", input: 104, expected: "104 дня"},
		{name: "Сто пять дней", input: 105, expected: "105 дней"},
		{name: "Сто одиннадцать дней", input: 111, expected: "111 дней"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetWateringInterval(tt.input)
			assert.Equal(t, tt.expected, result, "Для %d должно быть %q", tt.input, tt.expected)
		})
	}
}
