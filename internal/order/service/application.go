package service

import (
	"context"

	"github.com/getmelove/gorder2/internal/order/adapters"
	"github.com/getmelove/gorder2/internal/order/app"
)

// 将app包中的服务创建出来
func NewApplication(ctx context.Context) app.Application {
	orderRepo := adapters.NewOrderInMemRepoIt()
	return app.Application{}
}
