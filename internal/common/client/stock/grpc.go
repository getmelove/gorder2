package client

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/genproto/stockpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 使用grpc和stock服务通信的可以复用下面的代码，实现和stock grpc的链接
func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	grpcAddr := viper.Sub("stock").GetString("grpc-addr")
	opts, err := grpcDialOpts(grpcAddr)
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	conn, err := grpc.NewClient(grpcAddr, opts...)
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	return stockpb.NewStockServiceClient(conn), conn.Close, nil
}

func grpcDialOpts(addr string) ([]grpc.DialOption, error) {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}, nil
}
