package calendar

import (
	"gopkg.in/telebot.v4"
	"testing"
	"time"
)

// Тесты для метода validate
func TestOptions_Validate(t *testing.T) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := now.Month()

	tests := []struct {
		name    string
		opts    options
		wantErr bool
	}{
		{
			name: "Валидные значения по умолчанию",
			opts: options{
				initialYear:  currentYear,
				initialMonth: currentMonth,
				yearsRange:   [2]int{1970, 2050},
				language:     "en",
			},
			wantErr: false,
		},
		{
			name: "Год вне диапазона yearsRange",
			opts: options{
				initialYear:  1900,
				initialMonth: time.January,
				yearsRange:   [2]int{1970, 2050},
			},
			wantErr: true,
		},
		{
			name: "Год в пределах диапазона",
			opts: options{
				initialYear:  2000,
				initialMonth: time.March,
				yearsRange:   [2]int{1970, 2050},
			},
			wantErr: false,
		},
		{
			name: "Нижняя граница yearsRange — валидна",
			opts: options{
				initialYear:  1970,
				initialMonth: time.December,
				yearsRange:   [2]int{1970, 2050},
			},
			wantErr: false,
		},
		{
			name: "Верхняя граница yearsRange — валидна",
			opts: options{
				initialYear:  2050,
				initialMonth: time.January,
				yearsRange:   [2]int{1970, 2050},
			},
			wantErr: false,
		},
		{
			name: "yearsRange выходит за лимит Unix (слишком большой верх)",
			opts: options{
				initialYear:  2000,
				initialMonth: time.June,
				yearsRange:   [2]int{1970, MaxYearLimit + 1},
			},
			wantErr: true,
		},
		{
			name: "yearsRange выходит за лимит Unix (слишком маленький низ)",
			opts: options{
				initialYear:  2000,
				initialMonth: time.June,
				yearsRange:   [2]int{MinYearLimit - 1, 2050},
			},
			wantErr: true,
		},
		{
			name: "Месяц 0 — невалидный",
			opts: options{
				initialYear:  2000,
				initialMonth: time.Month(0),
				yearsRange:   [2]int{1970, 2050},
			},
			wantErr: true,
		},
		{
			name: "Месяц 13 — невалидный",
			opts: options{
				initialYear:  2000,
				initialMonth: time.Month(13),
				yearsRange:   [2]int{1970, 2050},
			},
			wantErr: true,
		},
		{
			name: "Пустой yearsRange (0,0) — но валидация не проходит, если не установлены вручную",
			opts: options{
				initialYear:  2000,
				initialMonth: time.June,
				yearsRange:   [2]int{0, 0},
			},
			wantErr: true, // потому что validate не знает, что это "по умолчанию"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.validate()

			if tt.wantErr && err == nil {
				t.Fatalf("Ожидалась ошибка, но её не было")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Не ожидалась ошибка, но получена: %v", err)
			}
		})
	}
}

// Тесты для функций-опций
func TestOptionFuncs(t *testing.T) {
	now := time.Now()
	currentYear := now.Year()
	currentMonth := now.Month()

	button := telebot.InlineButton{Unique: "back"}

	tests := []struct {
		name    string
		option  Option
		modify  func(*options)
		check   func(*options) bool
		wantErr bool
	}{
		{
			name:   "WithInitialYear: 0 → устанавливается текущий год",
			option: WithInitialYear(0),
			check: func(opts *options) bool {
				return opts.initialYear == currentYear
			},
			wantErr: false,
		},
		{
			name:   "WithInitialYear: 2020 → устанавливается 2020",
			option: WithInitialYear(2020),
			check: func(opts *options) bool {
				return opts.initialYear == 2020
			},
			wantErr: false,
		},
		{
			name:   "WithInitialMonth: 0 → устанавливается текущий месяц",
			option: WithInitialMonth(0),
			check: func(opts *options) bool {
				return opts.initialMonth == currentMonth
			},
			wantErr: false,
		},
		{
			name:   "WithInitialMonth: March → устанавливается March",
			option: WithInitialMonth(time.March),
			check: func(opts *options) bool {
				return opts.initialMonth == time.March
			},
			wantErr: false,
		},
		{
			name:   "WithYearsRange: [2000,2100] → устанавливается",
			option: WithYearsRange([2]int{2000, 2100}),
			check: func(opts *options) bool {
				return opts.yearsRange == [2]int{2000, 2100}
			},
			wantErr: false,
		},
		{
			name:   "WithYearsRange: [0,0] → устанавливается лимит Unix",
			option: WithYearsRange([2]int{0, 0}),
			check: func(opts *options) bool {
				return opts.yearsRange == [2]int{MinYearLimit, MaxYearLimit}
			},
			wantErr: false,
		},
		{
			name:   "WithLanguage: 'ru' → устанавливается",
			option: WithLanguage("ru"),
			check: func(opts *options) bool {
				return opts.language == "ru"
			},
			wantErr: false,
		},
		{
			name:   "WithLanguage: 'en' → устанавливается",
			option: WithLanguage("en"),
			check: func(opts *options) bool {
				return opts.language == "en"
			},
			wantErr: false,
		},
		{
			name:   "WithBackButton: передана кнопка → указатель сохранён",
			option: WithBackButton(button),
			check: func(opts *options) bool {
				return opts.backButton != nil && opts.backButton.Unique == button.Unique
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &options{}

			err := tt.option(opts)

			if tt.wantErr && err == nil {
				t.Fatalf("Ожидалась ошибка, но её не было")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Не ожидалась ошибка, но получена: %v", err)
			}

			if tt.check != nil && !tt.check(opts) {
				t.Errorf("Проверка не прошла: состояние options не соответствует ожидаемому")
			}
		})
	}
}
