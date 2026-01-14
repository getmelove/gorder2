package service

import (
	"context"
	"fmt"
	"time"

	"github.com/getmelove/gorder2/internal/common/broker"
	grpcClient "github.com/getmelove/gorder2/internal/common/client"
	"github.com/getmelove/gorder2/internal/common/metrics"
	"github.com/getmelove/gorder2/internal/order/adapters"
	"github.com/getmelove/gorder2/internal/order/adapters/grpc"
	"github.com/getmelove/gorder2/internal/order/app"
	"github.com/getmelove/gorder2/internal/order/app/command"
	"github.com/getmelove/gorder2/internal/order/app/query"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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
	// orderRepo := adapters.NewOrderInMemRepoIt()
	mongoClient := newMongoClient()
	orderRepo := adapters.NewOrderRepositoryMongo(mongoClient)
	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NewTodoMetrics()

	return app.Application{
		Commands: app.Commands{
			CreateOrderHandler: command.NewCreateOrderHandler(orderRepo, stockGRPC, ch, logger, metricsClient),
			UpdateOrderHandler: command.NewUpdateOrderHandler(orderRepo, logger, metricsClient),
		},
		Queries: app.Queries{
			GetCustomerOrderHandler: query.NewGetCustomerOrderHandler(orderRepo, logger, metricsClient),
		},
	}
}

func newMongoClient() *mongo.Client {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s",
		viper.GetString("mongo.user"),
		viper.GetString("mongo.password"),
		viper.GetString("mongo.host"),
		viper.GetString("mongo.port"),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err = c.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return c
}
