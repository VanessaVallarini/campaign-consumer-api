package repository

import "github.com/jackc/pgx/v5/pgxpool"

type ImpressionRepository struct {
	pool *pgxpool.Pool
}

func NewImpressionRepository(pool *pgxpool.Pool) ImpressionRepository {
	return ImpressionRepository{
		pool: pool,
	}
}
