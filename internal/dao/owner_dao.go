package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	easyzap "github.com/lockp111/go-easyzap"
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

const allOwnerFields = `
	id, 
	email, 
	status, 
	created_by, 
	updated_by, 
	created_at, 
	updated_at
`

var upsertOwnerQuery = `
	INSERT INTO owner (` + allOwnerFields + `)
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
		easyzap.Error(err, "failed to create or update merchant in database")

		return errors.Wrap(err, "Failed to create or update owner in database")
	}

	return nil
}
