package entities

type Temporary struct {
	ID        int
	UserID    int
	Step      int
	MessageID int
	Data      []byte
}
