package processor

import (
	"context"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
)

type OwnerService interface {
	CreateOrUpdate(context.Context, model.Owner) error
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
	oep.ownerService.CreateOrUpdate(context.Background(), model.Owner{
		Id:        uuid.MustParse(message.Id),
		Email:     message.Email,
		Status:    model.OwnerStatus(message.Status),
		CreatedBy: message.CreatedBy,
		UpdatedBy: message.UpdatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return nil
}
