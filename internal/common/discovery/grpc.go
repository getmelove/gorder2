package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/getmelove/gorder2/internal/common/discovery/consul"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
)

func RegisterToConsul(ctx context.Context, serviceName string) (func() error, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return func() error {
			return nil
		}, err
	}
	instanceID := GenerateInstanceID(serviceName)
	grpcAddr := viper.Sub(serviceName).GetString("grpc-addr")
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		return func() error {
			return nil
		}, err
	}
	// 连接好后开一个协程监听连接的grpc服务是否还存在
	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				logrus.Panicf("no heartbeat from %s to registry, err=%v", serviceName, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	logrus.WithFields(logrus.Fields{
		"serviceName": serviceName,
		"addr":        grpcAddr,
	}).Info("register to consul")
	// return cancel method
	return func() error {
		return registry.Deregister(ctx, instanceID, serviceName)
	}, nil
}

func GetServiceAddr(ctx context.Context, serviceName string) (string, error) {
	registry, err := consul.New(viper.GetString("consul.addr"))
	if err != nil {
		return "", err
	}
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", fmt.Errorf("got empty %s addrs from consul", serviceName)
	}
	i := rand.Intn(len(addrs))
	logrus.Infof("discovered %d instance of %s, addr=%v)", len(addrs), serviceName, addrs)
	return addrs[i], nil
}
