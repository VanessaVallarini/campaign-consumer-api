package transaction

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	easyzap "github.com/lockp111/go-easyzap"
	"golang.org/x/net/context"
)

// unit of work

type DBConn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type TransactionManager struct {
	rw DBConn
}

func NewTransactionManager(rw *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{rw: rw}
}

func (tm TransactionManager) Execute(ctx context.Context, fn func(ctx context.Context, tx Transaction) error) error {
	tx, err := tm.rw.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				easyzap.Error(errRollback, "[INCONSISTENCY] Failed to rollback operation")
			}
		}
	}()

	err = fn(ctx, NewTransaction(tx))
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	return err
}
