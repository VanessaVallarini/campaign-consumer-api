package processor

import (
	"context"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
)

type CampaignService interface {
	CreateOrUpdate(context.Context, model.Campaign) error
}

type CampaignProcessor struct {
	campaignService CampaignService
}

func NewCampaignProcessor(campaignService CampaignService) CampaignProcessor {
	return CampaignProcessor{
		campaignService: campaignService,
	}
}

func (cp CampaignProcessor) CampaignProcessor(message model.CampaignEvent) (returnErr error) {
	cp.campaignService.CreateOrUpdate(context.Background(), model.Campaign{
		Id:         uuid.MustParse(message.Id),
		MerchantId: uuid.MustParse(message.MerchantId),
		Status:     model.CampaignStatus(message.Status),
		Lat:        message.Lat,
		Long:       message.Long,
		Budget:     message.Budget,
		CreatedBy:  message.CreatedBy,
		UpdatedBy:  message.UpdatedBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	return nil
}
