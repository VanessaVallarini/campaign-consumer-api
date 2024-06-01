package processor

import (
	"context"
	"strconv"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type SlugService interface {
	CreateOrUpdate(context.Context, model.Slug) error
}

type SlugProcessor struct {
	slugService SlugService
}

func NewSlugProcessor(slugService SlugService) SlugProcessor {
	return SlugProcessor{
		slugService: slugService,
	}
}

func (sp SlugProcessor) SlugProcessor(message model.SlugEvent) (returnErr error) {
	active, err := strconv.ParseBool(message.Active)
	if err != nil {
		easyzap.Error(err, "error converting string to bool")

		return
	}

	sp.slugService.CreateOrUpdate(context.Background(), model.Slug{
		Id:        uuid.MustParse(message.Id),
		Name:      message.Name,
		Active:    active,
		Cost:      message.Cost,
		CreatedBy: message.CreatedBy,
		UpdatedBy: message.UpdatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return nil
}
