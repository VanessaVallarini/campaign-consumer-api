package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

const allMerchantFields = `
	id, 
	owner_id, 
	region_id, 
	slugs, 
	name,
	status,
	created_by,
	updated_by, 
	created_at, 
	updated_at
`

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
		merchant.slugs <> EXCLUDED.slugs
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

func (m MerchantRepository) Fetch(ctx context.Context, id uuid.UUID) (model.Merchant, error) {
	var merchant model.Merchant

	query := `SELECT ` + allMerchantFields + ` from merchant WHERE id = $1`

	row := m.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&merchant.Id, &merchant.OwnerId, &merchant.RegionId, &merchant.Slugs,
		&merchant.Name, &merchant.Status, &merchant.CreatedBy,
		&merchant.UpdatedBy, &merchant.CreatedAt, &merchant.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Merchant{}, errors.Wrap(err, "Merchant not found")
		}
		return model.Merchant{}, errors.Wrap(err, "Failed to fetch merchant in database")
	}

	return merchant, nil
}
