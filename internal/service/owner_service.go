package service

import (
	"context"
	"errors"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	easyzap "github.com/lockp111/go-easyzap"
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
	if err := o.isValidStatus(owner.Status); err != nil {
		return err
	}

	return o.ownerRepository.Upsert(ctx, model.Owner{
		Id:        owner.Id,
		Email:     strings.ToUpper(owner.Email),
		Status:    owner.Status,
		CreatedBy: owner.CreatedBy,
		UpdatedBy: owner.UpdatedBy,
		CreatedAt: owner.CreatedAt,
		UpdatedAt: owner.UpdatedAt,
	})
}

func (o OwnerService) isValidStatus(status string) error {
	modelStatus := model.OwnerStatus(status)
	if modelStatus != model.ActiveOwner && modelStatus != model.InactiveOwner {
		easyzap.Errorf("invalid owner status %s", status)

		return errors.New("Invalid owner status")
	}
	return nil
}
