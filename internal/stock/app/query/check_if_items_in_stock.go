package query

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/decorator"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/stock/domain/stock"
	"github.com/sirupsen/logrus"
)

// 1.定义一个查询

type CheckIfItemsInStock struct {
	ItemsWithQuantity []*orderpb.ItemWithQuantity
}

//type CheckIfItemsInStockResponse struct {
//	InStock int
//	Items   []*orderpb.Item
//}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*orderpb.Item]

type checkIfItemsInStockHandler struct {
	stockRepo stock.Repository
}

func NewCheckIfItemsInStockHandler(stockRepo stock.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("stock repository is nil")
	}
	return decorator.ApplyQueryDecorators(
		checkIfItemsInStockHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)
}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, q CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var items []*orderpb.Item
	for _, item := range q.ItemsWithQuantity {
		items = append(items, &orderpb.Item{
			ID:       item.ID,
			Name:     "",
			Quantity: item.Quantity,
			PriceID:  "",
		})
	}
	return items, nil
}
