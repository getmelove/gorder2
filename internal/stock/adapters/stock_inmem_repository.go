package adapters

import (
	"context"
	"sync"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	domain "github.com/getmelove/gorder2/internal/stock/domain/stock"
)

type StockInMemRepoIt struct {
	lock  *sync.RWMutex
	store map[string]*orderpb.Item
}

var stub = map[string]*orderpb.Item{
	"item_id": {
		ID:       "foo_item",
		Name:     "stub item",
		Quantity: 10000,
		PriceID:  "stub_price",
	},
	"test-1": {
		ID:       "test-1",
		Name:     "stub item",
		Quantity: 10000,
		PriceID:  "stub_price",
	},
	"test-2": {
		ID:       "test-2",
		Name:     "stub item",
		Quantity: 10000,
		PriceID:  "test-2",
	},
}

func NewStockInMemRepoIt() *StockInMemRepoIt {
	return &StockInMemRepoIt{
		lock:  &sync.RWMutex{},
		store: stub,
	}
}

func (s StockInMemRepoIt) GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	res, missId := make([]*orderpb.Item, 0), make([]string, 0)
	for _, id := range ids {
		if item, exist := s.store[id]; exist {
			res = append(res, item)
		} else {
			missId = append(missId, id)
		}
	}
	if len(res) == len(ids) {
		return res, nil
	}
	return res, domain.NotFoundError{missId}
}
