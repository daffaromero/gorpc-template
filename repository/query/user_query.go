package query

import (
	"context"

	"github.com/daffaromero/gorpc-template/protobuf/api"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserQuery interface {
	CreateUser(ctx context.Context, tx pgx.Tx, user *api.User) (*api.User, error)
	GetUser(ctx context.Context, id string) (*api.User, error)
	ListUsers(ctx context.Context) ([]*api.User, error)
	UpdateUser(ctx context.Context, tx pgx.Tx, user *api.User) (*api.User, error)
	DeleteUser(ctx context.Context, tx pgx.Tx, id string) error
}

type userQuery struct {
	db *pgxpool.Pool
}

func NewUserQuery(db *pgxpool.Pool) *userQuery {
	return &userQuery{db: db}
}

func (q *userQuery) CreateUser(ctx context.Context, tx pgx.Tx, user *api.User) (*api.User, error) {
	query := `INSERT INTO users (id, name, password) VALUES ($1, $2, $3) RETURNING id, name, password`

	var createdUser api.User
	err := tx.QueryRow(ctx, query, user.Id, user.Name, user.Password).Scan(&createdUser.Id, &createdUser.Name, &createdUser.Password)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (q *userQuery) GetUser(ctx context.Context, id string) (*api.User, error) {
	query := `SELECT id, name, password FROM users WHERE id = $1`

	row := q.db.QueryRow(ctx, query, id)

	var user api.User
	err := row.Scan(&user.Id, &user.Name, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (q *userQuery) ListUsers(ctx context.Context) ([]*api.User, error) {
	query := `SELECT id, name, password FROM users`

	rows, err := q.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*api.User
	for rows.Next() {
		var user api.User
		err := rows.Scan(&user.Id, &user.Name, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (q *userQuery) UpdateUser(ctx context.Context, tx pgx.Tx, user *api.User) (*api.User, error) {
	query := `UPDATE users SET name = $1, password = $2 WHERE id = $3 RETURNING id, name`

	var updatedUser api.User
	err := tx.QueryRow(ctx, query, user.Name, user.Password, user.Id).Scan(&updatedUser.Id, &updatedUser.Name)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}

func (q *userQuery) DeleteUser(ctx context.Context, tx pgx.Tx, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
