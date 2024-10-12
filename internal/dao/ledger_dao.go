package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type LedgerDao struct {
	pool *pgxpool.Pool
}

func NewLedgerDao(pool *pgxpool.Pool) LedgerDao {
	return LedgerDao{
		pool: pool,
	}
}

const allLedgerFields = `
	id,
	spent_id,
	campaign_id,
	merchant_id,
	slug_name,
	region_name,
	user_id,
	event_type,
	cost,
	ip,
	lat,
	long, 
	created_at, 
	event_time
`

var createLedgerQuery = `
	INSERT INTO ledger (` + allLedgerFields + `)
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
		$10,
		$11,
		$12,
		$13,
		$14
	);
`

func (ld LedgerDao) Create(ctx context.Context, tx transaction.Transaction, ledger model.Ledger) error {
	err := tx.Exec(
		ctx,
		createLedgerQuery,
		ledger.Id,
		ledger.SpentId,
		ledger.CampaignId,
		ledger.MerchantId,
		ledger.SlugName,
		ledger.RegionName,
		ledger.UserId,
		ledger.EventType,
		ledger.Cost,
		ledger.Ip,
		ledger.Lat,
		ledger.Long,
		ledger.CreatedAt,
		ledger.EventTime,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create ledger in database")
	}

	return nil
}
