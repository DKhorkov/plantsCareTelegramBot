package entities

import "time"

type User struct {
	ID         int       `json:"id"`
	TelegramID int       `json:"telegramId"`
	Username   string    `json:"username"`
	Firstname  string    `json:"firstname"`
	Lastname   string    `json:"lastname"`
	IsBot      bool      `json:"isBot"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
