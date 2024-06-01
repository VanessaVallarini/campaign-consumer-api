package processor

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
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
	active, err := strconv.ParseBool(message.Active)
	if err != nil {
		easyzap.Error(err, "error converting string to bool")

		return
	}

	name := strings.ToUpper(message.Name)

	rp.regionService.CreateOrUpdate(context.Background(), model.Region{
		Id:        uuid.MustParse(message.Id),
		Name:      name,
		Active:    active,
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
