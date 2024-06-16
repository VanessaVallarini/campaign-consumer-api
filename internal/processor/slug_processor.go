package processor

import (
	"context"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
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
	sp.slugService.CreateOrUpdate(context.Background(), model.Slug{
		Id:        message.Id,
		Name:      strings.ToUpper(message.Name),
		Status:    model.SlugStatus(message.Status),
		Cost:      message.Cost,
		CreatedBy: message.User,
		UpdatedBy: message.User,
		CreatedAt: message.EventTime,
		UpdatedAt: message.EventTime,
	})

	return nil
}
