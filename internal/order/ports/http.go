package ports

import (
	"fmt"
	"log"
	"net/http"

	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/order/app"
	"github.com/getmelove/gorder2/internal/order/app/command"
	"github.com/getmelove/gorder2/internal/order/app/query"
	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	app app.Application
}

func NewHTTPServer(app app.Application) *HTTPServer {
	return &HTTPServer{app: app}
}

func (H HTTPServer) PostCustomerCustomerIDOrders(c *gin.Context, customerID string) {
	var req orderpb.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	r, err := H.app.Commands.CreateOrderHandler.Handle(c, command.CreateOrder{
		CustomerId: req.CustomerID,
		Items:      req.Items,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"customer_id":  req.CustomerID,
		"order_id":     r.OrderId,
		"redirect_url": fmt.Sprintf("http://localhost:8282/success?cutomerID=%s&orderID=%s", req.CustomerID, r.OrderId),
	})
}

func (H HTTPServer) GetCustomerCustomerIDOrdersOrderID(c *gin.Context, customerID string, orderID string) {
	//
	log.Println("HTTP Response GetCustomerCustomerIDOrdersOrderID")
	//
	o, err := H.app.Queries.GetCustomerOrderHandler.Handle(c, query.GetCustomerOrder{
		CustomerId: customerID,
		OrderId:    orderID,
	})
	if err != nil {
		c.JSON(200, gin.H{
			"error": err,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "sucsess",
			"data": gin.H{
				"Order": o,
			},
		})
	}
}
