package query

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/decorator"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/stock/domain/stock"
	"github.com/getmelove/gorder2/internal/stock/infrastructure/integration"
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
	stripeAPI *integration.StripeAPI
}

func NewCheckIfItemsInStockHandler(stockRepo stock.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient, stripeAPI *integration.StripeAPI) CheckIfItemsInStockHandler {
	if stockRepo == nil {
		panic("stock repository is nil")
	}
	if stripeAPI == nil {
		panic("stripe api is nil")
	}
	return decorator.ApplyQueryDecorators(
		checkIfItemsInStockHandler{stockRepo: stockRepo, stripeAPI: stripeAPI},
		logger,
		metricsClient,
	)
}

// TODO: 确定商品priceID
//var stub = map[string]string{
//	"1": "price_1SkKOKEC0C5AuFWmhfgOrcU8",
//	"2": "price_1SkmAVEC0C5AuFWmSz8WzrpO",
//}

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, q CheckIfItemsInStock) ([]*orderpb.Item, error) {
	var items []*orderpb.Item
	for _, item := range q.ItemsWithQuantity {
		// TODO: 改为从数据库 or stripe获取
		priceID, err := c.stripeAPI.GetPriceByProductID(ctx, item.ID)
		if err != nil {
			logrus.Warnf("failed to get price by product id: %v", err)
			continue
			return nil, err
		}
		items = append(items, &orderpb.Item{
			ID:       item.ID,
			Name:     "",
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	return items, nil
}

//func getStubPriceID(id string) string {
//	priceID, ok := stub[id]
//	if !ok {
//		priceID = stub["1"]
//	}
//	return priceID
//}
