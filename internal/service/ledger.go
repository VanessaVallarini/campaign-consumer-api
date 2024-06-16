package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
)

type LedgerRepository interface {
	Create(context.Context, transaction.Transaction, model.Ledger) error
}

type LedgerService struct {
	ledgerRepository LedgerRepository
}

func NewLedgerService(ledgerRepository LedgerRepository) LedgerService {
	return LedgerService{
		ledgerRepository: ledgerRepository,
	}
}

func (s LedgerService) Create(ctx context.Context, tx transaction.Transaction, ledger model.Ledger) error {
	return s.ledgerRepository.Create(ctx, tx, ledger)
}
