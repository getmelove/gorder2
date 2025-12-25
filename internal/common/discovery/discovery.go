package discovery

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// 定义接口，什么东西才叫服务发现
type Registry interface {
	// 注册服务
	Register(ctx context.Context, instanceID, serviceName, hostPort string) error
	// 注销服务
	Deregister(ctx context.Context, instanceID, serviceName string) error
	// 发现服务
	Discover(ctx context.Context, serviceName string) ([]string, error)
	// 探活，看那些服务还健在
	HealthCheck(instanceID, serviceName string) error
}

func GenerateInstanceID(serviceName string) string {
	x := rand.New(rand.NewSource(time.Now().UnixNano())).Int()
	return fmt.Sprintf("%s_%d", serviceName, x)
}
