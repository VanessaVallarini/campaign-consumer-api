package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type CampaignDao struct {
	pool *pgxpool.Pool
}

func NewCampaignDao(pool *pgxpool.Pool) CampaignDao {
	return CampaignDao{
		pool: pool,
	}
}

const allCampaignFields = `
	id, 
	merchant_id, 
	status,
	budget,
	created_by,
	updated_by, 
	created_at, 
	updated_at
`

func (c CampaignDao) Fetch(ctx context.Context, id uuid.UUID) (model.Campaign, error) {
	var campaign model.Campaign

	query := `SELECT ` + allCampaignFields + ` from campaign WHERE id = $1`

	row := c.pool.QueryRow(ctx, query, id)
	err := row.Scan(
		&campaign.Id, &campaign.MerchantId, &campaign.Status,
		&campaign.Budget, &campaign.CreatedBy, &campaign.UpdatedBy,
		&campaign.CreatedAt, &campaign.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Campaign{}, model.ErrNotFound
		}

		return model.Campaign{}, errors.Wrap(err, "Failed to fetch campaign in database")
	}

	return campaign, nil
}

var createCampaignQuery = `
	INSERT INTO campaign (` + allCampaignFields + `)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	);
`

func (c CampaignDao) Create(ctx context.Context, campaign model.Campaign) error {
	_, err := c.pool.Exec(
		ctx,
		createCampaignQuery,
		campaign.Id,
		campaign.MerchantId,
		campaign.Status,
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

var udateCampaignQuery = `
	UPDATE 
		campaign 
	SET
		status = $1,
		budget = $2,
		updated_by = $3,
		updated_at = $4
	WHERE
		id = $5;
`

func (c CampaignDao) Update(ctx context.Context, campaign model.Campaign) error {
	_, err := c.pool.Exec(
		ctx,
		udateCampaignQuery,
		campaign.Status,
		campaign.Budget,
		campaign.UpdatedBy,
		campaign.UpdatedAt,
		campaign.Id,
	)
	if err != nil {

		return errors.Wrap(err, "Failed to update campaign in database")
	}

	return nil
}
