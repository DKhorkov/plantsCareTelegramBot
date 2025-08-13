package entities

import "encoding/json"

type Temporary struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Step      int    `json:"step"`
	MessageID *int   `json:"message_id,omitempty"`
	Data      []byte `json:"data"`
}

func (t *Temporary) GetGroup() (*Group, error) {
	group := &Group{}
	if err := json.Unmarshal(t.Data, group); err != nil {
		return nil, err
	}

	return group, nil
}
