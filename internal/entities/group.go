package entities

import "time"

type Group struct {
	ID                    int
	UserID                int
	Title                 string
	Description           *string
	LastWateringDate      time.Time
	NextWateringDate      time.Time
	WateringIntervalHours int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
