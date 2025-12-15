package service

import (
	"context"

	"github.com/getmelove/gorder2/internal/stock/app"
)

// 将app包中的服务创建出来
func NewApplication(ctx context.Context) app.Application {
	return app.Application{}
}
