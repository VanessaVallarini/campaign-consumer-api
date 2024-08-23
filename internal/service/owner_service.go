package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
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

	return os.ownerDao.Upsert(ctx, owner)
}
