package service

import (
	"context"
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
	slugDao        SlugDao
	slugHistoryDao SlugHistoryDao
	tm             TransactionManager
}

func NewSlugService(slugDao SlugDao, slugHistoryDao SlugHistoryDao, tm TransactionManager) SlugService {
	return SlugService{
		slugDao:        slugDao,
		slugHistoryDao: slugHistoryDao,
		tm:             tm,
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
