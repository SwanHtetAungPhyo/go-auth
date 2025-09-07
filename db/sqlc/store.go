package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	db *pgxpool.Pool
}

func NewStore(database *pgxpool.Pool) *Store {
	return &Store{
		Queries: New(database),
		db:      database,
	}
}

func (s *Store) Close() {
	s.db.Close()
}

type TxFunc func(*Queries) error

func (s *Store) WithTx(ctx context.Context, fn TxFunc) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("rolling back tx: %w", err)
		}
	}(tx, ctx)

	if err := fn(s.Queries.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
