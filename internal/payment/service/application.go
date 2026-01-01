package service

import (
	"context"

	grpcClient "github.com/getmelove/gorder2/internal/common/client/order_grpc"
	"github.com/getmelove/gorder2/internal/common/metrics"
	"github.com/getmelove/gorder2/internal/payment/app/command"
	"github.com/getmelove/gorder2/internal/payment/domain"
	"github.com/getmelove/gorder2/internal/payment/infrastructure/processor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/getmelove/gorder2/internal/payment/adapters"
	"github.com/getmelove/gorder2/internal/payment/app"
)

// 胶水层，将之前的这些抽象全部粘在一起，返回给业务层使用
// 将app包中的服务创建出来发给main使用
func NewApplication(ctx context.Context) (app.Application, func()) {
	orderClient, closeOrderClient, err := grpcClient.NewOrderGRPCClient(ctx)
	if err != nil {
		panic(err)
	}
	//memoryProcessor := processor.NewInmemProcessor()
	stripeProcessor := processor.NewStripeProcessor(viper.GetString("stripe-key"))
	orderGRPC := adapters.NewOrderGRPC(orderClient)
	return newApplication(ctx, orderGRPC, stripeProcessor), func() {
		_ = closeOrderClient()
	}
}

func newApplication(ctx context.Context, orderGRPC command.OrderService, memoryProcessor domain.Processor) app.Application {
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NewTodoMetrics()
	return app.Application{
		Commands: app.Commands{
			CreatePayment: command.NewCreatePaymentHandler(memoryProcessor, orderGRPC, logger, metricsClient),
		},
	}
}
