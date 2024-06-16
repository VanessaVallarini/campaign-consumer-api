package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/jackc/pgx/v5"
)

type SpentRepository interface {
	Upsert(context.Context, pgx.Tx, model.Spent)
}

type LedgerRepository interface {
	Create(context.Context, pgx.Tx, model.Ledger) error
}

type Transaction interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type ClickImpressionService struct {
	spentRepository  SpentRepository
	ledgerRepository LedgerRepository
	transaction      Transaction
}

func NewClickImpressionService(spentRepository SpentRepository,
	ledgerRepository LedgerRepository,
	transaction Transaction) ClickImpressionService {
	return ClickImpressionService{
		spentRepository:  spentRepository,
		ledgerRepository: ledgerRepository,
		transaction:      transaction,
	}
}

func (c CampaignService) ClickImpression(ctx context.Context, clickImpression model.ClickImpression) error {
	//get campaign e valida se é ativa
	//get merchant e valida se é ativo e valida se o slug_id do evento contém nos slugs do merchant
	//get slug e valida se é ativo
	//valida se o lat long da campanha batem com o lat long do merchant (region)
	//pega os custos do slug e da region e soma pra salvar na spent e ledger

	//criar indexes nas tabelas

	return nil
}
