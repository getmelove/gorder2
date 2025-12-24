package ports

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/genproto/stockpb"
	"github.com/getmelove/gorder2/internal/stock/app"
	"github.com/getmelove/gorder2/internal/stock/app/query"
	"github.com/sirupsen/logrus"
)

type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) GetItems(ctx context.Context, request *stockpb.GetItemsRequest) (*stockpb.GetItemsResponse, error) {
	items, err := G.app.Queries.GetItemsHandler.Handle(ctx, query.GetItems{
		ItemsID: request.ItemIDs,
	})
	if err != nil {
		logrus.Warnf("GetItems: %v", err)
		return nil, err
	}
	return &stockpb.GetItemsResponse{Items: items}, nil
}

func (G GRPCServer) CheckIfItemsInStock(ctx context.Context, request *stockpb.CheckIfItemsInStockRequest) (*stockpb.CheckIfItemsInStockResponse, error) {
	items, err := G.app.Queries.CheckIfItemsInStockHandler.Handle(ctx, query.CheckIfItemsInStock{
		ItemsWithQuantity: request.Items,
	})
	if err != nil {
		logrus.Warnf("CheckIfItemsInStock: %v", err)
		return nil, err
	}
	return &stockpb.CheckIfItemsInStockResponse{
		InStock: 1,
		Items:   items,
	}, nil
}
