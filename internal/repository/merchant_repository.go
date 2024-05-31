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

const allMerchantFields = `
	id, 
	owner_id, 
	slugs, 
	name, 
	active, 
	created_by, 
	updated_by, 
	created_at, 
	updated_at
`

var createMerchantQuery = `
	INSERT INTO merchant (` + allMerchantFields + `
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		$9
	);
`

func (m MerchantRepository) Create(ctx context.Context, merchant model.Merchant) error {
	_, err := m.pool.Exec(
		ctx,
		createMerchantQuery,
		merchant.Id,
		merchant.OwnerId,
		merchant.Slugs,
		merchant.Name,
		merchant.Active,
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

var updateMerchantQuery = `
	UPDATE merchant
	SET slugs=$2, "name"=$3, active=$4, updated_by=$5, updated_at=$6
	WHERE id=$1;
`

func (m MerchantRepository) Update(ctx context.Context, merchant model.Merchant) error {
	_, err := m.pool.Exec(
		ctx,
		updateMerchantQuery,
		merchant.Id,
		merchant.Slugs,
		merchant.Name,
		merchant.Active,
		merchant.UpdatedBy,
		merchant.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to update merchant in database")
	}

	return nil
}
