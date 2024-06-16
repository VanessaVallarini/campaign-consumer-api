package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type OwnerRepository interface {
	Upsert(context.Context, model.Owner) error
}

type OwnerService struct {
	ownerRepository OwnerRepository
}

func NewOwnerService(ownerRepository OwnerRepository) OwnerService {
	return OwnerService{
		ownerRepository: ownerRepository,
	}
}

func (o OwnerService) Upsert(ctx context.Context, owner model.Owner) error {
	return o.ownerRepository.Upsert(ctx, owner)
}
