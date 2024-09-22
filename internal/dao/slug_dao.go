package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type SlugDao struct {
	pool *pgxpool.Pool
}

func NewSlugDao(pool *pgxpool.Pool) SlugDao {

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

			return model.Slug{}, model.ErrNotFound
		}

		return model.Slug{}, errors.Wrap(err, "Failed to fetch slug in database")
	}

	return slug, nil
}

var createSlugQuery = `
	INSERT INTO slug (` + allSlugFields + `)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8
		);
`

func (sd SlugDao) Create(ctx context.Context, tx transaction.Transaction, slug model.Slug) error {
	err := tx.Exec(
		ctx,
		createSlugQuery,
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

		return errors.Wrap(err, "Failed to create slug in database")
	}

	return nil
}

var updateSlugQuery = `
	UPDATE
		slug
	SET 
		status = $1,
		cost = $2,
		updated_by = $3,
		updated_at = $4
	WHERE
		id = $5;
`

func (sd SlugDao) Update(ctx context.Context, tx transaction.Transaction, slug model.Slug) error {
	err := tx.Exec(
		ctx,
		updateSlugQuery,
		slug.Status,
		slug.Cost,
		slug.UpdatedBy,
		slug.UpdatedAt,
		slug.Id,
	)
	if err != nil {

		return errors.Wrap(err, "Failed to update slug in database")
	}

	return nil
}
