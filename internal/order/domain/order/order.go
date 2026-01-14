package order

import (
	"errors"
	"fmt"

	"github.com/getmelove/gorder2/internal/order/entity"
	"github.com/stripe/stripe-go/v84"
)

// Aggregate
type OrderAggregate struct {
	CustomerID  string         `json:"customerID,omitempty"`
	Id          string         `json:"id,omitempty"`
	Items       []*entity.Item `json:"items,omitempty"`
	PaymentLink string         `json:"paymentLink,omitempty"`
	Status      string         `json:"status,omitempty"`
}

func NewOrder(customerID string, id string, items []*entity.Item, paymentLink string, status string) (*OrderAggregate, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	if customerID == "" {
		return nil, errors.New("empty customerID")
	}
	if items == nil {
		return nil, errors.New("empty items")
	}
	if status == "" {
		return nil, errors.New("empty status")
	}
	// 订单刚创建的时候可以没有paymentLink
	return &OrderAggregate{
		CustomerID:  customerID,
		Id:          id,
		Items:       items,
		PaymentLink: paymentLink,
		Status:      status,
	}, nil
}

func (o *OrderAggregate) IsPaid() error {
	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
		return nil
	}
	return fmt.Errorf("order %s is not paid, status = %s", o.Id, o.Status)
}

func NewPendingOrder(customerID string, items []*entity.Item) (*OrderAggregate, error) {
	if customerID == "" {
		return nil, errors.New("empty customerID")
	}
	if items == nil {
		return nil, errors.New("empty items")
	}
	// 订单刚创建的时候可以没有paymentLink
	return &OrderAggregate{
		CustomerID: customerID,
		Items:      items,
		Status:     "pending",
	}, nil
}
