package entities

import "time"

type User struct {
	ID         int       `json:"id"`
	TelegramID int       `json:"telegram_id"`
	Username   string    `json:"username"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	IsBot      bool      `json:"is_bot"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
