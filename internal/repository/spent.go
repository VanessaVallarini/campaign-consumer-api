package repository

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type SpentRepository struct {
	pool *pgxpool.Pool
}

func NewSpentRepository(pool *pgxpool.Pool) SpentRepository {
	return SpentRepository{
		pool: pool,
	}
}

var upsertSpentQuery = `
	INSERT INTO spent (id, campaign_id, bucket, total_spent, total_clicks, total_impressions)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6
	)
	ON CONFLICT (id) DO UPDATE
	SET
		total_spent = EXCLUDED.total_spent,
		total_clicks = EXCLUDED.total_clicks,
		total_impressions = EXCLUDED.total_impressions;
`

func (s SpentRepository) Upsert(ctx context.Context, tx pgx.Tx, spent model.Spent) error {
	_, err := tx.Exec(
		ctx,
		upsertSpentQuery,
		spent.Id,
		spent.CampaignId,
		spent.Bucket,
		spent.TotalSpent,
		spent.TotalClicks,
		spent.TotalImpressions,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create or update spent in database")
	}

	return nil
}
