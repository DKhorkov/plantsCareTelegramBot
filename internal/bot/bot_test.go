//go:build integration

package bot_test

import (
	"fmt"
	"github.com/DKhorkov/libs/loadenv"
	"github.com/DKhorkov/plantsCareTelegramBot/internal/bot"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNew_TableDriven_Simple(t *testing.T) {
	// Запускать в корневой дирекетории проекта:
	token := loadenv.GetEnv("BOT_TOKEN", "")

	tests := []struct {
		name        string
		token       string
		timeout     time.Duration
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       token,
			timeout:     10 * time.Second,
			expectError: false,
		},
		{
			name:        "Empty token",
			token:       "",
			timeout:     10 * time.Second,
			expectError: true,
		},
		{
			name:        "Token with spaces",
			token:       fmt.Sprintf("  %s ", token),
			timeout:     5 * time.Second,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := bot.New(tt.token, tt.timeout)

			if tt.expectError {
				assert.Nil(t, b)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, b)
				assert.NoError(t, err)
			}
		})
	}
}
