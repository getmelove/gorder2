package convertor

import (
	client "github.com/getmelove/gorder2/internal/common/client/order"
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	domain "github.com/getmelove/gorder2/internal/order/domain/order"
	"github.com/getmelove/gorder2/internal/order/entity"
)

type OrderConvertor struct{}
type ItemConvertor struct{}
type ItemWithQuantityConvertor struct{}

func (c *OrderConvertor) EntityToProto(o domain.OrderAggregate) *orderpb.Order {
	c.check(o)
	return &orderpb.Order{
		ID:          o.Id,
		CustomerID:  o.CustomerID,
		Status:      o.Status,
		Items:       NewItemConvertor().EntityToProtos(o.Items),
		PaymentLink: o.PaymentLink,
	}
}

func (c *OrderConvertor) ProtoToEntity(o *orderpb.Order) *domain.OrderAggregate {
	c.check(o)
	return &domain.OrderAggregate{
		CustomerID:  o.CustomerID,
		Id:          o.ID,
		Items:       NewItemConvertor().ProtoToEntitys(o.Items),
		PaymentLink: o.PaymentLink,
		Status:      o.Status,
	}
}

func (c *OrderConvertor) ClientToEntity(o *client.Order) *domain.OrderAggregate {
	c.check(o)
	return &domain.OrderAggregate{
		CustomerID:  o.CustomerID,
		Id:          o.Id,
		Items:       NewItemConvertor().ClientToEntitys(o.Items),
		PaymentLink: o.PaymentLink,
		Status:      o.Status,
	}
}

func (c *OrderConvertor) EntityToClient(o *domain.OrderAggregate) *client.Order {
	c.check(o)
	return &client.Order{
		CustomerID:  o.CustomerID,
		Id:          o.Id,
		Items:       NewItemConvertor().EntityToClients(o.Items),
		PaymentLink: o.PaymentLink,
		Status:      o.Status,
	}
}

func (c *OrderConvertor) check(aggregate interface{}) {
	if aggregate == nil {
		panic("aggregate is nil")
	}
}

func (c *ItemConvertor) EntityToProtos(items []*entity.Item) (res []*orderpb.Item) {
	for _, i := range items {
		res = append(res, c.EntityToProto(i))
	}
	return
}

func (c *ItemConvertor) ProtoToEntitys(items []*orderpb.Item) (res []*entity.Item) {
	for _, i := range items {
		res = append(res, c.ProtoToEntity(i))
	}
	return
}

func (c *ItemConvertor) ClientToEntitys(items []client.Item) (res []*entity.Item) {
	for _, i := range items {
		res = append(res, c.ClientToEntity(i))
	}
	return
}

func (c *ItemConvertor) EntityToClients(items []*entity.Item) (res []client.Item) {
	for _, i := range items {
		res = append(res, c.EntityToClient(i))
	}
	return
}

func (c *ItemConvertor) EntityToProto(i *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       i.ID,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceID,
	}
}

func (c *ItemConvertor) ProtoToEntity(i *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       i.ID,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceID,
	}
}

func (c *ItemConvertor) ClientToEntity(i client.Item) *entity.Item {
	return &entity.Item{
		ID:       i.Id,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceID,
	}

}

func (c *ItemConvertor) EntityToClient(i *entity.Item) client.Item {
	return client.Item{
		Id:       i.ID,
		Name:     i.Name,
		PriceID:  i.PriceID,
		Quantity: i.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) EntityToProtos(items []entity.ItemWithQuantity) (res []*orderpb.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.EntityToProto(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) ProtoToEntities(items []*orderpb.ItemWithQuantity) (res []entity.ItemWithQuantity) {
	for _, i := range items {
		if i == nil {
			continue
		}
		res = append(res, c.ProtoToEntity(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) EntityToClients(items []entity.ItemWithQuantity) (res []client.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.EntityToClient(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) EntityToProto(i entity.ItemWithQuantity) *orderpb.ItemWithQuantity {
	return &orderpb.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ProtoToEntity(i *orderpb.ItemWithQuantity) entity.ItemWithQuantity {
	return entity.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) EntityToClient(i entity.ItemWithQuantity) client.ItemWithQuantity {
	return client.ItemWithQuantity{
		Id:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c *ItemWithQuantityConvertor) ClientToEntities(items []client.ItemWithQuantity) (res []entity.ItemWithQuantity) {
	for _, i := range items {
		res = append(res, c.ClientToEntity(i))
	}
	return
}

func (c *ItemWithQuantityConvertor) ClientToEntity(i client.ItemWithQuantity) entity.ItemWithQuantity {
	return entity.ItemWithQuantity{
		ID:       i.Id,
		Quantity: i.Quantity,
	}
}
