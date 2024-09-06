package query

import (
	"context"

	"github.com/daffaromero/gorpc-template/protobuf/api"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderQuery interface {
	CreateOrder(ctx context.Context, tx pgx.Tx, order *api.Order) (*api.Order, error)
	GetOrder(ctx context.Context, id string) (*api.Order, error)
	ListOrders(ctx context.Context) ([]*api.Order, error)
	UpdateOrder(ctx context.Context, tx pgx.Tx, order *api.Order) (*api.Order, error)
	DeleteOrder(ctx context.Context, tx pgx.Tx, id string) error
}

type orderQuery struct {
	db *pgxpool.Pool
}

func NewOrderQuery(db *pgxpool.Pool) *orderQuery {
	return &orderQuery{db: db}
}

func (q *orderQuery) CreateOrder(ctx context.Context, tx pgx.Tx, order *api.Order) (*api.Order, error) {
	query := `INSERT INTO orders (id, user_id, items) VALUES ($1, $2, $3) RETURNING id, user_id, items`

	var createdOrder api.Order
	err := tx.QueryRow(ctx, query, order.Id, order.UserId, order.Items).Scan(&createdOrder.Id, &createdOrder.UserId, &createdOrder.Items)
	if err != nil {
		return nil, err
	}

	return &createdOrder, nil

}

func (q *orderQuery) GetOrder(ctx context.Context, id string) (*api.Order, error) {
	query := `SELECT id, user_id, items FROM orders WHERE id = $1`

	var order api.Order
	err := q.db.QueryRow(ctx, query, id).Scan(&order.Id, &order.UserId, &order.Items)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (q *orderQuery) ListOrders(ctx context.Context) ([]*api.Order, error) {
	query := `SELECT id, user_id, items FROM orders`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var orders []*api.Order
	for rows.Next() {
		var order api.Order
		err := rows.Scan(&order.Id, &order.UserId, &order.Items)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (q *orderQuery) UpdateOrder(ctx context.Context, tx pgx.Tx, order *api.Order) (*api.Order, error) {
	query := `UPDATE orders SET user_id = $1, items = $2 WHERE id = $3 RETURNING id, user_id, items`

	var updatedOrder api.Order
	err := tx.QueryRow(ctx, query, order.UserId, order.Items, order.Id).Scan(&updatedOrder.Id, &updatedOrder.UserId, &updatedOrder.Items)
	if err != nil {
		return nil, err
	}

	return &updatedOrder, nil
}

func (q *orderQuery) DeleteOrder(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM orders WHERE id = $1`

	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
