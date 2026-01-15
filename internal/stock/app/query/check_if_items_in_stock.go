package query

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/decorator"
	"github.com/getmelove/gorder2/internal/stock/domain/stock"
	"github.com/getmelove/gorder2/internal/stock/entity"
	"github.com/getmelove/gorder2/internal/stock/infrastructure/integration"
	"github.com/sirupsen/logrus"
)

// 1.定义一个查询

type CheckIfItemsInStock struct {
	ItemsWithQuantity []*entity.ItemWithQuantity
}

//type CheckIfItemsInStockResponse struct {
//	InStock int
//	Items   []*orderpb.Item
//}

type CheckIfItemsInStockHandler decorator.QueryHandler[CheckIfItemsInStock, []*entity.Item]

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

func (c checkIfItemsInStockHandler) Handle(ctx context.Context, q CheckIfItemsInStock) ([]*entity.Item, error) {
	if err := c.checkStock(ctx, q.ItemsWithQuantity); err != nil {
		return nil, err
	}

	var items []*entity.Item
	for _, item := range q.ItemsWithQuantity {
		// TODO: 改为从数据库 or stripe获取
		priceID, err := c.stripeAPI.GetPriceByProductID(ctx, item.ID)
		if err != nil {
			logrus.Warnf("failed to get price by product id: %v", err)
			return nil, err
		}
		items = append(items, &entity.Item{
			ID:       item.ID,
			Name:     "",
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	return items, nil
}

func (c checkIfItemsInStockHandler) checkStock(ctx context.Context, query []*entity.ItemWithQuantity) error {
	var ids []string
	for _, i := range query {
		ids = append(ids, i.ID)
	}
	records, err := c.stockRepo.GetStock(ctx, ids)
	if err != nil {
		return err
	}
	idQuantityMap := make(map[string]int32)
	for _, r := range records {
		idQuantityMap[r.ID] += r.Quantity
	}
	var (
		ok       = true
		failedOn []struct {
			ID   string
			Want int32
			Have int32
		}
	)
	for _, item := range query {
		if item.Quantity > idQuantityMap[item.ID] {
			ok = false
			failedOn = append(failedOn, struct {
				ID   string
				Want int32
				Have int32
			}{ID: item.ID, Want: item.Quantity, Have: idQuantityMap[item.ID]})
		}
	}
	if ok {
		return nil
	}
	return stock.ExceedStockError{FailedOn: failedOn}
}
