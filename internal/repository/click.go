package repository

import "github.com/jackc/pgx/v5/pgxpool"

type ClickRepository struct {
	pool *pgxpool.Pool
}

func NewClickRepository(pool *pgxpool.Pool) ClickRepository {
	return ClickRepository{
		pool: pool,
	}
}
