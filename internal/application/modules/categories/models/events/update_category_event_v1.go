package events

import "time"

type UpdateCategoryEventV1 struct {
	ID        int64             `json:"id"`
	Changes   map[string]string `json:"changes"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func (v *UpdateCategoryEventV1) GetEventName() CategoryEventName {
	return UpdateCategoryEventV1_Name
}
