package models

import "time"

type BotStatus string

const (
	BotStatusIdle       BotStatus = "Idle"
	BotStatusProcessing BotStatus = "Processing"
)

type Bot struct {
	ID            string    `json:"id"`
	Status        BotStatus `json:"status"`
	CurrentOrder  *Order    `json:"current_order,omitempty"`
	ProcessingTime int      `json:"processing_time"`
	CompletedOrders int     `json:"completed_orders"`
}

func NewBot(id string) *Bot {
	return &Bot{
		ID:              id,
		Status:          BotStatusIdle,
		ProcessingTime:  10,
		CompletedOrders: 0,
	}
}

func (b *Bot) IsIdle() bool {
	return b.Status == BotStatusIdle
}

func (b *Bot) StartProcessing(order *Order) {
	b.Status = BotStatusProcessing
	b.CurrentOrder = order
	order.Status = OrderStatusProcessing
	order.BotID = b.ID
	now := time.Now()
	order.StartedAt = &now
}

func (b *Bot) CompleteProcessing() *Order {
	if b.CurrentOrder == nil {
		return nil
	}
	
	order := b.CurrentOrder
	order.Status = OrderStatusComplete
	now := time.Now()
	order.CompletedAt = &now
	
	b.Status = BotStatusIdle
	b.CurrentOrder = nil
	b.CompletedOrders++
	
	return order
}
