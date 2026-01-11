package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

const (
	DLX = "dlx"
	DLQ = "dlq"
)

func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	var err error
	defer func() {
		if err != nil {
			logrus.Fatal(err)
		}
	}()
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, func() error { return nil }
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, func() error { return nil }
	}
	// 声明两个exchange
	err = ch.ExchangeDeclare(EventOrderCreate, "direct", true, false, false, false, nil)
	if err != nil {
		return nil, func() error { return nil }
	}
	err = ch.ExchangeDeclare(EventOrderPaid, "fanout", true, false, false, false, nil)
	if err != nil {
		return nil, func() error { return nil }
	}
	if err = createDLX(ch); err != nil {
		return nil, func() error { return nil }
	}
	return ch, conn.Close
}

func createDLX(ch *amqp.Channel) error {

	q, err := ch.QueueDeclare("share_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(DLX, "fanout", true, false, false, false, nil)
	if err != nil {
		return err
	}
	err = ch.QueueBind(q.Name, "", DLX, false, nil)
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(DLQ, true, false, false, false, nil)

	return err
}

type RabbitMQHeaderCarrier map[string]interface{}

func (r RabbitMQHeaderCarrier) Get(key string) string {
	value, ok := r[key]
	if !ok {
		return ""
	}
	return value.(string)
}

func (r RabbitMQHeaderCarrier) Set(key, value string) {
	r[key] = value
}

func (r RabbitMQHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(r))
	i := 0
	for key := range r {
		keys[i] = key
		i++
	}
	return keys
}

func InjectRabbitMQHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(RabbitMQHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}

func ExtractRabbitMQHeaders(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, RabbitMQHeaderCarrier(headers))
}
