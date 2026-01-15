package query

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/decorator"
	domain "github.com/getmelove/gorder2/internal/stock/domain/stock"
	"github.com/getmelove/gorder2/internal/stock/entity"
	"github.com/sirupsen/logrus"
)

type GetItems struct {
	ItemsID []string `json:"items_id"`
}

type GetItemsHandler decorator.QueryHandler[GetItems, []*entity.Item]

type getItemsHandler struct {
	stockRepo domain.Repository
}

func NewGetItemsHandler(stockRepo domain.Repository, logger *logrus.Entry, metricsClient decorator.MetricsClient) GetItemsHandler {
	if stockRepo == nil {
		panic("stockRepo is nil")
	}
	return decorator.ApplyQueryDecorators[GetItems, []*entity.Item](
		getItemsHandler{stockRepo: stockRepo},
		logger,
		metricsClient,
	)

}

func (g getItemsHandler) Handle(ctx context.Context, q GetItems) ([]*entity.Item, error) {
	items, err := g.stockRepo.GetItems(ctx, q.ItemsID)
	if err != nil {
		return nil, err
	}
	return items, nil
}
