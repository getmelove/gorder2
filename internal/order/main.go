package main

import (
	"log"

	"github.com/getmelove/gorder2/common/config"
	"github.com/getmelove/gorder2/common/genproto/orderpb"
	"github.com/getmelove/gorder2/common/server"
	"github.com/getmelove/gorder2/order/ports"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// 初始化，读取服务配置
func init() {
	// 若没有读到服务配置则记录错误并退出
	if err := config.NewViperConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}

func main() {
	serviceName := viper.Sub("order").GetString("service-name")
	if serviceName == "" {
		log.Fatalf("Order service name is empty")
	}
	// serverType := viper.Sub("order").GetString("server-to-run")

	go server.RunGRPCServer(serviceName, func(server *grpc.Server) {
		svc := ports.NewGRPCServer()
		orderpb.RegisterOrderServiceServer(server, svc)
	})

	server.RunHttpServer(serviceName, func(router *gin.Engine) {
		ports.RegisterHandlersWithOptions(router, ports.NewHTTPServer(), ports.GinServerOptions{
			BaseURL:      "/api",
			Middlewares:  nil,
			ErrorHandler: nil,
		})
	})

}
