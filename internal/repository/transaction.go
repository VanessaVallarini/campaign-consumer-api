package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Transaction struct {
	pool *pgxpool.Pool
}

func NewTransaction(pool *pgxpool.Pool) Transaction {
	return Transaction{
		pool: pool,
	}
}

func (t Transaction) Begin(ctx context.Context) (pgx.Tx, error) {
	return t.pool.Begin(ctx)
}
