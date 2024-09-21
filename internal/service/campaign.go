package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type CampaignDao interface {
	Upsert(context.Context, model.Campaign) error
	Fetch(context.Context, uuid.UUID) (model.Campaign, error)
}

type CampaignService struct {
	campaignDao CampaignDao
}

func NewCampaignService(campaignDao CampaignDao) CampaignService {
	return CampaignService{
		campaignDao: campaignDao,
	}
}

func (cs CampaignService) Upsert(ctx context.Context, campaign model.Campaign) error {
	err := campaign.ValidateCampaign()
	if err != nil {
		easyzap.Error(err, "upsert merchant fail : %w", err)

		return model.ErrInvalid
	}

	return cs.campaignDao.Upsert(ctx, campaign)
}

func (cs CampaignService) Fetch(ctx context.Context, id uuid.UUID) (model.Campaign, error) {
	return cs.campaignDao.Fetch(ctx, id)
}
