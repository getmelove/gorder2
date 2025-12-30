package domain

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
)

type Processor interface {
	CreatePaymentLink(ctx context.Context, order *orderpb.Order) (string, error)
}
