package main

import (
	"context"
	"log"

	"github.com/getmelove/gorder2/internal/common/broker"
	_ "github.com/getmelove/gorder2/internal/common/config"
	"github.com/getmelove/gorder2/internal/common/discovery"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/common/logging"
	"github.com/getmelove/gorder2/internal/common/server"
	"github.com/getmelove/gorder2/internal/common/tracing"
	"github.com/getmelove/gorder2/internal/order/infrastructure/consumer"
	"github.com/getmelove/gorder2/internal/order/ports"
	"github.com/getmelove/gorder2/internal/order/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// 初始化，读取服务配置
func init() {
	// 初始化日志
	logging.Init()
	// 若没有读到服务配置则记录错误并退出
	// if err := config.NewViperConfig(); err != nil {
	// 	log.Fatalf("Error reading config file, %s", err)
	// }
}

func main() {
	serviceName := viper.Sub("order").GetString("service-name")
	if serviceName == "" {
		log.Fatalf("Order service name is empty")
	}
	// serverType := viper.Sub("order").GetString("server-to-run")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 使用jaeger
	shutdown, err := tracing.InitJaegerProvider(viper.GetString("jaeger.url"), serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = shutdown(ctx)
	}()
	// 拿取service组装好的下层组件
	application, cleanup := service.NewApplication(ctx)
	defer cleanup()
	// 注册grpc服务
	logrus.Info("start register to consul")
	deregisterFunc, err := discovery.RegisterToConsul(ctx, serviceName)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		_ = deregisterFunc()
	}()
	logrus.Info("start register to consul end")

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
	// 起一个协程消费paid事件
	go consumer.NewConsumer(application).Listen(ch)

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer(application)
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHttpServer(serviceName, func(router *gin.Engine) {
		router.StaticFile("/success", "../../public/success.html")
		ports.RegisterHandlersWithOptions(router, ports.NewHTTPServer(application), ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
