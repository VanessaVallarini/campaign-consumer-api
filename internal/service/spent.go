package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/google/uuid"
)

type SpentRepository interface {
	Upsert(context.Context, transaction.Transaction, model.Spent) error
	FetchByCampaignIdAndBucket(context.Context, uuid.UUID, string) (model.Spent, error)
}

type SpentService struct {
	spentRepository SpentRepository
}

func NewSpentService(spentRepository SpentRepository) SpentService {
	return SpentService{
		spentRepository: spentRepository,
	}
}

func (s SpentService) Upsert(ctx context.Context, tx transaction.Transaction, spent model.Spent) error {
	return s.spentRepository.Upsert(ctx, tx, spent)
}

func (s SpentService) FetchByCampaignIdAndBucket(ctx context.Context, id uuid.UUID, bucket string) (model.Spent, error) {
	return s.spentRepository.FetchByCampaignIdAndBucket(ctx, id, bucket)
}
