package main

import (
	"context"
	"log"

	"github.com/getmelove/gorder2/internal/common/config"
	"github.com/getmelove/gorder2/internal/common/genproto/stockpb"
	"github.com/getmelove/gorder2/internal/common/server"
	"github.com/getmelove/gorder2/internal/stock/ports"
	"github.com/getmelove/gorder2/internal/stock/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	// 若没有读到服务配置则记录错误并退出
	if err := config.NewViperConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	serviceName := viper.Sub("stock").GetString("service-name")
	serverType := viper.Sub("stock").GetString("server-to-run")

	// 创建服务
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	application := service.NewApplication(ctx)
	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer(application)
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
	// TODO run http server
	default:
		log.Fatalf("Unsupported server type: %s", serverType)
	}
}
