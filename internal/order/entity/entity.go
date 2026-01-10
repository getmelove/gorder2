package entity

type Order struct {
	CustomerID  string  `json:"customerID,omitempty"`
	Id          string  `json:"id,omitempty"`
	Items       []*Item `json:"items,omitempty"`
	PaymentLink string  `json:"paymentLink,omitempty"`
	Status      string  `json:"status,omitempty"`
}

type Item struct {
	ID       string
	Name     string
	Quantity int32
	PriceID  string
}

type ItemWithQuantity struct {
	ID       string
	Quantity int32
}
