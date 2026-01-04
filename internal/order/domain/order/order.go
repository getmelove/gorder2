package order

import (
	"errors"
	"fmt"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/stripe/stripe-go/v84"
)

type Order struct {
	CustomerID  string          `json:"customerID,omitempty"`
	Id          string          `json:"id,omitempty"`
	Items       []*orderpb.Item `json:"items,omitempty"`
	PaymentLink string          `json:"paymentLink,omitempty"`
	Status      string          `json:"status,omitempty"`
}

func NewOrder(customerID string, id string, items []*orderpb.Item, paymentLink string, status string) (*Order, error) {
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
	return &Order{
		CustomerID:  customerID,
		Id:          id,
		Items:       items,
		PaymentLink: paymentLink,
		Status:      status,
	}, nil
}

func (o *Order) ToProto() *orderpb.Order {
	return &orderpb.Order{
		ID:          o.Id,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       o.Items,
		PaymentLink: o.PaymentLink,
	}
}

func (o *Order) IsPaid() error {
	if o.Status == string(stripe.CheckoutSessionPaymentStatusPaid) {
		return nil
	}
	return fmt.Errorf("order %s is not paid, status = %s", o.Id, o.Status)
}
