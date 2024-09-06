package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/daffaromero/gorpc-template/config"
	"github.com/daffaromero/gorpc-template/helper/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error
	WithoutTx(ctx context.Context, fn func(pool *pgxpool.Pool) error) error
}

type store struct {
	db     *pgxpool.Pool
	config config.DBConfig
	logger *logger.Log
}

func NewStore(db *pgxpool.Pool, config config.DBConfig) Store {
	return &store{db: db, config: config}
}

func (s *store) WithTx(ctx context.Context, fn func(tx pgx.Tx) error) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.TimeOutDuration)*time.Second)
	defer cancel()

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				s.logger.Error("Rollback error", "error", rollbackErr, "original_error", err)

				err = fmt.Errorf("rollback error: %v (original error: %w)", rollbackErr, err)
			}
		}
	}()

	if err = fn(tx); err != nil {
		return fmt.Errorf("transaction function failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *store) WithoutTx(ctx context.Context, fn func(pool *pgxpool.Pool) error) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.TimeOutDuration)*time.Second)
	defer cancel()

	return fn(s.db)
}
