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

type RegionDao struct {
	pool *pgxpool.Pool
}

func NewRegionDao(pool *pgxpool.Pool) RegionDao {
	return RegionDao{
		pool: pool,
	}
}

const allRegionFields = `
	id,
	name,
	status,
	lat,
	long,
	cost,
	created_by,
	updated_by,
	created_at,
	updated_at
`

func (rd RegionDao) Fetch(ctx context.Context, id uuid.UUID) (model.Region, error) {
	var region model.Region

	query := `SELECT ` + allRegionFields + ` from region WHERE id = $1`

	row := rd.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&region.Id, &region.Name, &region.Status, &region.Lat,
		&region.Long, &region.Cost, &region.CreatedBy,
		&region.UpdatedBy, &region.CreatedAt, &region.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Region{}, model.ErrNotFound
		}

		return model.Region{}, errors.Wrap(err, "Failed to fetch region in database")
	}

	return region, nil
}

var createRegionQuery = `
	INSERT INTO region (` + allRegionFields + `)
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
	);
`

func (rd RegionDao) Create(ctx context.Context, tx transaction.Transaction, region model.Region) error {
	err := tx.Exec(
		ctx,
		createRegionQuery,
		region.Id,
		region.Name,
		region.Status,
		region.Lat,
		region.Long,
		region.Cost,
		region.CreatedBy,
		region.UpdatedBy,
		region.CreatedAt,
		region.UpdatedAt,
	)
	if err != nil {

		return errors.Wrap(err, "Failed to create region in database")
	}

	return nil
}

var updateRegionQuery = `
	UPDATE
		region
	SET
		status = $1,
		cost = $2,
		updated_by = $3,
		updated_at = $4
	WHERE
		id = $5;
`

func (rd RegionDao) Update(ctx context.Context, tx transaction.Transaction, region model.Region) error {
	err := tx.Exec(
		ctx,
		updateRegionQuery,
		region.Status,
		region.Cost,
		region.UpdatedBy,
		region.UpdatedAt,
		region.Id,
	)
	if err != nil {

		return errors.Wrap(err, "Failed to update region in database")
	}

	return nil
}
