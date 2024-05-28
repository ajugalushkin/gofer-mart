package queue

import (
	"sync"
	"time"

	"github.com/ajugalushkin/gofer-mart/internal/dto"
)

type OrderQueue struct {
	queue sync.Map
}

func ClearExpiredResult(q *OrderQueue) {
	go func() {
		timer := time.NewTicker(time.Second)
		for {
			select {
			case <-timer.C:
				q.queue.Range(func(key, value interface{}) bool {
					q.queue.Delete(key)
					return true
				})
			}
		}
	}()
}

func NewOrderQueue() *OrderQueue {
	queue := &OrderQueue{}
	ClearExpiredResult(queue)
	return queue
}

var q = NewOrderQueue()

func (t *OrderQueue) Add(order *dto.Order) *dto.Order {
	t.queue.Store(order.Number, order)
	return order
}

func (t *OrderQueue) Fetch() (*dto.Order, bool) {
	if lenQueue(&t.queue) == 0 {
		return nil, false
	}

	order := &dto.Order{}
	t.queue.Range(func(key, value interface{}) bool {
		order = value.(*dto.Order)
		t.queue.Delete(key)
		return false
	})
	return order, false
}

func lenQueue(m *sync.Map) int {
	count := 0
	m.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

func AddOrder(order *dto.Order) *dto.Order {
	return q.Add(order)
}

func FetchOrder() (*dto.Order, bool) {
	return q.Fetch()
}
