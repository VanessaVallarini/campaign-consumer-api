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

type SpentDao struct {
	pool *pgxpool.Pool
}

func NewSpentDao(pool *pgxpool.Pool) SpentDao {
	return SpentDao{
		pool: pool,
	}
}

const allSpentFields = `
	id,
	campaign_id,
	merchant_id,
	bucket,
	total_spent,
	total_clicks,
	total_impressions
`

var upsertSpentQuery = `
	INSERT INTO spent 
	(id, campaign_id, merchant_id,  bucket, total_spent, total_clicks, total_impressions)
	VALUES ($1,$2,$3,$4,$5,$6,$7)
	ON CONFLICT (merchant_id, bucket) 
	DO update set 
	total_spent = EXCLUDED.total_spent,
	total_clicks = EXCLUDED.total_clicks,
	total_impressions = EXCLUDED.total_impressions;
`

func (s SpentDao) Upsert(ctx context.Context, tx transaction.Transaction, spent model.Spent) error {
	err := tx.Exec(
		ctx,
		upsertSpentQuery,
		spent.Id,
		spent.CampaignId,
		spent.MerchantId,
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

func (sd SpentDao) FetchByMerchantIdAndBucket(ctx context.Context, id uuid.UUID, bucket string) (model.Spent, error) {
	var spent model.Spent

	query := `SELECT` + allSpentFields + ` from spent WHERE merchant_id = $1 AND bucket = $2`

	row := sd.pool.QueryRow(ctx, query, id, bucket)
	err := row.Scan(
		&spent.Id, &spent.CampaignId, &spent.MerchantId, &spent.Bucket,
		&spent.TotalSpent, &spent.TotalClicks, &spent.TotalImpressions,
	)

	if err != nil {
		if err == pgx.ErrNoRows {

			return model.Spent{}, model.ErrNotFound
		}

		return model.Spent{}, errors.Wrap(err, "Failed to fetch spent by campaign id in database")
	}

	return spent, nil
}
