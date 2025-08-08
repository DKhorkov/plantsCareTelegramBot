package entities

import "time"

type User struct {
	ID         int
	TelegramID int
	Username   string
	Firstname  string
	Lastname   string
	IsBot      bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
