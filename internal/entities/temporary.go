package entities

import "encoding/json"

type Temporary struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	Step      int    `json:"step"`
	MessageID *int   `json:"messageId,omitempty"`
	Data      []byte `json:"data"`
}

func (t *Temporary) GetGroup() (*Group, error) {
	group := &Group{}

	err := json.Unmarshal(t.Data, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (t *Temporary) GetPlant() (*Plant, error) {
	plant := &Plant{}

	err := json.Unmarshal(t.Data, plant)
	if err != nil {
		return nil, err
	}

	return plant, nil
}
