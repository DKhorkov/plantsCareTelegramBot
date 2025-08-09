package entities

type Temporary struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Step   int    `json:"step"`
	Data   []byte `json:"data"`
}
