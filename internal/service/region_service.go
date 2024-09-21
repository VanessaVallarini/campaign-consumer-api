package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type RegionDao interface {
	Upsert(context.Context, model.Region) error
	Fetch(context.Context, uuid.UUID) (model.Region, error)
}

type RegionService struct {
	regionDao RegionDao
}

func NewRegionService(regionDao RegionDao) RegionService {
	return RegionService{
		regionDao: regionDao,
	}
}

func (rs RegionService) Upsert(ctx context.Context, region model.Region) error {
	err := region.ValidateRegion()
	if err != nil {
		easyzap.Error(err, "upsert region fail : %w", err)

		return model.ErrInvalid
	}

	return rs.regionDao.Upsert(ctx, region)
}

func (rs RegionService) Fetch(ctx context.Context, id uuid.UUID) (model.Region, error) {
	return rs.regionDao.Fetch(ctx, id)
}
