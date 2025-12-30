package main

import (
	"context"
	"log"

	"github.com/getmelove/gorder2/internal/common/broker"
	"github.com/getmelove/gorder2/internal/common/config"
	"github.com/getmelove/gorder2/internal/common/logging"
	"github.com/getmelove/gorder2/internal/common/server"
	"github.com/getmelove/gorder2/internal/payment/infrastructure/consumer"
	"github.com/getmelove/gorder2/internal/payment/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	// 初始化日志
	logging.Init()
	// 若没有读到服务配置则记录错误并退出
	if err := config.NewViperConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	serverType := viper.Sub("payment").GetString("server-to-run")
	application, cleanup := service.NewApplication(ctx)
	defer cleanup()
	// 初始化消息队列
	ch, closeCh := broker.Connect(
		viper.Sub("rabbitmq").GetString("user"),
		viper.Sub("rabbitmq").GetString("password"),
		viper.Sub("rabbitmq").GetString("host"),
		viper.Sub("rabbitmq").GetString("port"),
	)
	defer func() {
		_ = closeCh()
		_ = ch.Close()
	}()

	// 起一个协程，在后台实现不停地消费queue中的消息
	go consumer.NewConsumer(application).Listen(ch)
	paymentHandler := NewPaymentHandler()
	switch serverType {
	case "http":
		server.RunHttpServer(viper.Sub("payment").GetString("service-name"), paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("payment service not yet implemented grpc")
	default:
		logrus.Panic("unreachable code")
	}

}
