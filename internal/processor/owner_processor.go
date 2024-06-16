package processor

import (
	"context"
	"strings"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
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
		Id:        message.Id,
		Email:     strings.ToUpper(message.Email),
		Status:    model.OwnerStatus(message.Status),
		CreatedBy: message.User,
		UpdatedBy: message.User,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return nil
}
