package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type MerchantRepository interface {
	Upsert(context.Context, model.Merchant) error
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
