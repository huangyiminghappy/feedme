package models

import (
	"time"
)

type OrderType string

const (
	OrderTypeNormal OrderType = "Normal"
	OrderTypeVIP    OrderType = "VIP"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "Pending"
	OrderStatusProcessing OrderStatus = "Processing"
	OrderStatusComplete   OrderStatus = "Complete"
)

type Order struct {
	ID        string      `json:"id"`
	Type      OrderType   `json:"type"`
	Status    OrderStatus `json:"status"`
	BotID     string      `json:"bot_id,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	StartedAt *time.Time  `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

func NewOrder(id string, orderType OrderType) *Order {
	return &Order{
		ID:        id,
		Type:      orderType,
		Status:    OrderStatusPending,
		CreatedAt: time.Now(),
	}
}

func (o *Order) IsVIP() bool {
	return o.Type == OrderTypeVIP
}

func (o *Order) Priority() int {
	if o.IsVIP() {
		return 1
	}
	return 0
}
