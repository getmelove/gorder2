package ports

import (
	"context"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/order/app"
	"github.com/getmelove/gorder2/internal/order/app/command"
	"github.com/getmelove/gorder2/internal/order/app/query"
	domain "github.com/getmelove/gorder2/internal/order/domain/order"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// domain层用于表示实际的业务逻辑
type GRPCServer struct {
	app app.Application
}

func NewGRPCServer(app app.Application) *GRPCServer {
	return &GRPCServer{app: app}
}

func (G GRPCServer) CreateOrder(ctx context.Context, request *orderpb.CreateOrderRequest) (*emptypb.Empty, error) {
	_, err := G.app.Commands.CreateOrderHandler.Handle(ctx, command.CreateOrder{
		CustomerId: request.CustomerID,
		Items:      request.Items,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (G GRPCServer) GetOrder(ctx context.Context, request *orderpb.GetOrderRequest) (*orderpb.Order, error) {
	o, err := G.app.Queries.GetCustomerOrderHandler.Handle(ctx, query.GetCustomerOrder{
		CustomerId: request.CustomerID,
		OrderId:    request.OrderID,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return o.ToProto(), nil
}

func (G GRPCServer) UpdateOrder(ctx context.Context, request *orderpb.Order) (*emptypb.Empty, error) {
	logrus.Infof("order_grpc || request_in || request=%+v", request)
	order, err := domain.NewOrder(request.CustomerID, request.ID, request.Items, request.PaymentLink, request.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// 这里就是更新逻辑
	updateFn := func(ctx context.Context, order *domain.Order) (*domain.Order, error) {
		// 表示彻底更换一个order
		return order, nil
	}
	_, err = G.app.Commands.UpdateOrderHandler.Handle(ctx, command.UpdateOrder{Order: order, UpdateFn: updateFn})
	return &emptypb.Empty{}, err
}
