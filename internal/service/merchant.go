package service

import (
	"context"
	"errors"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
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

func (m MerchantService) Upsert(ctx context.Context, merchant model.Merchant) error {
	if err := m.isValidStatus(merchant.Status); err != nil {
		return err
	}

	return m.merchantRepository.Upsert(ctx, model.Merchant{
		Id:        merchant.Id,
		OwnerId:   merchant.OwnerId,
		RegionId:  merchant.RegionId,
		Slugs:     merchant.Slugs,
		Name:      strings.ToUpper(merchant.Name),
		Status:    merchant.Status,
		CreatedBy: merchant.CreatedBy,
		UpdatedBy: merchant.UpdatedBy,
		CreatedAt: merchant.CreatedAt,
		UpdatedAt: merchant.UpdatedAt,
	})
}

func (m MerchantService) Fetch(ctx context.Context, id uuid.UUID) (model.Merchant, error) {
	return m.merchantRepository.Fetch(ctx, id)
}

func (m MerchantService) isValidStatus(status string) error {
	modelStatus := model.MerchantStatus(status)
	if modelStatus != model.ActiveMerchant && modelStatus != model.InactiveMerchant {
		easyzap.Errorf("invalid merchant status %s", status)

		return errors.New("Invalid merchant status")
	}
	return nil
}
