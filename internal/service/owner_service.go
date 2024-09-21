package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	easyzap "github.com/lockp111/go-easyzap"
)

type OwnerDao interface {
	Upsert(context.Context, model.Owner) error
}

type OwnerService struct {
	ownerDao OwnerDao
}

func NewOwnerService(ownerDao OwnerDao) OwnerService {

	return OwnerService{
		ownerDao: ownerDao,
	}
}

func (os OwnerService) Upsert(ctx context.Context, owner model.Owner) error {
	err := owner.ValidateOwner()
	if err != nil {
		easyzap.Error(err, "upsert owner fail : %w", err)

		return model.ErrInvalid
	}

	return os.ownerDao.Upsert(ctx, owner)
}
