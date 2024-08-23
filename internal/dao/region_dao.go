package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
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

var upsertRegionQuery = `
	INSERT INTO region (id, name, status, lat, long, cost, created_by, updated_by, created_at, updated_at)
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
		name = EXCLUDED.name,
		status = EXCLUDED.status,
		lat = EXCLUDED.lat,
		long = EXCLUDED.long,
		cost = EXCLUDED.cost,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	WHERE
		region.name <> EXCLUDED.name
		OR region.status <> EXCLUDED.status
		OR region.lat <> EXCLUDED.lat
		OR region.long <> EXCLUDED.long
		OR region.cost <> EXCLUDED.cost;
`

func (rd RegionDao) Upsert(ctx context.Context, region model.Region) error {
	_, err := rd.pool.Exec(
		ctx,
		upsertRegionQuery,
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
		return errors.Wrap(err, "Failed to create or update region in database")
	}

	return nil
}

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
			return model.Region{}, errors.Wrap(err, "Region not found")
		}
		return model.Region{}, errors.Wrap(err, "Failed to fetch region in database")
	}

	return region, nil
}
