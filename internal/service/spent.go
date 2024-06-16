package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
)

type SpentRepository interface {
	Upsert(context.Context, transaction.Transaction, model.Spent) error
}

type SpentService struct {
	spentRepository SpentRepository
}

func NewSpentService(spentRepository SpentRepository) SpentService {
	return SpentService{
		spentRepository: spentRepository,
	}
}

func (s SpentService) CreateOrUpdate(ctx context.Context, tx transaction.Transaction, spent model.Spent) error {
	return s.spentRepository.Upsert(ctx, tx, spent)
}
