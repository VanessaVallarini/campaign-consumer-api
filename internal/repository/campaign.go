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

const allCampaignFields = `
	id, 
	merchant_id, 
	active, 
	lat,
	long,
	created_by, 
	updated_by, 
	created_at, 
	updated_at
`

var createCampaignQuery = `
	INSERT INTO campaign (` + allCampaignFields + `
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

func (c CampaignRepository) Create(ctx context.Context, campaign model.Campaign) error {
	_, err := c.pool.Exec(
		ctx,
		createCampaignQuery,
		campaign.Id,
		campaign.MerchantId,
		campaign.Active,
		campaign.Lat,
		campaign.Long,
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

var updateCampaignQuery = `
	UPDATE campaign
	SET active=$2, lat=$3, long=$4, updated_by=$5, updated_at=$6
	WHERE id=$1;
`

func (c CampaignRepository) Update(ctx context.Context, campaign model.Campaign) error {
	_, err := c.pool.Exec(
		ctx,
		updateCampaignQuery,
		campaign.Id,
		campaign.Active,
		campaign.Lat,
		campaign.Long,
		campaign.UpdatedBy,
		campaign.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to update campaign in database")
	}

	return nil
}
