package processor

import (
	"context"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
)

type MerchantService interface {
	CreateOrUpdate(context.Context, model.Merchant) error
}

type MerchantProcessor struct {
	merchantService MerchantService
}

func NewMerchantProcessor(merchantService MerchantService) MerchantProcessor {
	return MerchantProcessor{
		merchantService: merchantService,
	}
}

func (mp MerchantProcessor) MerchantProcessor(message model.MerchantEvent) (returnErr error) {
	name := strings.ToUpper(message.Name)

	mp.merchantService.CreateOrUpdate(context.Background(), model.Merchant{
		Id:        message.Id,
		OwnerId:   message.OwnerId,
		RegionId:  message.RegionId,
		Slugs:     mp.convertSlugs(message.Slugs),
		Name:      strings.ToUpper(name),
		Status:    model.MerchantStatus(message.Status),
		CreatedBy: message.User,
		UpdatedBy: message.User,
		CreatedAt: message.EventTime,
		UpdatedAt: message.EventTime,
	})

	return nil
}

func (mp MerchantProcessor) convertSlugs(slugs []string) []uuid.UUID {
	var slugsUuid []uuid.UUID
	for _, slug := range slugs {
		slugsUuid = append(slugsUuid, uuid.MustParse(slug))
	}

	return slugsUuid

}
