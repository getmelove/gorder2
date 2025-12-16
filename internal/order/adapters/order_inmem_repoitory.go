package adapters

import (
	"context"
	"strconv"
	"sync"
	"time"

	domain "github.com/getmelove/gorder2/internal/order/domain/order"
	"github.com/sirupsen/logrus"
)

type OrderInMemRepoIt struct {
	lock  *sync.RWMutex
	store []*domain.Order
}

func NewOrderInMemRepoIt() *OrderInMemRepoIt {
	return &OrderInMemRepoIt{
		lock:  &sync.RWMutex{},
		store: make([]*domain.Order, 0),
	}
}

func (o *OrderInMemRepoIt) Create(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// 创建订单，直接上锁
	o.lock.Lock()
	defer o.lock.Unlock()
	// 创建订单
	newOrder := &domain.Order{
		CustomerID:  order.CustomerID,
		Id:          strconv.FormatInt(time.Now().Unix(), 10),
		Items:       order.Items,
		PaymentLink: order.PaymentLink,
		Status:      order.Status,
	}
	o.store = append(o.store, newOrder)
	logrus.WithFields(logrus.Fields{
		"input_order":        order,
		"store_after_create": o.store,
	}).Debug("memory_order_repo_create")
	return newOrder, nil
}

func (o *OrderInMemRepoIt) Get(ctx context.Context, id, customerID string) (*domain.Order, error) {
	o.lock.RLock()
	defer o.lock.RUnlock()

	for _, order := range o.store {
		if order.Id == id && order.CustomerID == customerID {
			logrus.Debugf("OrderInMemRepoIt.Get||found||id=%s||customerID=%s||res=%v", id, customerID, *o)
			return order, nil
		}
	}
	return nil, domain.NotFoundError{OrderID: id}
}

func (o *OrderInMemRepoIt) Update(ctx context.Context, order *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	o.lock.Lock()
	defer o.lock.Unlock()
	found := false
	for i, od := range o.store {
		if od.Id == order.Id && od.CustomerID == order.CustomerID {
			found = true
			updateOrder, err := updateFn(ctx, od)
			if err != nil {
				return err
			}
			o.store[i] = updateOrder
		}
	}
	if !found {
		logrus.Debugf("OrderInMemRepoIt.Update||found=%v", found)
		return domain.NotFoundError{OrderID: order.Id}
	}
	return nil
}
