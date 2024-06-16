package processor

import (
	"context"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type RegionService interface {
	Upsert(context.Context, model.Region) error
}

type RegionProcessor struct {
	regionService RegionService
}

func NewRegionProcessor(regionService RegionService) RegionProcessor {
	return RegionProcessor{
		regionService: regionService,
	}
}

func (rp RegionProcessor) RegionProcessor(message model.RegionEvent) (returnErr error) {
	return rp.regionService.Upsert(context.Background(), model.Region{
		Id:        message.Id,
		Name:      strings.ToUpper(message.Name),
		Status:    model.RegionStatus(message.Status),
		Lat:       message.Lat,
		Long:      message.Long,
		Cost:      message.Cost,
		CreatedBy: message.User,
		UpdatedBy: message.User,
		CreatedAt: message.EventTime,
		UpdatedAt: message.EventTime,
	})
}
