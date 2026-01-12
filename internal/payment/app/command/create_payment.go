package command

import (
	"context"
	"errors"

	"github.com/getmelove/gorder2/internal/common/decorator"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/payment/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type CreatePayment struct {
	Order *orderpb.Order
}

type CreatePaymentHandler decorator.QueryHandler[CreatePayment, string]

type createPaymentCommand struct {
	// 创建订单的第三方服务
	// 通过stripe提供
	processor domain.Processor
	orderGRPC OrderService
}

func NewCreatePaymentHandler(processor domain.Processor, orderGRPC OrderService, logger *logrus.Entry, metrics decorator.MetricsClient) CreatePaymentHandler {
	if processor == nil {
		panic("nil processor")
	}
	if orderGRPC == nil {
		panic("nil orderGRPC")
	}
	return decorator.ApplyCommandDecorators[CreatePayment, string](createPaymentCommand{
		processor: processor,
		orderGRPC: orderGRPC,
	},
		logger,
		metrics,
	)
}

func (c createPaymentCommand) Handle(ctx context.Context, cmd CreatePayment) (string, error) {
	// test DLQ
	return "test-retry-link", errors.New("test retry")
	//
	tracer := otel.Tracer("payment")
	ctx, span := tracer.Start(ctx, "create_payment")
	defer span.End()
	//

	link, err := c.processor.CreatePaymentLink(ctx, cmd.Order)
	if err != nil {
		return "", err
	}
	logrus.Infof("create payment link for orderID=%s success, payment link=%s", cmd.Order, link)
	newOrder := &orderpb.Order{
		ID:          cmd.Order.ID,
		CustomerID:  cmd.Order.CustomerID,
		Status:      "waiting_for_payment",
		Items:       cmd.Order.Items,
		PaymentLink: link,
	}
	err = c.orderGRPC.UpdateOrder(ctx, newOrder)
	return link, err
}
