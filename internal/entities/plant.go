package entities

import "time"

type Plant struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"groupId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Photo       []byte    `json:"photo"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
