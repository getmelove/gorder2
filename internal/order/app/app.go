package app

import (
	"github.com/getmelove/gorder2/internal/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	GetCustomerOrderHandler query.GetCustomerOrderHandler
}
