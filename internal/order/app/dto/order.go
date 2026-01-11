package dto

type CreateOrderResponse struct {
	CustomerId  string `json:"customer_id"`
	OrderId     string `json:"order_id"`
	RedirectURL string `json:"redirect_url"`
}
