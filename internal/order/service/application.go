package service

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/broker"
	grpcClient "github.com/getmelove/gorder2/internal/common/client/stock"
	"github.com/getmelove/gorder2/internal/common/metrics"
	"github.com/getmelove/gorder2/internal/order/adapters"
	"github.com/getmelove/gorder2/internal/order/adapters/grpc"
	"github.com/getmelove/gorder2/internal/order/app"
	"github.com/getmelove/gorder2/internal/order/app/command"
	"github.com/getmelove/gorder2/internal/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 胶水层，将之前的这些抽象全部粘在一起，返回给业务层使用
// 将app包中的服务创建出来发给main使用
func NewApplication(ctx context.Context) (app.Application, func()) {
	stockClient, closeStockClient, err := grpcClient.NewStockGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	stockGRPC := grpc.NewStockGrpc(stockClient)
	ch, closeCh := broker.Connect(
		viper.Sub("rabbitmq").GetString("user"),
		viper.Sub("rabbitmq").GetString("password"),
		viper.Sub("rabbitmq").GetString("host"),
		viper.Sub("rabbitmq").GetString("port"),
	)
	return newApplication(ctx, stockGRPC, ch), func() {
		_ = closeStockClient()
		_ = closeCh()
		_ = ch.Close()
	}
}

func newApplication(_ context.Context, stockGRPC query.StockService, ch *amqp.Channel) app.Application {
	orderInmemRepo := adapters.NewOrderInMemRepoIt()
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NewTodoMetrics()

	return app.Application{
		Commands: app.Commands{
			CreateOrderHandler: command.NewCreateOrderHandler(orderInmemRepo, stockGRPC, ch, logger, metricsClient),
			UpdateOrderHandler: command.NewUpdateOrderHandler(orderInmemRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrderHandler: query.NewGetCustomerOrderHandler(orderInmemRepo, logger, metricsClient),
		},
	}
}
