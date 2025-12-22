package command

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/decorator"
	domain "github.com/getmelove/gorder2/internal/order/domain/order"
	"github.com/sirupsen/logrus"
)

// 1.定义一个cmd，也就是C。
type UpdateOrder struct {
	// 更新订单需要的信息
	// 客户的ID，以及订单的内容是什么，即客户下单了什么
	Order    *domain.Order `json:"order"`
	UpdateFn func(context.Context, *domain.Order) (*domain.Order, error)
}

type UpdateOrderHandler decorator.CommandHandler[UpdateOrder, interface{}]

type updateOrderHandler struct {
	orderRepo domain.Repository
	// stockGRPC
}

func NewUpdateOrderHandler(orderRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) UpdateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	return decorator.ApplyCommandDecorators[UpdateOrder, interface{}](
		updateOrderHandler{orderRepo: orderRepo},
		logger,
		metricsClient,
	)
}

func (c updateOrderHandler) Handle(ctx context.Context, cmd UpdateOrder) (interface{}, error) {
	// 进行兜底，没有实现UpdateFn的话就什么都不干
	if cmd.UpdateFn == nil {
		logrus.Warnf("UpdateOrder command is missing UpdateFn function, order=%#v", cmd.Order)
		cmd.UpdateFn = func(_ context.Context, order *domain.Order) (*domain.Order, error) {
			return order, nil
		}
	}
	err := c.orderRepo.Update(ctx, cmd.Order, cmd.UpdateFn)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
