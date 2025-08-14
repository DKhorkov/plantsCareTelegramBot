package entities

import "time"

type Group struct {
	ID                    int       `json:"id"`
	UserID                int       `json:"userId"`
	Title                 string    `json:"title"`
	Description           string    `json:"description"`
	LastWateringDate      time.Time `json:"lastWateringDate"`
	NextWateringDate      time.Time `json:"nextWateringDate"`
	WateringIntervalHours int       `json:"wateringIntervalHours"`
	CreatedAt             time.Time `json:"createdAt"`
	UpdatedAt             time.Time `json:"updatedAt"`
}
