package processor

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
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
		Id:         message.Id,
		MerchantId: message.MerchantId,
		Status:     model.CampaignStatus(message.Status),
		Lat:        message.Lat,
		Long:       message.Long,
		Budget:     message.Budget,
		CreatedBy:  message.User,
		UpdatedBy:  message.User,
		CreatedAt:  message.EventTime,
		UpdatedAt:  message.EventTime,
	})

	return nil
}
