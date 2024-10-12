package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
)

type LedgerDao interface {
	Create(context.Context, transaction.Transaction, model.Ledger) error
}

type LedgerService struct {
	ledgerDao LedgerDao
}

func NewLedgerService(ledgerDao LedgerDao) LedgerService {
	return LedgerService{
		ledgerDao: ledgerDao,
	}
}

func (ls LedgerService) Create(ctx context.Context, tx transaction.Transaction, ledger model.Ledger) error {
	return ls.ledgerDao.Create(ctx, tx, ledger)
}
