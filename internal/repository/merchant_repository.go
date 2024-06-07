package repository

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type MerchantRepository struct {
	pool *pgxpool.Pool
}

func NewMerchantRepository(pool *pgxpool.Pool) MerchantRepository {
	return MerchantRepository{
		pool: pool,
	}
}

var upsertMerchantQuery = `
	INSERT INTO merchant (id, owner_id, region_id, slugs, name, status, created_by, updated_by, created_at, updated_at)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9,
		$10
	)
	ON CONFLICT (id) DO UPDATE
	SET
		slugs = EXCLUDED.slugs,
		name = EXCLUDED.name,
		status = EXCLUDED.status,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	WHERE
		OR merchant.slugs <> EXCLUDED.slugs
		OR merchant.name <> EXCLUDED.name
		OR merchant.status <> EXCLUDED.status;
`

func (m MerchantRepository) Upsert(ctx context.Context, merchant model.Merchant) error {
	_, err := m.pool.Exec(
		ctx,
		upsertMerchantQuery,
		merchant.Id,
		merchant.OwnerId,
		merchant.RegionId,
		merchant.Slugs,
		merchant.Name,
		merchant.Status,
		merchant.CreatedBy,
		merchant.UpdatedBy,
		merchant.CreatedAt,
		merchant.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create merchant in database")
	}

	return nil
}
