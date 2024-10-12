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

type SlugDao interface {
	Fetch(context.Context, uuid.UUID) (model.Slug, error)
	Create(context.Context, transaction.Transaction, model.Slug) error
	Update(context.Context, transaction.Transaction, model.Slug) error
}

type SlugHistoryDao interface {
	Create(context.Context, transaction.Transaction, model.SlugHistory) error
}

type SlugService struct {
	slugDao            SlugDao
	slugHistoryDao     SlugHistoryDao
	tm                 TransactionManager
	slugRedisValidator RedisValidator
}

func NewSlugService(slugDao SlugDao, slugHistoryDao SlugHistoryDao, tm TransactionManager, slugRedisValidator RedisValidator) SlugService {
	return SlugService{
		slugDao:            slugDao,
		slugHistoryDao:     slugHistoryDao,
		tm:                 tm,
		slugRedisValidator: slugRedisValidator,
	}
}

func (ss SlugService) Upsert(ctx context.Context, slug model.Slug) error {
	err := slug.ValidateSlug()
	if err != nil {
		easyzap.Error(err, "invalid slug")

		return model.ErrInvalid
	}

	slugDb, err := ss.slugDao.Fetch(ctx, slug.Id)
	if err != nil && err != model.ErrNotFound {
		easyzap.Errorf("fail to fetch slug by slugId %s: %v", slug.Id.String(), err)

		return err
	}

	if err != nil && err == model.ErrNotFound {
		err := ss.createAndRegistryHistory(ctx, slug, &slugDb)
		if err != nil {
			easyzap.Errorf("fail to create slug %v: %v", slug, err)

			return err
		}
	} else {
		err := ss.updateAndRegistryHistory(ctx, slug, &slugDb)
		if err != nil {
			easyzap.Errorf("fail to registry history slugDb %v to slug %v: %v", slugDb, slug, err)

			return err
		}
	}

	return nil
}

func (ss SlugService) createAndRegistryHistory(ctx context.Context, slug model.Slug, slugDb *model.Slug) error {
	funcWithTransaction := func(ctx context.Context, tx transaction.Transaction) error {
		err := ss.slugDao.Create(ctx, tx, slug)
		if err != nil {

			return err
		}

		history := ss.buildHistory(slug, slugDb)

		err = ss.slugHistoryDao.Create(ctx, tx, history)
		if err != nil {

			return err
		}

		return err
	}

	return ss.tm.Execute(ctx, funcWithTransaction)
}

func (ss SlugService) updateAndRegistryHistory(ctx context.Context, slug model.Slug, slugDb *model.Slug) error {
	funcWithTransaction := func(ctx context.Context, tx transaction.Transaction) error {
		err := ss.slugDao.Update(ctx, tx, slug)
		if err != nil {

			return err
		}

		history := ss.buildHistory(slug, slugDb)

		err = ss.slugHistoryDao.Create(ctx, tx, history)
		if err != nil {

			return err
		}

		return err
	}

	return ss.tm.Execute(ctx, funcWithTransaction)
}

func (ss SlugService) buildHistory(slug model.Slug, slugDb *model.Slug) model.SlugHistory {
	history := model.SlugHistory{
		Id:          uuid.New(),
		SlugId:      slug.Id,
		Status:      slug.Status,
		Description: model.SlugCreatedAndActive,
		CreatedBy:   slug.UpdatedBy,
		CreatedAt:   slug.UpdatedAt,
	}

	if slugDb.Id == uuid.Nil {
		history.Description = model.SlugCreatedAndActive
	} else {
		if slugDb.Status != slug.Status {
			history.Description = fmt.Sprintf(model.SlugUpdateStatus, slugDb.Status, slug.Status)
		}
		if slugDb.Cost != slug.Cost {
			history.Description = fmt.Sprintf(model.SlugUpdateCost, slugDb.Cost, slug.Cost)
		}
		if slugDb.Cost != slug.Cost && slugDb.Status != slug.Status {
			history.Description = fmt.Sprintf("%v E %v",
				fmt.Sprintf(model.SlugUpdateCost, slugDb.Cost, slug.Cost),
				fmt.Sprintf(model.SlugUpdateStatus, slugDb.Status, slug.Status))
		}
	}

	return history
}

// Fetch slug iin redis
func (ss SlugService) Fetch(ctx context.Context, name string) (model.Slug, error) {

	value, err := ss.slugRedisValidator.Get(ctx, ss.uniqueKey(name))
	if err != nil {

		return model.Slug{}, err
	}

	slugJSON, err := json.Marshal(value)
	if err != nil {

		return model.Slug{}, fmt.Errorf("failed to marshal slug: %w", err)
	}

	var slug model.Slug
	err = json.Unmarshal(slugJSON, &slug)
	if err != nil {

		return model.Slug{}, fmt.Errorf("failed to unmarshal into slug: %w", err)
	}

	return slug, nil
}

// Operation performed at dawn
// Postgres is loaded into Redis
func (ss SlugService) SaveInRedis(ctx context.Context, slug model.Slug) (bool, error) {

	slugJSON, err := json.Marshal(slug)
	if err != nil {

		return false, fmt.Errorf("failed to marshal slug: %w", err)
	}

	ok, err := ss.slugRedisValidator.SetIfNotExists(ctx, ss.uniqueKey(slug.Name), slugJSON)
	if err != nil {

		return ok, err
	}

	return ok, nil
}

func (ss SlugService) uniqueKey(slugName string) string {

	return fmt.Sprintf("%s::slug", slugName)
}
