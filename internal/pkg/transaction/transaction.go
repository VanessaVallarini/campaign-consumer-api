package transaction

import (
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

type Transaction struct {
	it pgx.Tx
}

func NewTransaction(it pgx.Tx) Transaction {
	return Transaction{it: it}
}
func (t Transaction) Commit(ctx context.Context) error {
	return t.it.Commit(ctx)
}

func (t Transaction) Rollback(ctx context.Context) error {
	return t.it.Rollback(ctx)
}

func (t Transaction) Exec(ctx context.Context, query string, arguments ...interface{}) error {
	_, err := t.it.Exec(ctx, query, arguments...)
	return err
}
