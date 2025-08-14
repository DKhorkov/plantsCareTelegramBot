package entities

type Notification struct {
	ID        int `json:"id"`
	GroupID   int `json:"groupId"`
	MessageID int `json:"messageId"`
}
