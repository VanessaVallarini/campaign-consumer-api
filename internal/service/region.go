package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/VanessaVallarini/campaign-consumer-api/internal/pkg/transaction"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type RegionDao interface {
	Fetch(context.Context, uuid.UUID) (model.Region, error)
	Create(context.Context, transaction.Transaction, model.Region) error
	Update(context.Context, transaction.Transaction, model.Region) error
}

type RegionHistoryDao interface {
	Create(context.Context, transaction.Transaction, model.RegionHistory) error
}

type RegionService struct {
	regionDao             RegionDao
	regionHistoryDao      RegionHistoryDao
	tm                    TransactionManager
	regioinRedisValidator RedisValidator
}

func NewRegionService(regionDao RegionDao, regionHistoryDao RegionHistoryDao, tm TransactionManager, regioinRedisValidator RedisValidator) RegionService {
	return RegionService{
		regionDao:             regionDao,
		regionHistoryDao:      regionHistoryDao,
		tm:                    tm,
		regioinRedisValidator: regioinRedisValidator,
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
		err := rs.createAndRegistryHistory(ctx, region, &regionDb)
		if err != nil {
			easyzap.Errorf("fail to create region %v: %v", region, err)

			return err
		}
	} else {
		err := rs.updateAndRegistryHistory(ctx, region, &regionDb)
		if err != nil {
			easyzap.Errorf("fail to update regionDb %v to region %v: %v", regionDb, region, err)

			return err
		}
	}

	return nil
}

func (rs RegionService) createAndRegistryHistory(ctx context.Context, region model.Region, regionDb *model.Region) error {
	funcWithTransaction := func(ctx context.Context, tx transaction.Transaction) error {
		err := rs.regionDao.Create(ctx, tx, region)
		if err != nil {

			return err
		}

		history := rs.buildHistory(region, regionDb)

		err = rs.regionHistoryDao.Create(ctx, tx, history)
		if err != nil {

			return err
		}

		return err
	}

	return rs.tm.Execute(ctx, funcWithTransaction)
}

func (rs RegionService) updateAndRegistryHistory(ctx context.Context, region model.Region, regionDb *model.Region) error {
	funcWithTransaction := func(ctx context.Context, tx transaction.Transaction) error {
		err := rs.regionDao.Update(ctx, tx, region)
		if err != nil {

			return err
		}

		history := rs.buildHistory(region, regionDb)

		err = rs.regionHistoryDao.Create(ctx, tx, history)
		if err != nil {

			return err
		}

		return err
	}

	return rs.tm.Execute(ctx, funcWithTransaction)
}

func (rs RegionService) buildHistory(region model.Region, regionDb *model.Region) model.RegionHistory {
	history := model.RegionHistory{
		Id:          uuid.New(),
		RegionId:    region.Id,
		Status:      region.Status,
		Description: model.RegionCreatedAndActive,
		CreatedBy:   region.UpdatedBy,
		CreatedAt:   region.UpdatedAt,
	}

	if regionDb.Id == uuid.Nil {
		history.Description = model.RegionCreatedAndActive
	} else {
		if regionDb.Status != region.Status {
			history.Description = fmt.Sprintf(model.RegionUpdateStatus, regionDb.Status, region.Status)
		}

		if regionDb.Cost != region.Cost {
			history.Description = fmt.Sprintf(model.RegionUpdateCost, regionDb.Cost, region.Cost)
		}

		if regionDb.Cost != region.Cost && regionDb.Status != region.Status {
			history.Description = fmt.Sprintf("%v E %v",
				fmt.Sprintf(model.RegionUpdateCost, regionDb.Cost, region.Cost),
				fmt.Sprintf(model.RegionUpdateStatus, regionDb.Status, region.Status))
		}
	}

	return history
}

// Fetch region iin redis
func (rs RegionService) Fetch(ctx context.Context, regionId uuid.UUID) (model.Region, error) {

	value, err := rs.regioinRedisValidator.Get(ctx, rs.uniqueKey(regionId))
	if err != nil {

		return model.Region{}, err
	}

	regionJSON, err := json.Marshal(value)
	if err != nil {

		return model.Region{}, fmt.Errorf("failed to marshal region: %w", err)
	}

	var region model.Region
	err = json.Unmarshal(regionJSON, &region)
	if err != nil {

		return model.Region{}, fmt.Errorf("failed to unmarshal into region: %w", err)
	}

	return region, nil
}

// Operation performed at dawn
// Postgres is loaded into Redis
func (rs RegionService) SaveInRedis(ctx context.Context, region model.Region) (bool, error) {

	regionJSON, err := json.Marshal(region)
	if err != nil {

		return false, fmt.Errorf("failed to marshal region: %w", err)
	}

	ok, err := rs.regioinRedisValidator.SetIfNotExists(ctx, rs.uniqueKey(region.Id), regionJSON)
	if err != nil {

		return ok, err
	}

	return ok, nil
}

func (rs RegionService) uniqueKey(regionId uuid.UUID) string {

	return fmt.Sprintf("%s::region", regionId.String())
}
