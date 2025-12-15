package main

import (
	"log"

	"github.com/getmelove/gorder2/common/genproto/stockpb"
	"github.com/getmelove/gorder2/common/server"
	"github.com/getmelove/gorder2/stock/ports"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	serviceName := viper.Sub("stock").GetString("service-name")
	serverType := viper.Sub("stock").GetString("server-to-run")

	switch serverType {
	case "grpc":
		server.RunGRPCServer(serviceName, func(server *grpc.Server) {
			svc := ports.NewGRPCServer()
			stockpb.RegisterStockServiceServer(server, svc)
		})
	case "http":
	// TODO run http server
	default:
		log.Fatalf("Unsupported server type: %s", serverType)
	}
}
