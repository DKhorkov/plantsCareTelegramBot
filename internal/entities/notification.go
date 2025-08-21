package entities

import "time"

type Notification struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"groupId"`
	MessageID int       `json:"messageId"`
	Text      string    `json:"text"`
	SentAt    time.Time `json:"sentAt"`
}
