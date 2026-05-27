package storage

import (
	"sync"
	"time"
)

type Market struct {
	mu   sync.RWMutex
	Data map[string]*Spot
}

type SpotStatus string

const (
	StatusActive SpotStatus = "SPOT_STATUS_ACTIVE"

	StatusDisabled SpotStatus = "SPOT_STATUS_DISABLED"
)

type Spot struct {
	ID     string //уникальный идентификатор
	Symbol string //идентификатор спота (BTC/USDT)

	Base_asset  string //актив который продают или покупаеют
	Quote_asset string //в чем выражается цена

	Price_precision    int32  // количество знаков после запятой у цены
	Quantity_precision int32  // количество знаков после запятой у актива
	Min_order_size     string // минимальное возможное количество покупки/продажи
	Max_order_size     string // максимальное количество покупки актива

	Status       SpotStatus // статус спота
	AllowedRoles []string   //разрешенные роли у пользователей

	Name        string // название спота
	Description string // описание спота

	Created_at time.Time // время создания
	Updated_at time.Time // время обновления
	Deleted_at time.Time // время удаления
}

func NewMarketStorage() *Market {
	return &Market{Data: make(map[string]*Spot)}
}

func (r *Market) Create(
	spot *Spot,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Data[spot.ID] = spot
}

func (r *Market) Get(
	id string,
) (
	*Spot,
	bool,
) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	spot, ok := r.Data[id]
	if !ok {
		return &Spot{}, false
	}

	return spot, true
}

func (r *Market) Delete(
	id string,
) (
	bool,
	time.Time,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	spot, ok := r.Data[id]
	if !ok {
		return false, time.Time{}
	}

	deleted := time.Now()

	spot.Deleted_at = deleted
	spot.Status = StatusDisabled

	return true, deleted
}

func (r *Market) List() []*Spot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	spots := make([]*Spot, 0, len(r.Data))

	for _, spot := range r.Data {

		spots = append(spots, spot)

	}

	return spots
}
