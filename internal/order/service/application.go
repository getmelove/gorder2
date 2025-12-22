package service

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/metrics"
	"github.com/getmelove/gorder2/internal/order/adapters"
	"github.com/getmelove/gorder2/internal/order/app"
	"github.com/getmelove/gorder2/internal/order/app/command"
	"github.com/getmelove/gorder2/internal/order/app/query"
	"github.com/sirupsen/logrus"
)

// 胶水层，将之前的这些抽象全部粘在一起，返回给业务层使用
// 将app包中的服务创建出来发给main使用
func NewApplication(ctx context.Context) app.Application {
	orderInmemRepo := adapters.NewOrderInMemRepoIt()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NewTodoMetrics()
	return app.Application{
		Commands: app.Commands{
			CreateOrderHandler: command.NewCreateOrderHandler(orderInmemRepo, logger, metricsClient),
			UpdateOrderHandler: command.NewUpdateOrderHandler(orderInmemRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrderHandler: query.NewGetCustomerOrderHandler(orderInmemRepo, logger, metricsClient),
		},
	}
}
