package stock

import (
	"context"
	"fmt"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
)

type Repository interface {
	GetItems(ctx context.Context, ids []string) ([]*orderpb.Item, error)
}

type NotFoundError struct {
	Missing []string
}

func (e NotFoundError) Error() string {

}
