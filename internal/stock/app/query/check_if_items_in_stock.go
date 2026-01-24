package query

import (
	"context"
	"strings"
	"time"

	"github.com/getmelove/gorder2/internal/common/decorator"
	"github.com/getmelove/gorder2/internal/common/handler/redis"
	"github.com/getmelove/gorder2/internal/stock/domain/stock"
	"github.com/getmelove/gorder2/internal/stock/entity"
	"github.com/getmelove/gorder2/internal/stock/infrastructure/integration"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	redisLockPrefix = "check_stock"
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
	// 开启事务
	// 1.上redis锁
	if err := lock(ctx, getLockKey(q)); err != nil {
		return nil, errors.Wrapf(err, "redis lock error: key=%s", getLockKey(q))
	}
	defer func() {
		if err := unlock(ctx, getLockKey(q)); err != nil {
			logrus.Warnf("redis unlock fail, err=%v", err)
		}
	}()
	//
	var items []*entity.Item
	// 2. 使用新封装的 CacheClient 进行价格缓存
	cacheCli := redis.DefaultCacheClient()

	for _, item := range q.ItemsWithQuantity {
		// 缓存 Key: price_id:{product_id}
		cacheKey := "price_id:" + item.ID
		// TTL: 24小时 (因为价格不常变动)
		ttl := 24 * time.Hour

		// 使用 GetOrSet 获取价格 ID
		// 如果缓存命中，直接返回；如果未命中，执行闭包里的逻辑去 Stripe 查询
		priceID, err := cacheCli.GetOrSet(ctx, cacheKey, ttl, func(ctx context.Context) (string, error) {
			pid, err := c.stripeAPI.GetPriceByProductID(ctx, item.ID)
			if err != nil {
				return "", err
			}
			return pid, nil
		})

		if err != nil {
			// 如果查询失败，记录日志并返回错误
			logrus.Warnf("failed to get price for product %s: %v", item.ID, err)
			return nil, err
		}

		items = append(items, &entity.Item{
			ID:       item.ID,
			Name:     "",
			Quantity: item.Quantity,
			PriceID:  priceID,
		})
	}
	// 扣减库存
	if err := c.checkStock(ctx, q.ItemsWithQuantity); err != nil {
		return nil, err
	}
	return items, nil
}

func getLockKey(query CheckIfItemsInStock) string {
	var ids []string
	for _, i := range query.ItemsWithQuantity {
		ids = append(ids, i.ID)
	}
	return redisLockPrefix + strings.Join(ids, "_")
}

func lock(ctx context.Context, key string) error {
	return redis.SetNX(ctx, redis.LocalClient(), key, "1", 5*time.Minute)
}

func unlock(ctx context.Context, key string) error {
	return redis.Del(ctx, redis.LocalClient(), key)
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
		return c.stockRepo.UpdateStock(ctx, query, func(
			ctx context.Context,
			existing []*entity.ItemWithQuantity,
			query []*entity.ItemWithQuantity,
		) ([]*entity.ItemWithQuantity, error) {
			var newItems []*entity.ItemWithQuantity
			for _, e := range existing {
				for _, q := range query {
					if e.ID == q.ID {
						newItems = append(newItems, &entity.ItemWithQuantity{
							ID:       e.ID,
							Quantity: e.Quantity - q.Quantity,
						})
					}
				}
			}
			return newItems, nil
		})
	}
	return stock.ExceedStockError{FailedOn: failedOn}
}

// func getStubPriceID(id string) string {
// 	priceID, ok := stub[id]
// 	if !ok {
// 		priceID = stub["1"]
// 	}
// 	return priceID
// }
