package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type SlugDao struct {
	pool *pgxpool.Pool
}

func NewSlugRepository(pool *pgxpool.Pool) SlugDao {

	return SlugDao{
		pool: pool,
	}
}

const allSlugFields = `
	id, 
	name, 
	status, 
	cost, 
	created_by, 
	updated_by, 
	created_at, 
	updated_at
`

var upsertSlugQuery = `
	INSERT INTO slug (id, name, status, cost, created_by, updated_by, created_at, updated_at)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)
	ON CONFLICT (id) DO UPDATE
	SET
		name = EXCLUDED.name,
		status = EXCLUDED.status,
		cost = EXCLUDED.cost,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	WHERE
		slug.name <> EXCLUDED.name
		OR slug.status <> EXCLUDED.status
		OR slug.cost <> EXCLUDED.cost;
`

func (sd SlugDao) Upsert(ctx context.Context, slug model.Slug) error {
	_, err := sd.pool.Exec(
		ctx,
		upsertSlugQuery,
		slug.Id,
		slug.Name,
		slug.Status,
		slug.Cost,
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

func (sd SlugDao) Fetch(ctx context.Context, id uuid.UUID) (model.Slug, error) {
	var slug model.Slug

	query := `SELECT ` + allSlugFields + ` from slug WHERE id = $1`

	row := sd.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&slug.Id, &slug.Name, &slug.Status, &slug.Cost, &slug.CreatedBy,
		&slug.UpdatedBy, &slug.CreatedAt, &slug.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Slug{}, errors.Wrap(err, "Slug not found")
		}

		return model.Slug{}, errors.Wrap(err, "Failed to fetch slug in database")
	}

	return slug, nil
}
