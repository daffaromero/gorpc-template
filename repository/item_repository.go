package repository

import (
	"context"
	"fmt"

	"github.com/daffaromero/gorpc-template/protobuf/api"
	"github.com/daffaromero/gorpc-template/repository/query"

	"github.com/jackc/pgx/v5"
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
	db        Store
	itemQuery query.ItemQuery
}

func NewItemRepository(db *pgxpool.Pool, itemQuery query.ItemQuery) ItemRepository {
	return &itemRepository{db: db, itemQuery: itemQuery}
}

func (r *itemRepository) CreateItem(ctx context.Context, item *api.Item) (*api.Item, error) {
	var createdItem *api.Item

	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		var err error
		createdItem, err = r.itemQuery.CreateItem(ctx, tx, item)
		if err != nil {
			return fmt.Errorf("failed to create item: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return createdItem, nil
}

func (r *itemRepository) GetItem(ctx context.Context, id string) (*api.Item, error) {
	var item *api.Item

	err := r.db.WithoutTx(ctx, func(ctx context.Context) error {
		var err error
		item, err = r.itemQuery.GetItem(ctx, id)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return item, nil
}
