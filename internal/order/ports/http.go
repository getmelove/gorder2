package ports

import (
	"log"

	"github.com/getmelove/gorder2/internal/order/app"
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
	//TODO implement me
	panic("implement me")
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
			"data":    o,
		})
	}
}
