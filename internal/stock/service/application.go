package service

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/metrics"
	"github.com/getmelove/gorder2/internal/stock/adapters"
	"github.com/getmelove/gorder2/internal/stock/app"
	"github.com/getmelove/gorder2/internal/stock/app/query"
	"github.com/sirupsen/logrus"
)

// 将app包中的服务创建出来
func NewApplication(ctx context.Context) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NewTodoMetrics()
	stockRepo := adapters.NewStockInMemRepoIt()
	return app.Application{
		Commands: app.Commands{},
		Queries: app.Queries{
			GetItemsHandler:            query.NewGetItemsHandler(stockRepo, logger, metricsClient),
			CheckIfItemsInStockHandler: query.NewCheckIfItemsInStockHandler(stockRepo, logger, metricsClient),
		},
	}
}
