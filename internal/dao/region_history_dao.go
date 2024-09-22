package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type RegionHistoryDao struct {
	pool *pgxpool.Pool
}

func NewRegionHistoryDao(pool *pgxpool.Pool) RegionHistoryDao {
	return RegionHistoryDao{
		pool: pool,
	}
}

const allRegionHistoryFields = `
	id, 
	region_id, 
	status,
	description,
	created_by,
	created_at
`

var createRegionHistoryQuery = `
	INSERT INTO region_history (` + allRegionHistoryFields + `)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
	);
`

func (rh RegionHistoryDao) Create(ctx context.Context, history model.RegionHistory) error {
	_, err := rh.pool.Exec(
		ctx,
		createRegionHistoryQuery,
		history.Id,
		history.RegionId,
		history.Status,
		history.Description,
		history.CreatedBy,
		history.CreatedAt,
	)
	if err != nil {

		return errors.Wrap(err, "Failed to create region history in database")
	}

	return nil
}
