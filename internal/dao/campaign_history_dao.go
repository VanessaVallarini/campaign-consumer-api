package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type CampaignHistoryDao struct {
	pool *pgxpool.Pool
}

func NewCampaignHistoryDao(pool *pgxpool.Pool) CampaignHistoryDao {
	return CampaignHistoryDao{
		pool: pool,
	}
}

const allCampaignHistoryFields = `
	id, 
	campaign_id, 
	status,
	description,
	created_by,
	created_at
`

var createCampaignHistoryQuery = `
	INSERT INTO campaign_history (` + allCampaignHistoryFields + `)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
	);
`

func (ch CampaignHistoryDao) Create(ctx context.Context, history model.CampaignHistory) error {
	_, err := ch.pool.Exec(
		ctx,
		createCampaignHistoryQuery,
		history.Id,
		history.CampaignId,
		history.Status,
		history.Description,
		history.CreatedBy,
		history.CreatedAt,
	)
	if err != nil {

		return errors.Wrap(err, "Failed to create campaign history in database")
	}

	return nil
}
