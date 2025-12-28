package main

import (
	"log"

	"github.com/getmelove/gorder2/internal/common/broker"
	"github.com/getmelove/gorder2/internal/common/config"
	"github.com/getmelove/gorder2/internal/common/logging"
	"github.com/getmelove/gorder2/internal/common/server"
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
	serverType := viper.Sub("payment").GetString("server-to-run")
	paymentHandler := NewPaymentHandler()
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

	switch serverType {
	case "http":
		server.RunHttpServer(viper.Sub("payment").GetString("service-name"), paymentHandler.RegisterRoutes)
	case "grpc":
		logrus.Panic("payment service not yet implemented grpc")
	default:
		logrus.Panic("unreachable code")
	}

}
