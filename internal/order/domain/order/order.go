package order

import (
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
)

type Order struct {
	CustomerID  string          `json:"customerID,omitempty"`
	Id          string          `json:"id,omitempty"`
	Items       *[]orderpb.Item `json:"items,omitempty"`
	PaymentLink string          `json:"paymentLink,omitempty"`
	Status      string          `json:"status,omitempty"`
}
