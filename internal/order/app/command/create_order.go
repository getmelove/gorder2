package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/getmelove/gorder2/internal/common/broker"
	"github.com/getmelove/gorder2/internal/common/decorator"
	"github.com/getmelove/gorder2/internal/order/app/query"
	"github.com/getmelove/gorder2/internal/order/convertor"
	domain "github.com/getmelove/gorder2/internal/order/domain/order"
	"github.com/getmelove/gorder2/internal/order/entity"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

// 1.定义一个cmd，也就是C。
type CreateOrder struct {
	// 创建订单需要的信息
	// 客户的ID，已经订单的内容是什么，即客户下单了什么
	CustomerId string                    `json:"customer_id"` // 客户ID
	Items      []entity.ItemWithQuantity `json:"items"`       // 客户下单的东西，前端传回来的是商品和数量
}

// 2. 定义R
type CreateOrderResult struct {
	OrderId string `json:"order_id"`
}

type CreateOrderHandler decorator.CommandHandler[CreateOrder, *CreateOrderResult]

type createOrderHandler struct {
	orderRepo domain.Repository
	stockGRPC query.StockService
	channel   *amqp.Channel
}

func NewCreateOrderHandler(orderRepo domain.Repository, stockGRPC query.StockService, channel *amqp.Channel, logger *logrus.Entry, metricsClient decorator.MetricsClient) CreateOrderHandler {
	if orderRepo == nil {
		panic("orderRepo is nil")
	}
	if stockGRPC == nil {
		panic("sotckgRPC is nil")
	}
	if channel == nil {
		panic("Channel is nil")
	}
	return decorator.ApplyCommandDecorators[CreateOrder, *CreateOrderResult](
		createOrderHandler{orderRepo: orderRepo, stockGRPC: stockGRPC, channel: channel},
		logger,
		metricsClient,
	)
}

func (c createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) (*CreateOrderResult, error) {
	// 创建好订单以后，转发给MQ
	// 1.创建queue
	q, err := c.channel.QueueDeclare(broker.EventOrderCreate, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	// 创建span
	t := otel.Tracer("rabbitmq")
	ctx, span := t.Start(ctx, fmt.Sprintf("rabbitmq.%s.publish", q.Name))
	defer span.End()
	// 处理校验,查看仓储是否足够
	validItems, err := c.validata(ctx, cmd.Items)
	if err != nil {
		return nil, err
	}
	pendingOrder, err := domain.NewPendingOrder(cmd.CustomerId, validItems)
	if err != nil {
		return nil, err
	}
	o, err := c.orderRepo.Create(ctx, pendingOrder)
	if err != nil {
		return &CreateOrderResult{
			OrderId: o.Id,
		}, err
	}
	// 检查通过，发布到mq中
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		logrus.Error("marshall Order to queue error", err)
		return nil, err
	}
	header := broker.InjectRabbitMQHeaders(ctx)
	err = c.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalledOrder,
		Headers:      header,
	})
	if err != nil {
		logrus.Error("publish Order to queue error", err)
		return nil, err
	}
	return &CreateOrderResult{
		OrderId: o.Id,
	}, nil
}

// 验证用户下单的东西，是否仓库中还有
func (c createOrderHandler) validata(ctx context.Context, items []entity.ItemWithQuantity) ([]*entity.Item, error) {
	if len(items) == 0 {
		return nil, errors.New("must have at least one item")
	}
	items = packItems(items)
	resp, err := c.stockGRPC.CheckIfItemsInStock(ctx, convertor.NewItemWithQuantityConvertor().EntityToProtos(items))
	if err != nil {
		return nil, err
	}
	return convertor.NewItemConvertor().ProtoToEntitys(resp.GetItems()), nil
}

// 将用户下单的相同ID的商品合并起来
func packItems(items []entity.ItemWithQuantity) []entity.ItemWithQuantity {
	if len(items) == 0 {
		return items
	}
	packed := make([]entity.ItemWithQuantity, 0, len(items))
	indexByID := make(map[string]int, len(items))
	for _, item := range items {
		if idx, ok := indexByID[item.ID]; ok {
			packed[idx].Quantity += item.Quantity
			continue
		}
		packed = append(packed, item)
		indexByID[item.ID] = len(packed) - 1
	}
	return packed
}
