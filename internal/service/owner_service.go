package service

import (
	"context"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type OwnerDao interface {
	Upsert(context.Context, model.Owner) error
}

type OwnerService struct {
	ownerRepository OwnerDao
}

func NewOwnerService(ownerRepository OwnerDao) OwnerService {
	return OwnerService{
		ownerRepository: ownerRepository,
	}
}

func (os OwnerService) Upsert(ctx context.Context, owner model.Owner) error {
	return os.ownerRepository.Upsert(ctx, model.Owner{
		Id:        owner.Id,
		Email:     strings.ToUpper(owner.Email),
		Status:    owner.Status,
		CreatedBy: owner.CreatedBy,
		UpdatedBy: owner.UpdatedBy,
		CreatedAt: owner.CreatedAt,
		UpdatedAt: owner.UpdatedAt,
	})
}
