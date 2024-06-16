package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
)

type MerchantRepository interface {
	Upsert(context.Context, model.Merchant) error
	Fetch(context.Context, uuid.UUID) (model.Merchant, error)
}

type MerchantService struct {
	merchantRepository MerchantRepository
}

func NewMerchantService(merchantRepository MerchantRepository) MerchantService {
	return MerchantService{
		merchantRepository: merchantRepository,
	}
}

func (m MerchantService) CreateOrUpdate(ctx context.Context, merchant model.Merchant) error {
	return m.merchantRepository.Upsert(ctx, merchant)
}

func (m MerchantService) Fetch(ctx context.Context, id uuid.UUID) (model.Merchant, error) {
	return m.merchantRepository.Fetch(ctx, id)
}
