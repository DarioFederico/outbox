package models

import "time"

type DeliveryStatus string

const (
	OutboxStatus_Pending   DeliveryStatus = "pending"
	OutboxStatus_OnProcess DeliveryStatus = "onProcess"
	OutboxStatus_Processed DeliveryStatus = "processed"
)

type Outbox struct {
	ID        int64          `json:"id"`
	Type      string         `json:"type"`
	Message   string         `json:"message"`
	Status    DeliveryStatus `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
