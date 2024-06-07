package repository

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type CampaignRepository struct {
	pool *pgxpool.Pool
}

func NewCampaignRepository(pool *pgxpool.Pool) CampaignRepository {
	return CampaignRepository{
		pool: pool,
	}
}

var upsertCampaignQuery = `
	INSERT INTO campaign (id, merchant_id, status, lat, long, budget, created_by, updated_by, created_at, updated_at)
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
		status = EXCLUDED.status,
		lat = EXCLUDED.lat,
		long = EXCLUDED.long,
		budget = EXCLUDED.budget,
		updated_by = EXCLUDED.updated_by,
		updated_at = EXCLUDED.updated_at
	WHERE
		campaign.status <> EXCLUDED.status
		OR campaign.lat <> EXCLUDED.lat
		OR campaign.long <> EXCLUDED.long
		OR campaign.budget <> EXCLUDED.budget;
`

func (c CampaignRepository) Upsert(ctx context.Context, campaign model.Campaign) error {
	_, err := c.pool.Exec(
		ctx,
		upsertCampaignQuery,
		campaign.Id,
		campaign.MerchantId,
		campaign.Status,
		campaign.Lat,
		campaign.Long,
		campaign.Budget,
		campaign.CreatedBy,
		campaign.UpdatedBy,
		campaign.CreatedAt,
		campaign.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create campaign in database")
	}

	return nil
}
