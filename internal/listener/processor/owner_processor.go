package processor

import (
	"context"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type OwnerService interface {
	Upsert(context.Context, model.Owner) error
}

type OwnerProcessor struct {
	ownerService OwnerService
}

func NewOwnerProcessor(ownerService OwnerService) OwnerProcessor {
	return OwnerProcessor{
		ownerService: ownerService,
	}
}

func (oep OwnerProcessor) OwnerProcessor(message model.OwnerEvent) (returnErr error) {
	return oep.ownerService.Upsert(context.Background(), model.Owner{
		Id:        message.Id,
		Email:     strings.ToUpper(message.Email),
		Status:    model.OwnerStatus(message.Status),
		CreatedBy: message.User,
		UpdatedBy: message.User,
		CreatedAt: message.EventTime,
		UpdatedAt: message.EventTime,
	})
}
