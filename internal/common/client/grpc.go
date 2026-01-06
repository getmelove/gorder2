package client

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/getmelove/gorder2/internal/common/discovery"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 使用grpc和order服务通信的可以复用下面的代码，实现和order grpc的链接

func NewOrderGRPCClient(ctx context.Context) (client orderpb.OrderServiceClient, close func() error, err error) {
	if !waitForOrderGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("order grpc not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("order.service-name"))
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	if grpcAddr == "" {
		logrus.Warn("no order grpc service found")
	}
	opts := grpcDialOpts(grpcAddr)
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
	return orderpb.NewOrderServiceClient(conn), conn.Close, nil
}

func NewStockGRPCClient(ctx context.Context) (client stockpb.StockServiceClient, close func() error, err error) {
	if !waitForStockGRPCClient(viper.GetDuration("dial-grpc-timeout") * time.Second) {
		return nil, nil, errors.New("stock grpc not available")
	}
	grpcAddr, err := discovery.GetServiceAddr(ctx, viper.GetString("stock.service-name"))
	if err != nil {
		return nil, func() error {
			return nil
		}, err
	}
	if grpcAddr == "" {
		logrus.Warn("no stock grpc service found")
	}
	opts := grpcDialOpts(grpcAddr)
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

func grpcDialOpts(_ string) []grpc.DialOption {
	return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

func waitForOrderGRPCClient(timeout time.Duration) bool {
	logrus.Infof("waiting for order grpc client connection timeout %v", timeout)
	return waitFor(viper.GetString("order.grpc-addr"), timeout)
}

func waitForStockGRPCClient(timeout time.Duration) bool {
	logrus.Infof("waiting for order grpc client connection timeout %v", timeout)
	return waitFor(viper.GetString("stock.grpc-addr"), timeout)
}

func waitFor(addr string, timeout time.Duration) bool {
	portAvailable := make(chan struct{})
	timeoutCh := time.After(timeout)

	go func() {
		for {
			select {
			case <-timeoutCh:
				return
			default:
			}
			_, err := net.Dial("tcp", addr)
			if err == nil {
				portAvailable <- struct{}{}
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	select {
	case <-portAvailable:
		return true
	case <-timeoutCh:
		return false
	}
}
