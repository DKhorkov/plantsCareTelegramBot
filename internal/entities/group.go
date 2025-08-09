package entities

import "time"

type Group struct {
	ID                    int       `json:"id"`
	UserID                int       `json:"user_id"`
	Title                 string    `json:"title"`
	Description           *string   `json:"description"`
	LastWateringDate      time.Time `json:"last_watering_date"`
	NextWateringDate      time.Time `json:"next_watering_date"`
	WateringIntervalHours int       `json:"watering_interval_hours"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
