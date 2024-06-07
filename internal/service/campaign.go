package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type CampaignRepository interface {
	Upsert(context.Context, model.Campaign) error
}

type CampaignService struct {
	campaignRepository CampaignRepository
}

func NewCampaignService(campaignRepository CampaignRepository) CampaignService {
	return CampaignService{
		campaignRepository: campaignRepository,
	}
}

func (c CampaignService) CreateOrUpdate(ctx context.Context, campaign model.Campaign) error {
	return c.campaignRepository.Upsert(ctx, campaign)
}
