package storage

import (
	"sync"
	"time"
)

type OrderData struct {
	mu   sync.RWMutex
	Data map[string]*Order
}

type OrderTypeStorage string

const (
	OrderTypeUnspecified OrderTypeStorage = "ORDER_TYPE_UNSPECIFIED"
	Buy                  OrderTypeStorage = "BUY"
	Sell                 OrderTypeStorage = "SELL"
)

type OrderStatusStorage string

const (
	OrderStatusUnspecified OrderStatusStorage = "ORDER_STATUS_UNSPECIFIED"
	Created                OrderStatusStorage = "CREATED"
	Canceled               OrderStatusStorage = "CANCELED"
	Rejected               OrderStatusStorage = "REJECTED"
	Completed              OrderStatusStorage = "COMPLETED"
)

type Order struct {
	ID     string
	UserId string
	SpotId string

	Type   OrderTypeStorage
	Status OrderStatusStorage

	Price    string
	Quantity string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewOrderStorage() *OrderData {
	return &OrderData{Data: make(map[string]*Order)}
}

func (r *OrderData) Set(
	order *Order,
) {

	r.mu.Lock()
	defer r.mu.Unlock()

	r.Data[order.ID] = order
}

func (r *OrderData) Get(
	id string,
) (
	*Order,
	bool,
) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.Data[id]
	if !ok {
		return nil, false
	}

	return order, true
}
