package query

import (
	"context"
	"errors"
	"fmt"

	"github.com/daffaromero/gorpc-template/protobuf/api"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemQuery interface {
	CreateItem(ctx context.Context, tx pgx.Tx, item *api.Item) (*api.Item, error)
	GetItem(ctx context.Context, id string) (*api.Item, error)
	ListItems(ctx context.Context) ([]*api.Item, error)
	UpdateItem(ctx context.Context, tx pgx.Tx, item *api.Item) (*api.Item, error)
	DeleteItem(ctx context.Context, tx pgx.Tx, id string) error
}

type itemQuery struct {
	db *pgxpool.Pool
}

func NewItemQuery(db *pgxpool.Pool) *itemQuery {
	return &itemQuery{db: db}
}

func (q *itemQuery) CreateItem(ctx context.Context, tx pgx.Tx, item *api.Item) (*api.Item, error) {
	if item == nil {
		return nil, errors.New("item cannot be nil")
	}

	query := `INSERT INTO items (id, name, description) VALUES ($1, $2, $3) RETURNING id, name, description`

	var createdItem api.Item
	err := tx.QueryRow(ctx, query, item.Id, item.Name, item.Description).Scan(&createdItem.Id, &createdItem.Name, &createdItem.Description)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			switch pgErr.Code {
			case "23505": // unique_violation
				return nil, fmt.Errorf("item with ID %s already exists: %w", item.Id, err)
			case "23502": // not_null_violation
				return nil, fmt.Errorf("missing required field: %w", err)
			}
		}
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	return &createdItem, nil
}

func (q *itemQuery) GetItem(ctx context.Context, id string) (*api.Item, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	query := `SELECT id, name, description FROM items WHERE id = $1`

	row := q.db.QueryRow(ctx, query, id)

	var item api.Item
	err := row.Scan(&item.Id, &item.Name, &item.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("item with ID %s not found", id)
		}
		return nil, err
	}

	return &item, nil
}

func (q *itemQuery) ListItems(ctx context.Context) ([]*api.Item, error) {
	query := `SELECT id, name, description FROM items`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var items []*api.Item

	for rows.Next() {
		var item api.Item
		err := rows.Scan(&item.Id, &item.Name, &item.Description)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	return items, nil
}

func (q *itemQuery) UpdateItem(ctx context.Context, tx pgx.Tx, item *api.Item) (*api.Item, error) {
	query := `UPDATE items SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description`

	var updatedItem api.Item
	err := tx.QueryRow(ctx, query, item.Name, item.Description, item.Id).Scan(&updatedItem.Id, &updatedItem.Name, &updatedItem.Description)
	if err != nil {
		return nil, err
	}

	return &updatedItem, nil
}

func (q *itemQuery) DeleteItem(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM items WHERE id = $1`

	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
