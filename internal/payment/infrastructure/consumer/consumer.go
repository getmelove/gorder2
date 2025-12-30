package consumer

import (
	"context"
	"encoding/json"

	"github.com/getmelove/gorder2/internal/common/broker"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/payment/app"
	"github.com/getmelove/gorder2/internal/payment/app/command"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
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
			c.handleMessage(msg, q, ch)
		}
	}()
	<-forever
}

func (c *Consumer) handleMessage(msg amqp.Delivery, q amqp.Queue, ch *amqp.Channel) {
	logrus.Infof("Received payment message from queue=%s, msg=%s", q.Name, msg.Body)
	o := &orderpb.Order{}
	if err := json.Unmarshal(msg.Body, o); err != nil {
		logrus.Infof("fail to unmarshal order, err=%+v", err)
		_ = msg.Nack(false, false)
	}
	// 创建链接
	if _, err := c.app.Commands.CreatePayment.Handle(context.TODO(), command.CreatePayment{
		Order: o,
	}); err != nil {
		// TODO: retry
		logrus.Infof("fail to create payment, err=%+v", err)
		_ = msg.Nack(false, false)
		return
	}

	_ = msg.Ack(false)
	logrus.Info("consume success")
}
