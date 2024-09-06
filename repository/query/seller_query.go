package query

import (
	"context"

	"github.com/daffaromero/gorpc-template/protobuf/api"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SellerQuery interface {
	CreateSeller(ctx context.Context, tx pgx.Tx, seller *api.Seller) (*api.Seller, error)
	GetSeller(ctx context.Context, id string) (*api.Seller, error)
	ListSellers(ctx context.Context) ([]*api.Seller, error)
	UpdateSeller(ctx context.Context, tx pgx.Tx, seller *api.Seller) (*api.Seller, error)
	DeleteSeller(ctx context.Context, tx pgx.Tx, id string) error
}

type sellerQuery struct {
	db *pgxpool.Pool
}

func NewSellerQuery(db *pgxpool.Pool) *sellerQuery {
	return &sellerQuery{db: db}
}

func (q *sellerQuery) CreateSeller(ctx context.Context, tx pgx.Tx, seller *api.Seller) (*api.Seller, error) {
	query := `INSERT INTO sellers (id, name) VALUES ($1, $2) RETURNING id, name`

	var createdSeller api.Seller
	err := tx.QueryRow(ctx, query, seller.Id, seller.Name).Scan(&createdSeller.Id, &createdSeller.Name)
	if err != nil {
		return nil, err
	}

	return &createdSeller, nil
}

func (q *sellerQuery) GetSeller(ctx context.Context, id string) (*api.Seller, error) {
	query := `SELECT id, name FROM sellers WHERE id = $1`

	row := q.db.QueryRow(ctx, query, id)

	var seller api.Seller
	err := row.Scan(&seller.Id, &seller.Name)
	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func (q *sellerQuery) ListSellers(ctx context.Context) ([]*api.Seller, error) {
	query := `SELECT id, name FROM sellers`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var sellers []*api.Seller
	for rows.Next() {
		var seller api.Seller
		err := rows.Scan(&seller.Id, &seller.Name)
		if err != nil {
			return nil, err
		}
		sellers = append(sellers, &seller)
	}

	return sellers, nil
}

func (q *sellerQuery) UpdateSeller(ctx context.Context, tx pgx.Tx, seller *api.Seller) (*api.Seller, error) {
	query := `UPDATE sellers SET name = $1 WHERE id = $2 RETURNING id, name`

	var updatedSeller api.Seller
	err := tx.QueryRow(ctx, query, seller.Name, seller.Id).Scan(&updatedSeller.Id, &updatedSeller.Name)
	if err != nil {
		return nil, err
	}

	return &updatedSeller, nil
}

func (q *sellerQuery) DeleteSeller(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM sellers WHERE id = $1`

	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
