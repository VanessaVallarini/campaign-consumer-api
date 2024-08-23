package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type OwnerDao struct {
	pool *pgxpool.Pool
}

func NewOwnerDao(pool *pgxpool.Pool) OwnerDao {
	return OwnerDao{
		pool: pool,
	}
}

var upsertOwnerQuery = `
	INSERT INTO owner (id, email, status, created_by, updated_by, created_at, updated_at)
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
		status = EXCLUDED.status,
		email = EXCLUDED.email,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	WHERE
		owner.email <> EXCLUDED.email
		OR owner.status <> EXCLUDED.status;
`

func (od OwnerDao) Upsert(ctx context.Context, owner model.Owner) error {
	_, err := od.pool.Exec(
		ctx,
		upsertOwnerQuery,
		owner.Id,
		owner.Email,
		owner.Status,
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
