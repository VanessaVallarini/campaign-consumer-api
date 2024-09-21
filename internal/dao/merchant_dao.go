package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	easyzap "github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
)

type MerchantDao struct {
	pool *pgxpool.Pool
}

func NewMerchantDao(pool *pgxpool.Pool) MerchantDao {

	return MerchantDao{
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
	INSERT INTO merchant (` + allMerchantFields + `)
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

func (md MerchantDao) Upsert(ctx context.Context, merchant model.Merchant) error {
	_, err := md.pool.Exec(
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
		easyzap.Error(err, "failed to create or update merchant in database")

		return errors.Wrap(err, "Failed to create or update merchant in database")
	}

	return nil
}

func (md MerchantDao) Fetch(ctx context.Context, id uuid.UUID) (model.Merchant, error) {
	var merchant model.Merchant

	query := `SELECT ` + allMerchantFields + ` from merchant WHERE id = $1`

	row := md.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&merchant.Id, &merchant.OwnerId, &merchant.RegionId, &merchant.Slugs,
		&merchant.Name, &merchant.Status, &merchant.CreatedBy,
		&merchant.UpdatedBy, &merchant.CreatedAt, &merchant.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Merchant{}, model.ErrNotFound
		}
		easyzap.Error(err, "failed to fetch merchant in database")

		return model.Merchant{}, errors.Wrap(err, "Failed to fetch merchant in database")
	}

	return merchant, nil
}
