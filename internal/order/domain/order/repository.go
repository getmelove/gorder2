package order

import (
	"context"
	"fmt"
)

type Repository interface {
	Create(context.Context, *OrderAggregate) (*OrderAggregate, error)
	Get(ctx context.Context, id, customerID string) (*OrderAggregate, error)
	Update(
		ctx context.Context,
		o *OrderAggregate,
		updateFn func(context.Context, *OrderAggregate) (*OrderAggregate, error),
	) error
}
type NotFoundError struct {
	OrderID string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("order with id %s not found", e.OrderID)
}
