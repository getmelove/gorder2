package ports

import (
	"fmt"
	"net/http"

	client "github.com/getmelove/gorder2/internal/common/client/order"
	"github.com/getmelove/gorder2/internal/common/tracing"
	"github.com/getmelove/gorder2/internal/order/app"
	"github.com/getmelove/gorder2/internal/order/app/command"
	"github.com/getmelove/gorder2/internal/order/app/query"
	"github.com/getmelove/gorder2/internal/order/convertor"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app app.Application
}

func NewHTTPServer(app app.Application) *HTTPServer {
	return &HTTPServer{app: app}
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	ctx, span := tracing.Start(c, "PostCustomerCustomerIDOrders")
	defer span.End()
	//var req orderpb.CreateOrderRequest
	var req client.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := H.app.Commands.CreateOrderHandler.Handle(ctx, command.CreateOrder{
		CustomerId: req.CustomerID,
		Items:      convertor.NewItemWithQuantityConvertor().ClientToEntities(req.Items),
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	traceID := tracing.TraceID(ctx)
	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"trace_id":     traceID,
		"customer_id":  req.CustomerID,
		"order_id":     r.OrderId,
		"redirect_url": fmt.Sprintf("http://10.11.71.154:8282/success?customerID=%s&orderID=%s", req.CustomerID, r.OrderId),
	})
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	//
	ctx, span := tracing.Start(c, "GetCustomerCustomerIDOrdersOrderID")
	defer span.End()
	o, err := H.app.Queries.GetCustomerOrderHandler.Handle(ctx, query.GetCustomerOrder{
		CustomerId: customerID,
		OrderId:    orderID,
	})
	traceID := tracing.TraceID(ctx)
	if err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"message":  "sucsess",
			"trace_id": traceID,
			"data": gin.H{
				"Order": o,
			},
		})
	}
}
