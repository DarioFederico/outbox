package models

import (
	"errors"
	"time"
)

var (
	NotFoundError error = errors.New("category not found")
)

type Category struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
