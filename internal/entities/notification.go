package entities

type Notification struct {
	ID        int `json:"id"`
	GroupID   int `json:"group_id"`
	MessageID int `json:"message_id"`
}
