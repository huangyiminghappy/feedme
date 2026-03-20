package queue

import (
	"container/heap"
	"mcdonalds-bot-backend/internal/models"
	"sync"
)

type OrderItem struct {
	Order    *models.Order
	Priority int
	Index    int
}

type PriorityQueue []*OrderItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].Priority != pq[j].Priority {
		return pq[i].Priority > pq[j].Priority
	}
	return pq[i].Order.CreatedAt.Before(pq[j].Order.CreatedAt)
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*OrderItem)
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

type OrderQueue struct {
	pq    PriorityQueue
	mutex sync.RWMutex
}

func NewOrderQueue() *OrderQueue {
	oq := &OrderQueue{
		pq: make(PriorityQueue, 0),
	}
	heap.Init(&oq.pq)
	return oq
}

func (oq *OrderQueue) Push(order *models.Order) {
	oq.mutex.Lock()
	defer oq.mutex.Unlock()

	item := &OrderItem{
		Order:    order,
		Priority: order.Priority(),
	}
	heap.Push(&oq.pq, item)
}

func (oq *OrderQueue) Pop() *models.Order {
	oq.mutex.Lock()
	defer oq.mutex.Unlock()

	if oq.pq.Len() == 0 {
		return nil
	}

	item := heap.Pop(&oq.pq).(*OrderItem)
	return item.Order
}

func (oq *OrderQueue) Peek() *models.Order {
	oq.mutex.RLock()
	defer oq.mutex.RUnlock()

	if oq.pq.Len() == 0 {
		return nil
	}

	return oq.pq[0].Order
}

func (oq *OrderQueue) Len() int {
	oq.mutex.RLock()
	defer oq.mutex.RUnlock()
	return oq.pq.Len()
}

func (oq *OrderQueue) GetAll() []*models.Order {
	oq.mutex.RLock()
	defer oq.mutex.RUnlock()

	if len(oq.pq) == 0 {
		return []*models.Order{}
	}

	tempPQ := make(PriorityQueue, len(oq.pq))
	for i, item := range oq.pq {
		tempPQ[i] = &OrderItem{
			Order:    item.Order,
			Priority: item.Priority,
			Index:    i,
		}
	}
	heap.Init(&tempPQ)

	orders := make([]*models.Order, 0, len(tempPQ))
	for tempPQ.Len() > 0 {
		item := heap.Pop(&tempPQ).(*OrderItem)
		orders = append(orders, item.Order)
	}

	return orders
}

func (oq *OrderQueue) Remove(orderID string) bool {
	oq.mutex.Lock()
	defer oq.mutex.Unlock()

	for i, item := range oq.pq {
		if item.Order.ID == orderID {
			heap.Remove(&oq.pq, i)
			return true
		}
	}
	return false
}
