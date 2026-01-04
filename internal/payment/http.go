package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/getmelove/gorder2/internal/common/broker"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/payment/domain"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/webhook"
)

type PaymentHandler struct {
	ch *amqp.Channel
}

func NewPaymentHandler(ch *amqp.Channel) *PaymentHandler {
	return &PaymentHandler{ch: ch}
}

func (h *PaymentHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/api/webhook", h.handleWebhook)
}

func (h *PaymentHandler) handleWebhook(c *gin.Context) {
	logrus.Info("got webhook from stripe")
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Infof("error reading request body: %v\n", err)
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
		c.JSON(http.StatusServiceUnavailable, err.Error())
		return
	}

	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		viper.GetString("ENDPOINT_STRIPE_SECRET"))

	if err != nil {
		logrus.Infof("error constructing event: %v\n", err)
		//fmt.Fprintf(os.Stderr, "⚠️  Webhook signature verification failed. %v\n", err)
		//c.Writer.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	logrus.Debugf("event: %s\n", event.Type)
	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			logrus.Warnf("error unmarshaling checkout session: %v\n", err)
			return
		}
		// 说明支付成功了
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			logrus.Info("payment for checkout session %v success!", session.ID)
			ctx, cancel := context.WithCancel(context.TODO())
			defer cancel()

			var items []*orderpb.Item
			_ = json.Unmarshal([]byte(session.Metadata["items"]), &items)

			marshalledOrder, err := json.Marshal(&domain.Order{
				CustomerID:  session.Metadata["customerID"],
				Id:          session.Metadata["orderID"],
				Items:       items,
				PaymentLink: session.Metadata["PaymentLink"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
			})

			if err != nil {
				logrus.Warnf("error marshalling domain.order: %v\n", err)
				return
			}

			_ = h.ch.PublishWithContext(ctx, broker.EventOrderPaid, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         marshalledOrder,
			})
			logrus.Infof("message published to %s, body: %s", broker.EventOrderPaid, string(marshalledOrder))
		}
	default:
		logrus.Warnf("unknown event type: %v", event.Type)
	}

	c.JSON(http.StatusOK, nil)
}
