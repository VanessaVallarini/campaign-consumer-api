package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type MerchantDao interface {
	Upsert(context.Context, model.Merchant) error
	Fetch(context.Context, uuid.UUID) (model.Merchant, error)
}

type MerchantService struct {
	merchantDao MerchantDao
}

func NewMerchantService(merchantDao MerchantDao) MerchantService {
	return MerchantService{
		merchantDao: merchantDao,
	}
}

func (ms MerchantService) Upsert(ctx context.Context, merchant model.Merchant) error {
	err := merchant.ValidateMerchant()
	if err != nil {
		easyzap.Error("upsert merchant for merchant id %s fail: %v", merchant.Id.String(), err)

		return model.ErrInvalid
	}

	return ms.merchantDao.Upsert(ctx, merchant)
}

func (ms MerchantService) Fetch(ctx context.Context, id uuid.UUID) (model.Merchant, error) {

	return ms.merchantDao.Fetch(ctx, id)
}
