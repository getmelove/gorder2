package query

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
)

// 定一下和stock的通信协议
type StockService interface {
	CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) error
	GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error)
}
