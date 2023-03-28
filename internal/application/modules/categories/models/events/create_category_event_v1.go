package events

import "time"

type CreateCategoryEventV1 struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (v *CreateCategoryEventV1) GetEventName() CategoryEventName {
	return CreateCategoryEventV1_Name
}
