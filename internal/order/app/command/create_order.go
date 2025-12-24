package command

import (
	"context"
	"errors"

	"github.com/getmelove/gorder2/internal/common/decorator"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/order/app/query"
	domain "github.com/getmelove/gorder2/internal/order/domain/order"
	"github.com/sirupsen/logrus"
)

// 1.定义一个cmd，也就是C。
type CreateOrder struct {
	// 创建订单需要的信息
	// 客户的ID，已经订单的内容是什么，即客户下单了什么
	CustomerId string                      `json:"customer_id"` // 客户ID
	Items      []*orderpb.ItemWithQuantity `json:"items"`       // 客户下单的东西，前端传回来的是商品和数量
}

// 2. 定义R
type CreateOrderResult struct {
	OrderId string `json:"order_id"`
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
}

func NewCreateOrderHandler(orderRepo domain.Repository, stockGRPC query.StockService, logger *logrus.Entry, metricsClient decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	if stockGRPC == nil {
		panic("sotckgRPC is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC},
		logger,
		metricsClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// 1.创建订单前，需要判断库存是否足够
	//err := c.stockGRPC.CheckIfItemsInStock(ctx, cmd.Items)
	//resp, err := c.stockGRPC.GetItems(ctx, []string{"123"})
	//logrus.Info("fail to create conn to stock grpc", resp)
	//
	//var stockResponse []*orderpb.Item
	//for _, item := range cmd.Items {
	//	stockResponse = append(stockResponse, &orderpb.Item{
	//		ID:       item.ID,
	//		Name:     "",
	//		Quantity: item.Quantity,
	//		PriceID:  "",
	//	})
	//}
	// 先处理校验
	validItems, err := c.validata(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	o, err := c.orderRepo.Create(ctx, &domain.Order{
		CustomerID:  cmd.CustomerId,
		Id:          "",
		Items:       validItems,
		PaymentLink: "",
		Status:      "",
	})
	if err != nil {
		return &CreateOrderResult{
			OrderId: o.Id,
		}, err
	}
	return &CreateOrderResult{
		OrderId: o.Id,
	}, nil
}

// 验证用户下单的东西，是否仓库中还有
func (c createOrderHandler) validata(ctx context.Context, items []*orderpb.ItemWithQuantity) ([]*orderpb.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	items = packItems(items)
	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, items)
	if err != nil {
		return nil, err
	}
	return resp.GetItems(), nil
}

// 将用户下单的相同ID的商品合并起来
func packItems(items []*orderpb.ItemWithQuantity) []*orderpb.ItemWithQuantity {
	if len(items) == 0 {
		return items
	}
	packed := make([]*orderpb.ItemWithQuantity, 0, len(items))
	indexByID := make(map[string]int, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		if idx, ok := indexByID[item.ID]; ok {
			packed[idx].Quantity += item.Quantity
			continue
		}
		clone := *item
		packed = append(packed, &clone)
		indexByID[item.ID] = len(packed) - 1
	}
	return packed
}
