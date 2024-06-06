package processor

import (
	"context"
	"strings"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
)

type RegionService interface {
	CreateOrUpdate(context.Context, model.Region) error
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
	name := strings.ToUpper(message.Name)

	rp.regionService.CreateOrUpdate(context.Background(), model.Region{
		Id:        uuid.MustParse(message.Id),
		Name:      name,
		Status:    model.RegionStatus(message.Status),
		Lat:       message.Lat,
		Long:      message.Long,
		Cost:      message.Cost,
		CreatedBy: message.CreatedBy,
		UpdatedBy: message.UpdatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return nil
}
