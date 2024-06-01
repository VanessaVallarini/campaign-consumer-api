package repository

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type RegionRepository struct {
	pool *pgxpool.Pool
}

func NewRegionRepository(pool *pgxpool.Pool) RegionRepository {
	return RegionRepository{
		pool: pool,
	}
}

var upsertRegionQuery = `
	INSERT INTO region (id, name, active, lat, long, cost, created_by, updated_by, created_at, updated_at)
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
		active = EXCLUDED.active,
		lat = EXCLUDED.lat,
		long = EXCLUDED.long,
		cost = EXCLUDED.cost,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	WHERE
		slug.name <> EXCLUDED.name
		OR region.active <> EXCLUDED.active
		OR region.lat <> EXCLUDED.lat
		OR region.long <> EXCLUDED.long
		OR region.cost <> EXCLUDED.cost;
`

func (s RegionRepository) Upsert(ctx context.Context, region model.Region) error {
	_, err := s.pool.Exec(
		ctx,
		upsertRegionQuery,
		region.Id,
		region.Name,
		region.Active,
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
