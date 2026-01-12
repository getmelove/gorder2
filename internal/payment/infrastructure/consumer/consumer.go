package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getmelove/gorder2/internal/common/broker"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/payment/app"
	"github.com/getmelove/gorder2/internal/payment/app/command"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type Consumer struct {
	app app.Application
}

func NewConsumer(app app.Application) *Consumer {
	return &Consumer{
		app: app,
	}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.EventOrderCreate, true, false, false, false, nil)
	if err != nil {
		logrus.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		logrus.Warnf("consume fail : queue=%s, err=%v", q.Name, err)
	}
	// 永久阻塞
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			c.handleMessage(ch, msg, q)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(ch *amqp.Channel, msg amqp.Delivery, q amqp.Queue) {
	//
	ctx := broker.ExtractRabbitMQHeaders(context.Background(), msg.Headers)
	t := otel.Tracer("rabbitmq")
	_, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.consume", q.Name))
	defer span.End()

	var err error
	defer func() {
		if err != nil {
			_ = msg.Nack(false, false)
		} else {
			_ = msg.Ack(false)
		}
	}()
	//
	logrus.Infof("Received payment message from queue=%s, msg=%s", q.Name, msg.Body)
	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("fail to unmarshal order, err=%+v", err)
		return
	}
	logrus.Info("Test retry, sleep for 5s, kill order now")
	time.Sleep(5 * time.Second)
	// 创建链接
	if _, err = c.app.Commands.CreatePayment.Handle(ctx, command.CreatePayment{

		Order: o,
	}); err != nil {
		// TODO: retry
		logrus.Infof("fail to create payment, err=%+v", err)
		if err = broker.HandleRetry(ctx, ch, &msg); err != nil {
			logrus.Warnf("retry_error, error handling retry, messageID = %s, err = %v", msg.MessageId, err)
		}
		return
	}
	span.AddEvent("payment.created")
	logrus.Info("consume success")
}
