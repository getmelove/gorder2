package adapters

import (
	"context"
	"time"

	_ "github.com/getmelove/gorder2/internal/common/config"
	domain "github.com/getmelove/gorder2/internal/order/domain/order"
	"github.com/getmelove/gorder2/internal/order/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	dbName   = viper.GetString("mongo.db-name")
	collName = viper.GetString("mongo.coll-name")
)

type OrderRepositoryMongo struct {
	db *mongo.Client
}

func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
	return &OrderRepositoryMongo{db: db}
}

func (r *OrderRepositoryMongo) collecting() *mongo.Collection {
	return r.db.Database(dbName).Collection(collName)
}

type orderModle struct {
	MongoID     primitive.ObjectID `bson:"_id"`
	ID          string             `bson:"id"`
	CustomerID  string             `bson:"customer_id"`
	Status      string             `bson:"status"`
	PaymentLink string             `bson:"payment_link"`
	Items       []*entity.Item     `bson:"items"`
}

func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.OrderAggregate) (created *domain.OrderAggregate, err error) {
	defer r.logWithTag("create", err, created)
	wm := r.marshalToModel(order)
	res, err := r.collecting().InsertOne(ctx, wm)
	if err != nil {
		return nil, err
	}
	created = order
	created.Id = res.InsertedID.(primitive.ObjectID).Hex()
	return
}

func (r *OrderRepositoryMongo) logWithTag(tag string, err error, result interface{}) {
	l := logrus.WithFields(logrus.Fields{
		"tag":         "order_repository_mongo",
		"create_time": time.Now().Unix(),
		"err":         err,
		"result":      result,
	})
	if err != nil {
		l.Infof("%s_failed", tag)
	} else {
		l.Infof("%s_success", tag)
	}
}

func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.OrderAggregate, err error) {
	defer r.logWithTag("get", nil, got)
	read := &orderModle{}
	mongoID, _ := primitive.ObjectIDFromHex(id)
	// condition
	cond := bson.M{"_id": mongoID}
	err = r.collecting().FindOne(ctx, cond).Decode(read)
	if err != nil {
		return
	}
	if read == nil {
		return nil, domain.NotFoundError{OrderID: id}
	}
	return r.unmarshal(read), nil
}

// Update 先查找对应的order， 然后apply updateFn, 再写入回去
func (r *OrderRepositoryMongo) Update(
	ctx context.Context,
	order *domain.OrderAggregate,
	updateFn func(context.Context, *domain.OrderAggregate) (*domain.OrderAggregate, error),
) (err error) {
	defer r.logWithTag("after_update", err, nil)
	if order == nil {
		panic("got nil order")
	}
	// 启用事务
	session, err := r.db.StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = session.CommitTransaction(ctx)
		} else {
			session.AbortTransaction(ctx)
		}
	}()
	// inside transaction
	oldOrder, err := r.Get(ctx, order.Id, order.CustomerID)
	if err != nil {
		return
	}
	updated, err := updateFn(ctx, oldOrder)
	if err != nil {
		return
	}
	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.Id)
	res, err := r.collecting().UpdateOne(
		ctx,
		bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
		bson.M{"$set": bson.M{
			"status":       updated.Status,
			"payment_link": updated.PaymentLink,
		}},
	)
	if err != nil {
		return
	}
	r.logWithTag("finish update", err, res)
	return
}

func (r *OrderRepositoryMongo) marshalToModel(order *domain.OrderAggregate) *orderModle {
	return &orderModle{
		MongoID:     primitive.NewObjectID(),
		ID:          order.Id,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
}

func (r *OrderRepositoryMongo) unmarshal(m *orderModle) *domain.OrderAggregate {
	return &domain.OrderAggregate{
		CustomerID:  m.MongoID.Hex(),
		Id:          m.ID,
		Items:       m.Items,
		PaymentLink: m.PaymentLink,
		Status:      m.Status,
	}
}
