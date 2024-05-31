package processor

import (
	"context"
	"strconv"
	"time"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
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
	active, err := strconv.ParseBool(message.Active)
	if err != nil {
		easyzap.Error(err, "error converting string to bool")

		return
	}

	oep.ownerService.CreateOrUpdate(context.Background(), model.Owner{
		Id:        uuid.MustParse(message.Id),
		Email:     message.Email,
		Active:    active,
		CreatedBy: message.CreatedBy,
		UpdatedBy: message.UpdatedBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	return nil
}
