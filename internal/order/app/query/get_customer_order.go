package query

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/decorator"
	domain "github.com/getmelove/gorder2/internal/order/domain/order"
	"github.com/sirupsen/logrus"
)

type GetCustomerOrder struct {
	CustomerId string `json:"customer_id"`
	OrderId    string `json:"order_id"`
}

type GetCustomerOrderHandler decorator.QueryHandler[GetCustomerOrder, *domain.Order]

type getCustomerOrderHandler struct {
	orderRepo domain.Repository
}

func NewGetCustomerOrderHandler(orderRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) GetCustomerOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyQueryDecotators[GetCustomerOrder, *domain.Order](
		getCustomerOrderHandler{orderRepo: orderRepo},
		logger,
		metricsClient,
	)
}

// 实现具体的查询
func (g getCustomerOrderHandler) Handle(ctx context.Context, q GetCustomerOrder) (*domain.Order, error) {
	o, err := g.orderRepo.Get(ctx, q.OrderId, q.CustomerId)
	if err != nil {
		return nil, err
	}
	return o, nil
}
