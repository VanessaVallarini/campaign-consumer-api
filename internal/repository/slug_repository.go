package repository

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type SlugRepository struct {
	pool *pgxpool.Pool
}

func NewSlugRepository(pool *pgxpool.Pool) SlugRepository {
	return SlugRepository{
		pool: pool,
	}
}

var upsertSlugQuery = `
	INSERT INTO slug (id, name, active, created_by, updated_by, created_at, updated_at)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7
	)
	ON CONFLICT (id) DO UPDATE
	SET
		email = EXCLUDED.name,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	WHERE
		slug.email <> EXCLUDED.name
		OR slug.active <> EXCLUDED.active;
`

func (s SlugRepository) Upsert(ctx context.Context, slug model.Slug) error {
	_, err := s.pool.Exec(
		ctx,
		upsertSlugQuery,
		slug.Id,
		slug.Name,
		slug.Active,
		slug.CreatedBy,
		slug.UpdatedBy,
		slug.CreatedAt,
		slug.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create or update slug in database")
	}

	return nil
}
