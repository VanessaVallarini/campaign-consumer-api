package dao

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type SlugHistoryDao struct {
	pool *pgxpool.Pool
}

func NewSlugHistoryDao(pool *pgxpool.Pool) SlugHistoryDao {
	return SlugHistoryDao{
		pool: pool,
	}
}

const allSlugHistoryFields = `
	id, 
	slug_id, 
	status,
	description,
	created_by,
	created_at
`

var createSlugHistoryQuery = `
	INSERT INTO slug_history (` + allSlugHistoryFields + `)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		);
`

func (sh SlugHistoryDao) Create(ctx context.Context, tx transaction.Transaction, history model.SlugHistory) error {
	err := tx.Exec(
		ctx,
		createSlugHistoryQuery,
		history.Id,
		history.SlugId,
		history.Status,
		history.Description,
		history.CreatedBy,
		history.CreatedAt,
	)
	if err != nil {

		return errors.Wrap(err, "Failed to create slug history in database")
	}

	return nil
}
