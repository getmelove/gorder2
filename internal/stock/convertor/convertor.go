package convertor

import (
	"github.com/getmelove/gorder2/internal/common/genproto/orderpb"
	"github.com/getmelove/gorder2/internal/stock/entity"
)

type ItemConvertor struct{}

func NewItemConvertor() *ItemConvertor {
	return &ItemConvertor{}
}

func (c *ItemConvertor) ToProto(i *entity.Item) *orderpb.Item {
	return &orderpb.Item{
		ID:       i.ID,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceID,
	}
}

func (c *ItemConvertor) ToProtos(items []*entity.Item) []*orderpb.Item {
	var res []*orderpb.Item
	for _, item := range items {
		res = append(res, c.ToProto(item))
	}
	return res
}

func (c *ItemConvertor) ToEntity(i *orderpb.Item) *entity.Item {
	return &entity.Item{
		ID:       i.ID,
		Name:     i.Name,
		Quantity: i.Quantity,
		PriceID:  i.PriceID,
	}
}

func (c *ItemConvertor) ProtoToEntity(i *orderpb.ItemWithQuantity) *entity.ItemWithQuantity {
	return &entity.ItemWithQuantity{
		ID:       i.ID,
		Quantity: i.Quantity,
	}
}

func (c *ItemConvertor) ProtoToEntities(items []*orderpb.ItemWithQuantity) []*entity.ItemWithQuantity {
	var res []*entity.ItemWithQuantity
	for _, item := range items {
		res = append(res, c.ProtoToEntity(item))
	}
	return res
}
