package service

import (
	"context"
	"fmt"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type RegionDao interface {
	Fetch(context.Context, uuid.UUID) (model.Region, error)
	Create(context.Context, model.Region) error
	Update(context.Context, model.Region) error
}

type RegionHistoryDao interface {
	Create(context.Context, model.RegionHistory) error
}

type RegionService struct {
	regionDao        RegionDao
	regionHistoryDao RegionHistoryDao
}

func NewRegionService(regionDao RegionDao, regionHistoryDao RegionHistoryDao) RegionService {
	return RegionService{
		regionDao:        regionDao,
		regionHistoryDao: regionHistoryDao,
	}
}

func (rs RegionService) Upsert(ctx context.Context, region model.Region) error {
	err := region.ValidateRegion()
	if err != nil {
		easyzap.Error(err, "invalid region")

		return model.ErrInvalid
	}

	regionDb, err := rs.regionDao.Fetch(ctx, region.Id)
	if err != nil && err != model.ErrNotFound {
		easyzap.Errorf("fail to fetch region by regionId %s: %v", region.Id.String(), err)

		return err
	}

	if err != nil && err == model.ErrNotFound {
		err := rs.regionDao.Create(ctx, region)
		if err != nil {
			easyzap.Errorf("fail to create region %v: %v", region, err)

			return err
		}
	} else {
		err = rs.regionDao.Update(ctx, region)
		if err != nil {
			easyzap.Errorf("fail to update regionDb %v to region %v: %v", regionDb, region, err)

			return err
		}
	}

	err = rs.registryHistory(ctx, region, &regionDb)
	if err != nil {
		easyzap.Errorf("fail to registry history regionDb %v to region %v: %v", regionDb, region, err)
		region.Status = string(model.Cancelled)
		errRollback := rs.regionDao.Update(ctx, regionDb)
		if errRollback != nil {
			easyzap.Errorf("[INCONSISTENT] fail to rollback region %v to regionDb %v: %v", region, regionDb, err)

			return err
		}

		return err
	}

	return nil
}

func (rs RegionService) registryHistory(ctx context.Context, region model.Region, regionDb *model.Region) error {
	if regionDb.Id == uuid.Nil {
		err := rs.regionHistoryDao.Create(ctx, model.RegionHistory{
			Id:          uuid.New(),
			RegionId:    region.Id,
			Status:      region.Status,
			Description: model.RegionCreatedAndActive,
			CreatedBy:   region.UpdatedBy,
			CreatedAt:   region.UpdatedAt,
		})
		if err != nil {
			easyzap.Errorf("fail to registry history region create %v: %v", region, region)

			return err
		}
	} else {
		if regionDb.Status != region.Status {
			err := rs.regionHistoryDao.Create(ctx, model.RegionHistory{
				Id:          uuid.New(),
				RegionId:    region.Id,
				Status:      region.Status,
				Description: fmt.Sprintf(model.RegionUpdateStatus, regionDb.Status, region.Status),
				CreatedBy:   region.UpdatedBy,
				CreatedAt:   region.UpdatedAt,
			})
			if err != nil {
				easyzap.Errorf("fail to registry history region status from %s to %s for regionId %s: %v", regionDb.Status, region.Status, region.Id.String(), err)

				return err
			}
		}

		if regionDb.Cost != region.Cost {
			err := rs.regionHistoryDao.Create(ctx, model.RegionHistory{
				Id:          uuid.New(),
				RegionId:    region.Id,
				Status:      region.Status,
				Description: fmt.Sprintf(model.RegionUpdateCost, regionDb.Cost, region.Cost),
				CreatedBy:   region.UpdatedBy,
				CreatedAt:   region.UpdatedAt,
			})
			if err != nil {
				easyzap.Errorf("fail to registry history region cost from %.2f to %.2f for regionId %s: %v", regionDb.Cost, region.Cost, region.Id.String(), err)

				return err
			}
		}
	}

	return nil
}
