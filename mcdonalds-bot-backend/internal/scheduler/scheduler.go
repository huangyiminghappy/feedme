package scheduler

import (
	"fmt"
	"mcdonalds-bot-backend/internal/models"
	"mcdonalds-bot-backend/internal/queue"
	"sync"
	"time"
)

type EventType string

const (
	EventOrderCreated    EventType = "order_created"
	EventOrderProcessing EventType = "order_processing"
	EventOrderComplete   EventType = "order_complete"
	EventBotAdded        EventType = "bot_added"
	EventBotRemoved      EventType = "bot_removed"
	EventQueueUpdated    EventType = "queue_updated"
)

type Event struct {
	Type      EventType     `json:"type"`
	Order     *models.Order `json:"order,omitempty"`
	Bot       *models.Bot   `json:"bot,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Message   string        `json:"message,omitempty"`
}

type EventListener func(event Event)

type Scheduler struct {
	bots             map[string]*models.Bot
	orderQueue       *queue.OrderQueue
	orders           map[string]*models.Order
	completedOrders  []*models.Order
	mutex            sync.RWMutex
	listeners        []EventListener
	stopChan         chan bool
	processingTimers map[string]*time.Timer
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		bots:             make(map[string]*models.Bot),
		orderQueue:       queue.NewOrderQueue(),
		orders:           make(map[string]*models.Order),
		completedOrders:  make([]*models.Order, 0),
		listeners:        make([]EventListener, 0),
		stopChan:         make(chan bool),
		processingTimers: make(map[string]*time.Timer),
	}
}

func (s *Scheduler) AddEventListener(listener EventListener) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.listeners = append(s.listeners, listener)
}

func (s *Scheduler) emitEvent(event Event) {
	s.mutex.RLock()
	if len(s.listeners) == 0 {
		s.mutex.RUnlock()
		return
	}
	listeners := make([]EventListener, len(s.listeners))
	copy(listeners, s.listeners)
	s.mutex.RUnlock()

	for _, listener := range listeners {
		go listener(event)
	}
}

func (s *Scheduler) AddBot(botID string) error {
	s.mutex.Lock()
	if _, exists := s.bots[botID]; exists {
		s.mutex.Unlock()
		return fmt.Errorf("bot %s already exists", botID)
	}

	bot := models.NewBot(botID)
	s.bots[botID] = bot
	s.mutex.Unlock()

	s.emitEvent(Event{
		Type:      EventBotAdded,
		Bot:       bot,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Bot %s added", botID),
	})

	go s.tryAssignOrder()

	return nil
}

func (s *Scheduler) RemoveBot(botID string) error {
	s.mutex.Lock()
	bot, exists := s.bots[botID]
	if !exists {
		s.mutex.Unlock()
		return fmt.Errorf("bot %s not found", botID)
	}

	if bot.CurrentOrder != nil {
		order := bot.CurrentOrder
		order.Status = models.OrderStatusPending
		order.BotID = ""
		order.StartedAt = nil
		s.orderQueue.Push(order)

		if timer, ok := s.processingTimers[order.ID]; ok {
			timer.Stop()
			delete(s.processingTimers, order.ID)
		}
	}

	delete(s.bots, botID)
	s.mutex.Unlock()

	s.emitEvent(Event{
		Type:      EventBotRemoved,
		Bot:       bot,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Bot %s removed", botID),
	})

	return nil
}

func (s *Scheduler) CreateOrder(orderID string, orderType models.OrderType) error {
	s.mutex.Lock()
	if _, exists := s.orders[orderID]; exists {
		s.mutex.Unlock()
		return fmt.Errorf("order %s already exists", orderID)
	}

	order := models.NewOrder(orderID, orderType)
	s.orders[orderID] = order
	s.orderQueue.Push(order)
	s.mutex.Unlock()

	s.emitEvent(Event{
		Type:      EventOrderCreated,
		Order:     order,
		Timestamp: time.Now(),
		Message:   fmt.Sprintf("Order %s created", orderID),
	})

	s.emitEvent(Event{
		Type:      EventQueueUpdated,
		Timestamp: time.Now(),
	})

	go s.tryAssignOrder()

	return nil
}

func (s *Scheduler) tryAssignOrder() {
	s.mutex.Lock()

	var assignedPairs []struct {
		bot   *models.Bot
		order *models.Order
	}

	for {
		idleBot := s.findIdleBot()
		if idleBot == nil {
			break
		}

		order := s.orderQueue.Pop()
		if order == nil {
			break
		}

		idleBot.StartProcessing(order)

		timer := time.AfterFunc(time.Duration(idleBot.ProcessingTime)*time.Second, func() {
			s.completeOrder(idleBot.ID)
		})
		s.processingTimers[order.ID] = timer

		assignedPairs = append(assignedPairs, struct {
			bot   *models.Bot
			order *models.Order
		}{bot: idleBot, order: order})
	}

	s.mutex.Unlock()

	for _, pair := range assignedPairs {
		s.emitEvent(Event{
			Type:      EventOrderProcessing,
			Order:     pair.order,
			Bot:       pair.bot,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("Order %s assigned to Bot %s", pair.order.ID, pair.bot.ID),
		})
	}

	if len(assignedPairs) > 0 {
		s.emitEvent(Event{
			Type:      EventQueueUpdated,
			Timestamp: time.Now(),
		})
	}
}

func (s *Scheduler) findIdleBot() *models.Bot {
	for _, bot := range s.bots {
		if bot.IsIdle() {
			return bot
		}
	}
	return nil
}

func (s *Scheduler) completeOrder(botID string) {
	s.mutex.Lock()
	bot, exists := s.bots[botID]
	if !exists || bot.CurrentOrder == nil {
		s.mutex.Unlock()
		return
	}

	order := bot.CompleteProcessing()
	if order != nil {
		s.completedOrders = append(s.completedOrders, order)
		delete(s.processingTimers, order.ID)
		s.mutex.Unlock()

		s.emitEvent(Event{
			Type:      EventOrderComplete,
			Order:     order,
			Bot:       bot,
			Timestamp: time.Now(),
			Message:   fmt.Sprintf("Order %s completed by Bot %s", order.ID, botID),
		})

		go s.tryAssignOrder()
	} else {
		s.mutex.Unlock()
	}
}

func (s *Scheduler) GetBots() []*models.Bot {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	bots := make([]*models.Bot, 0, len(s.bots))
	for _, bot := range s.bots {
		bots = append(bots, bot)
	}
	return bots
}

func (s *Scheduler) GetPendingOrders() []*models.Order {
	return s.orderQueue.GetAll()
}

func (s *Scheduler) GetCompletedOrders() []*models.Order {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	completed := make([]*models.Order, len(s.completedOrders))
	copy(completed, s.completedOrders)
	return completed
}

func (s *Scheduler) GetAllOrders() []*models.Order {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	orders := make([]*models.Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	return orders
}

func (s *Scheduler) GetStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	totalBots := len(s.bots)
	idleBots := 0
	processingBots := 0

	for _, bot := range s.bots {
		if bot.IsIdle() {
			idleBots++
		} else {
			processingBots++
		}
	}

	return map[string]interface{}{
		"total_bots":       totalBots,
		"idle_bots":        idleBots,
		"processing_bots":  processingBots,
		"pending_orders":   s.orderQueue.Len(),
		"completed_orders": len(s.completedOrders),
		"total_orders":     len(s.orders),
	}
}

func (s *Scheduler) Reset() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, timer := range s.processingTimers {
		timer.Stop()
	}

	s.bots = make(map[string]*models.Bot)
	s.orderQueue = queue.NewOrderQueue()
	s.orders = make(map[string]*models.Order)
	s.completedOrders = make([]*models.Order, 0)
	s.processingTimers = make(map[string]*time.Timer)
}
