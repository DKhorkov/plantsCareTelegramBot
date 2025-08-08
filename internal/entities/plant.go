package entities

import "time"

type Plant struct {
	ID          int
	GroupID     int
	Title       string
	Description *string
	Photo       []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
