package repository

import (
	"context"

	"github.com/daffaromero/gorpc-template/protobuf/api"
	"github.com/daffaromero/gorpc-template/repository/query"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *api.Order) (*api.Order, error)
	GetOrder(ctx context.Context, id string) (*api.Order, error)
	ListOrders(ctx context.Context) ([]*api.Order, error)
	UpdateOrder(ctx context.Context, order *api.Order) (*api.Order, error)
	DeleteOrder(ctx context.Context, id string) error
}

type orderRepository struct {
	db         Store
	orderQuery query.OrderQuery
}

func NewOrderRepository(db Store, orderQuery query.OrderQuery) OrderRepository {
	return &orderRepository{db: db, orderQuery: orderQuery}
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *api.Order) (*api.Order, error) {

}
