package entities

import "time"

type Group struct {
	ID               int       `json:"id"`
	UserID           int       `json:"userId"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	LastWateringDate time.Time `json:"lastWateringDate"`
	NextWateringDate time.Time `json:"nextWateringDate"`
	WateringInterval int       `json:"wateringInterval"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}
