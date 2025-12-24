package app

import "github.com/getmelove/gorder2/internal/stock/app/query"

// CQRS
type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
}

type Queries struct {
	GetItemsHandler            query.GetItemsHandler
	CheckIfItemsInStockHandler query.CheckIfItemsInStockHandler
}
