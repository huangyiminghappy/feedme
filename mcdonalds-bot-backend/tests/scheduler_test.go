package tests

import (
	"mcdonalds-bot-backend/internal/models"
	"mcdonalds-bot-backend/internal/scheduler"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {
	s := scheduler.NewScheduler()

	err := s.CreateOrder("order1", models.OrderTypeNormal)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	orders := s.GetAllOrders()
	if len(orders) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(orders))
	}

	if orders[0].ID != "order1" {
		t.Fatalf("Expected order ID 'order1', got '%s'", orders[0].ID)
	}
}

func TestVIPPriority(t *testing.T) {
	s := scheduler.NewScheduler()

	s.CreateOrder("normal1", models.OrderTypeNormal)
	s.CreateOrder("vip1", models.OrderTypeVIP)
	s.CreateOrder("normal2", models.OrderTypeNormal)

	pending := s.GetPendingOrders()
	if len(pending) != 3 {
		t.Fatalf("Expected 3 pending orders, got %d", len(pending))
	}

	if pending[0].ID != "vip1" {
		t.Fatalf("Expected VIP order first, got '%s'", pending[0].ID)
	}
}

func TestBotProcessing(t *testing.T) {
	s := scheduler.NewScheduler()

	s.AddBot("bot1")
	s.CreateOrder("order1", models.OrderTypeNormal)

	time.Sleep(100 * time.Millisecond)

	bots := s.GetBots()
	if len(bots) != 1 {
		t.Fatalf("Expected 1 bot, got %d", len(bots))
	}

	if bots[0].Status != models.BotStatusProcessing {
		t.Fatalf("Expected bot to be processing, got status '%s'", bots[0].Status)
	}

	if bots[0].CurrentOrder == nil {
		t.Fatal("Expected bot to have current order")
	}

	if bots[0].CurrentOrder.ID != "order1" {
		t.Fatalf("Expected bot to process 'order1', got '%s'", bots[0].CurrentOrder.ID)
	}
}

func TestOrderCompletion(t *testing.T) {
	s := scheduler.NewScheduler()

	s.AddBot("bot1")
	s.CreateOrder("order1", models.OrderTypeNormal)

	time.Sleep(11 * time.Second)

	completed := s.GetCompletedOrders()
	if len(completed) != 1 {
		t.Fatalf("Expected 1 completed order, got %d", len(completed))
	}

	if completed[0].ID != "order1" {
		t.Fatalf("Expected completed order 'order1', got '%s'", completed[0].ID)
	}

	bots := s.GetBots()
	if bots[0].Status != models.BotStatusIdle {
		t.Fatalf("Expected bot to be idle after completion, got status '%s'", bots[0].Status)
	}
}

func TestRemoveBot(t *testing.T) {
	s := scheduler.NewScheduler()

	s.AddBot("bot1")
	s.CreateOrder("order1", models.OrderTypeNormal)

	time.Sleep(100 * time.Millisecond)

	err := s.RemoveBot("bot1")
	if err != nil {
		t.Fatalf("Failed to remove bot: %v", err)
	}

	bots := s.GetBots()
	if len(bots) != 0 {
		t.Fatalf("Expected 0 bots after removal, got %d", len(bots))
	}

	pending := s.GetPendingOrders()
	if len(pending) != 1 {
		t.Fatalf("Expected order to return to queue, got %d pending orders", len(pending))
	}

	if pending[0].ID != "order1" {
		t.Fatalf("Expected 'order1' back in queue, got '%s'", pending[0].ID)
	}
}

func TestMultipleBotsAndOrders(t *testing.T) {
	s := scheduler.NewScheduler()

	s.AddBot("bot1")
	s.AddBot("bot2")

	s.CreateOrder("order1", models.OrderTypeNormal)
	s.CreateOrder("order2", models.OrderTypeVIP)
	s.CreateOrder("order3", models.OrderTypeNormal)

	time.Sleep(100 * time.Millisecond)

	bots := s.GetBots()
	processingCount := 0
	for _, bot := range bots {
		if bot.Status == models.BotStatusProcessing {
			processingCount++
		}
	}

	if processingCount != 2 {
		t.Fatalf("Expected 2 bots processing, got %d", processingCount)
	}

	pending := s.GetPendingOrders()
	if len(pending) != 1 {
		t.Fatalf("Expected 1 pending order, got %d", len(pending))
	}
}

func TestVIPOrderInsertedCorrectly(t *testing.T) {
	s := scheduler.NewScheduler()

	s.CreateOrder("normal1", models.OrderTypeNormal)
	time.Sleep(1 * time.Millisecond)
	s.CreateOrder("vip1", models.OrderTypeVIP)
	time.Sleep(1 * time.Millisecond)
	s.CreateOrder("vip2", models.OrderTypeVIP)
	time.Sleep(1 * time.Millisecond)
	s.CreateOrder("normal2", models.OrderTypeNormal)
	time.Sleep(10 * time.Millisecond)

	pending := s.GetPendingOrders()

	if len(pending) != 4 {
		t.Fatalf("Expected 4 pending orders, got %d", len(pending))
	}

	t.Logf("Order 0: %s (%s)", pending[0].ID, pending[0].Type)
	t.Logf("Order 1: %s (%s)", pending[1].ID, pending[1].Type)
	t.Logf("Order 2: %s (%s)", pending[2].ID, pending[2].Type)
	t.Logf("Order 3: %s (%s)", pending[3].ID, pending[3].Type)

	if pending[0].Type != models.OrderTypeVIP {
		t.Fatalf("First order should be VIP, got %s (ID: %s)", pending[0].Type, pending[0].ID)
	}
	if pending[1].Type != models.OrderTypeVIP {
		t.Fatalf("Second order should be VIP, got %s (ID: %s)", pending[1].Type, pending[1].ID)
	}
	if pending[2].Type != models.OrderTypeNormal {
		t.Fatalf("Third order should be Normal, got %s (ID: %s)", pending[2].Type, pending[2].ID)
	}
	if pending[3].Type != models.OrderTypeNormal {
		t.Fatalf("Fourth order should be Normal, got %s (ID: %s)", pending[3].Type, pending[3].ID)
	}

	if pending[0].ID != "vip1" {
		t.Fatalf("First VIP order should be vip1, got %s", pending[0].ID)
	}
	if pending[1].ID != "vip2" {
		t.Fatalf("Second VIP order should be vip2, got %s", pending[1].ID)
	}
}
