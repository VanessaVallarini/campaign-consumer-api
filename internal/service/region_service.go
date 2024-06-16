package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
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

func (r RegionService) CreateOrUpdate(ctx context.Context, region model.Region) error {
	return r.regionRepository.Upsert(ctx, region)
}

func (r RegionService) Fetch(ctx context.Context, id uuid.UUID) (model.Region, error) {
	return r.regionRepository.Fetch(ctx, id)
}
