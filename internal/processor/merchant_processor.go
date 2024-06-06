package processor

import (
	"context"
	"strings"
	"time"

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
		Id:        uuid.MustParse(message.Id),
		OwnerId:   uuid.MustParse(message.OwnerId),
		RegionId:  uuid.MustParse(message.RegionId),
		Slugs:     mp.convertSlugs(message.Slugs),
		Name:      name,
		Status:    model.MerchantStatus(message.Status),
		CreatedBy: message.CreatedBy,
		UpdatedBy: message.UpdatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
