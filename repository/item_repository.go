package repository

import (
	"context"

	"github.com/daffaromero/gorpc-template/protobuf/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemRepository interface {
	CreateItem(ctx context.Context, item *api.Item) (*api.Item, error)
	GetItem(ctx context.Context, id string) (*api.Item, error)
	ListItems(ctx context.Context) ([]*api.Item, error)
	UpdateItem(ctx context.Context, item *api.Item) (*api.Item, error)
	DeleteItem(ctx context.Context, id string) error
}

type itemRepository struct {
	db Store
}

func NewItemRepository(db *pgxpool.Pool) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) CreateItem(ctx context.Context, item *api.Item) (*api.Item, error) {
	query := `INSERT INTO items (id, name, description) VALUES ($1, $2, $3) RETURNING id, name, description`

	var createdItem api.Item
	err := r.db.WithTx(ctx, query, item.Id, item.Name, item.Description).Scan(&createdItem.Id, &createdItem.Name, &createdItem.Description)
	if err != nil {
		return nil, err
	}

	return &createdItem, nil
}
