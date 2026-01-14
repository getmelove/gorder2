package ports

import (
	"fmt"

	"github.com/getmelove/gorder2/internal/common"
	client "github.com/getmelove/gorder2/internal/common/client/order"
	"github.com/getmelove/gorder2/internal/order/app"
	"github.com/getmelove/gorder2/internal/order/app/command"
	"github.com/getmelove/gorder2/internal/order/app/dto"
	"github.com/getmelove/gorder2/internal/order/app/query"
	"github.com/getmelove/gorder2/internal/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	common.BaseResponse
	app app.Application
}

func NewHTTPServer(app app.Application) *HTTPServer {
	return &HTTPServer{app: app}
}

func (H HTTPServer) PostCustomerCustomerIdOrders(c *gin.Context, customerID string) {
	//ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
	//defer span.End()
	//var req orderpb.CreateOrderRequest
	var (
		req  client.CreateOrderRequest
		err  error
		resp dto.CreateOrderResponse
	)
	defer func() {
		H.Response(c, err, &resp)
	}()
	if err = c.ShouldBindJSON(&req); err != nil {
		return
	}
	r, err := H.app.Commands.CreateOrderHandler.Handle(c.Request.Context(), command.CreateOrder{
		CustomerId: req.CustomerId,
		Items:      convertor.NewItemWithQuantityConvertor().ClientToEntities(req.Items),
	})
	if err != nil {
		return
	}
	resp = dto.CreateOrderResponse{
		CustomerId:  req.CustomerId,
		OrderId:     r.OrderId,
		RedirectURL: fmt.Sprintf("http://10.11.69.39:8282/success?customerID=%s&orderID=%s", req.CustomerId, r.OrderId),
	}
}

func (H HTTPServer) GetCustomerCustomerIdOrdersOrderId(c *gin.Context, customerID string, orderID string) {
	//ctx, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
	//defer span.End()
	var (
		err  error
		resp struct {
			Order interface{}
		}
	)
	defer func() {
		H.Response(c, err, resp)
	}()

	o, err := H.app.Queries.GetCustomerOrderHandler.Handle(c.Request.Context(), query.GetCustomerOrder{
		CustomerId: customerID,
		OrderId:    orderID,
	})
	resp.Order = convertor.NewOrderConvertor().EntityToClient(o)
}
