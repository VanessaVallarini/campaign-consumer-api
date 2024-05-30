package repository

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type LedgerRepository struct {
	pool *pgxpool.Pool
}

func NewLedgerRepository(pool *pgxpool.Pool) LedgerRepository {
	return LedgerRepository{
		pool: pool,
	}
}

var createLedgerQuery = `
	INSERT INTO ledger (id, campaign_id, spent_id, user_id, event_type, cost, lat, long, created_at)
	VALUES (
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

func (s LedgerRepository) Create(ctx context.Context, tx pgx.Tx, ledger model.Ledger) error {
	_, err := tx.Exec(
		ctx,
		createLedgerQuery,
		ledger.Id,
		ledger.CampaignId,
		ledger.SpentId,
		ledger.UserId,
		ledger.EventType,
		ledger.Cost,
		ledger.Lat,
		ledger.Long,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create ledger in database")
	}

	return nil
}
