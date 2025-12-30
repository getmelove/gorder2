package command

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
)

// 所有的第三方服务都放在这里

type OrderService interface {
	UpdateOrder(ctx context.Context, order *orderpb.Order) error
}
