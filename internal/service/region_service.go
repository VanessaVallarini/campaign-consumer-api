package service

import (
	"context"
	"errors"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type RegionRepository interface {
	Upsert(context.Context, model.Region) error
	Fetch(context.Context, uuid.UUID) (model.Region, error)
}

type RegionService struct {
	regionRepository RegionRepository
}

func NewRegionService(regionRepository RegionRepository) RegionService {
	return RegionService{
		regionRepository: regionRepository,
	}
}

func (r RegionService) Upsert(ctx context.Context, region model.Region) error {
	if err := r.isValidStatus(region.Status); err != nil {
		return err
	}

	return r.regionRepository.Upsert(ctx, region)
}

func (r RegionService) Fetch(ctx context.Context, id uuid.UUID) (model.Region, error) {
	return r.regionRepository.Fetch(ctx, id)
}

func (r RegionService) isValidStatus(status string) error {
	modelStatus := model.RegionStatus(status)
	if modelStatus != model.ActiveRegion && modelStatus != model.InactiveRegion {
		easyzap.Errorf("invalid region status %s", status)

		return errors.New("Invalid region status")
	}
	return nil
}
