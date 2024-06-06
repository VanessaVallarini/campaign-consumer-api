package processor

import (
	"context"
	"strings"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
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
	name := strings.ToUpper(message.Name)

	sp.slugService.CreateOrUpdate(context.Background(), model.Slug{
		Id:        uuid.MustParse(message.Id),
		Name:      name,
		Status:    model.SlugStatus(message.Status),
		Cost:      message.Cost,
		CreatedBy: message.CreatedBy,
		UpdatedBy: message.UpdatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return nil
}
