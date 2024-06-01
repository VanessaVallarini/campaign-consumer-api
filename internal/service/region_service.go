package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type RegionRepository interface {
	Upsert(context.Context, model.Region) error
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
