package grpc

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/common/genproto/stockpb"
	"github.com/sirupsen/logrus"
)

// 和stock通信的底层协议
type StockGrpc struct {
	client stockpb.StockServiceClient
}

func NewStockGrpc(client stockpb.StockServiceClient) *StockGrpc {
	return &StockGrpc{client: client}
}

func (s StockGrpc) CheckIfItemsInStock(ctx context.Context, items []*orderpb.ItemWithQuantity) (*stockpb.CheckIfItemsInStockResponse, error) {
	resp, err := s.client.CheckIfItemsInStock(ctx, &stockpb.CheckIfItemsInStockRequest{
		Items: items,
	})
	logrus.Info("stock grpc CheckIfItemsInStock resp:", resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s StockGrpc) GetItems(ctx context.Context, itemIDs []string) ([]*orderpb.Item, error) {
	resp, err := s.client.GetItems(ctx, &stockpb.GetItemsRequest{
		ItemIDs: itemIDs,
	})
	if err != nil {
		return nil, err
	}
	logrus.Info("stock grpc GetItems resp:", resp)
	return resp.Items, nil
}
