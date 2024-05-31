package repository

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type OwnerRepository struct {
	pool *pgxpool.Pool
}

func NewOwnerRepository(pool *pgxpool.Pool) OwnerRepository {
	return OwnerRepository{
		pool: pool,
	}
}

var upsertOwnerQuery = `
	INSERT INTO owner (id, email, active, created_by, updated_by, created_at, updated_at)
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
		active = EXCLUDED.active,
		email = EXCLUDED.email,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	WHERE
		owner.email <> EXCLUDED.email
		OR owner.active <> EXCLUDED.active;
`

func (o OwnerRepository) Upsert(ctx context.Context, owner model.Owner) error {
	_, err := o.pool.Exec(
		ctx,
		upsertOwnerQuery,
		owner.Id,
		owner.Email,
		owner.Active,
		owner.CreatedBy,
		owner.UpdatedBy,
		owner.CreatedAt,
		owner.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create or update owner in database")
	}

	return nil
}
