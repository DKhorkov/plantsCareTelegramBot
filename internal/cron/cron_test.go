package cron

import (
	"errors"
	mocklogging "github.com/DKhorkov/libs/logging/mocks"
	customerrors "github.com/DKhorkov/plantsCareTelegramBot/internal/errors"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/interfaces"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"sync"
	"testing"
	"time"
)

func TestNew_CronInitialization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		logger   *mocklogging.MockLogger
		callback interfaces.Callback
		interval time.Duration
	}{
		{
			name:     "Valid dependencies",
			logger:   mocklogging.NewMockLogger(ctrl),
			callback: func() error { return nil },
			interval: 100 * time.Millisecond,
		},
		{
			name:     "Nil logger",
			logger:   nil,
			callback: func() error { return nil },
			interval: 100 * time.Millisecond,
		},
		{
			name:     "Nil callback",
			logger:   mocklogging.NewMockLogger(ctrl),
			callback: nil,
			interval: 100 * time.Millisecond,
		},
		{
			name:     "Zero interval",
			logger:   mocklogging.NewMockLogger(ctrl),
			callback: func() error { return nil },
			interval: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cron := New(tt.logger, tt.callback, tt.interval)

			assert.NotNil(t, cron, "Cron должен быть создан")
			assert.Equal(t, tt.logger, cron.logger, "Logger должен совпадать")
			assert.Equal(t, tt.interval, cron.interval, "Interval должен совпадать")
			assert.NotNil(t, cron.stopChan, "stopChan должен быть инициализирован")

			if tt.callback == nil {
				assert.Nil(t, cron.callback, "Callback должен быть nil")
			} else {
				assert.NotNil(t, cron.callback, "Callback не должен быть nil")
				// Дополнительно: можно вызвать, если это безопасно
				// (но в этом тесте — не нужно)
			}
		})
	}
}

func TestCron_Run_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocklogging.NewMockLogger(ctrl)
	called := false
	var mu sync.Mutex

	callback := func() error {
		mu.Lock()
		called = true
		mu.Unlock()
		return nil
	}

	cron := New(logger, callback, 10*time.Millisecond)

	// Запускаем Run в отдельной горутине
	go func() {
		_ = cron.Run()
	}()

	// Даём время на запуск
	time.Sleep(25 * time.Millisecond)

	// Проверяем, что callback вызвался
	mu.Lock()
	assert.True(t, called, "Callback должен быть вызван хотя бы раз")
	mu.Unlock()

	// Останавливаем
	err := cron.Stop()
	assert.NoError(t, err, "Stop не должен возвращать ошибку")

	// Ждём, чтобы убедиться, что после остановки ничего не происходит
	time.Sleep(10 * time.Millisecond)
}

func TestCron_Run_CallbackError(t *testing.T) {
	ctrl := gomock.NewController(t)
	logger := mocklogging.NewMockLogger(ctrl)
	expectedErr := errors.New("callback failed")

	cron := New(logger, func() error {
		return expectedErr
	}, 10*time.Millisecond)

	done := make(chan error, 1)
	go func() {
		done <- cron.Run()
	}()

	select {
	case err := <-done:
		assert.Equal(t, expectedErr, err, "Run должен вернуть ошибку от callback")
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Run не вернул ошибку вовремя")
	}
}

func TestCron_Run_PanicRecovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogging.NewMockLogger(ctrl)
	logger.EXPECT().
		Error("Recovered from panic", "Recovered", "test panic").
		Times(1)

	cron := New(logger, func() error {
		panic("test panic")
	}, 10*time.Millisecond)

	err := cron.Run()

	assert.Error(t, err)
	assert.Equal(t, customerrors.ErrPanic, err, "Должна быть возвращена ошибка паники")
}
