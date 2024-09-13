package repository

import (
	"context"
	"fmt"

	"github.com/daffaromero/gorpc-template/protobuf/api"
	"github.com/daffaromero/gorpc-template/repository/query"

	"github.com/jackc/pgx/v5"
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

func NewItemRepository(db Store, itemQuery query.ItemQuery) ItemRepository {
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
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return item, nil
}

func (r *itemRepository) ListItems(ctx context.Context) ([]*api.Item, error) {
	var items []*api.Item

	err := r.db.WithoutTx(ctx, func(ctx context.Context) error {
		var err error
		items, err = r.itemQuery.ListItems(ctx)
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}
	return items, nil
}

func (r *itemRepository) UpdateItem(ctx context.Context, item *api.Item) (*api.Item, error) {
	var updatedItem *api.Item

	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		var err error
		updatedItem, err = r.itemQuery.UpdateItem(ctx, tx, item)
		if err != nil {
			return fmt.Errorf("failed to update item: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("transaction failed: %w", err)
	}

	return updatedItem, nil
}

func (r *itemRepository) DeleteItem(ctx context.Context, id string) error {
	err := r.db.WithTx(ctx, func(tx pgx.Tx) error {
		err := r.itemQuery.DeleteItem(ctx, tx, id)
		if err != nil {
			return fmt.Errorf("failed to delete item: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	return nil
}
