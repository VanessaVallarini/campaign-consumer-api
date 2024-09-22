package service

import (
	"context"
	"fmt"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
)

type SlugDao interface {
	Fetch(ctx context.Context, id uuid.UUID) (model.Slug, error)
	Create(ctx context.Context, slug model.Slug) error
	Update(ctx context.Context, slug model.Slug) error
}

type SlugHistoryDao interface {
	Create(context.Context, model.SlugHistory) error
}

type SlugService struct {
	slugDao        SlugDao
	slugHistoryDao SlugHistoryDao
}

func NewSlugService(slugDao SlugDao, slugHistoryDao SlugHistoryDao) SlugService {

	return SlugService{
		slugDao:        slugDao,
		slugHistoryDao: slugHistoryDao,
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
		err := ss.slugDao.Create(ctx, slug)
		if err != nil {
			easyzap.Errorf("fail to create slug %v: %v", slug, err)

			return err
		}
	} else {
		err = ss.slugDao.Update(ctx, slug)
		if err != nil {
			easyzap.Errorf("fail to update slugDb %v to slug %v: %v", slugDb, slug, err)

			return err
		}
	}

	err = ss.registryHistory(ctx, slug, &slugDb)
	if err != nil {
		easyzap.Errorf("fail to registry history slugDb %v to slug %v: %v", slugDb, slug, err)
		slug.Status = string(model.Cancelled)
		errRollback := ss.slugDao.Update(ctx, slugDb)
		if errRollback != nil {
			easyzap.Errorf("[INCONSISTENT] fail to rollback slug %v to slugDb %v: %v", slug, slugDb, err)

			return err
		}

		return err
	}

	return nil
}

func (ss SlugService) registryHistory(ctx context.Context, slug model.Slug, slugDb *model.Slug) error {
	if slugDb.Id == uuid.Nil {
		err := ss.slugHistoryDao.Create(ctx, model.SlugHistory{
			Id:          uuid.New(),
			SlugId:      slug.Id,
			Status:      slug.Status,
			Description: model.SlugCreatedAndActive,
			CreatedBy:   slug.UpdatedBy,
			CreatedAt:   slug.UpdatedAt,
		})
		if err != nil {
			easyzap.Errorf("fail to registry history slug create %v: %v", slug, err)

			return err
		}
	} else {
		if slugDb.Status != slug.Status {
			err := ss.slugHistoryDao.Create(ctx, model.SlugHistory{
				Id:          uuid.New(),
				SlugId:      slug.Id,
				Status:      slug.Status,
				Description: fmt.Sprintf(model.SlugUpdateStatus, slugDb.Status, slug.Status),
				CreatedBy:   slug.UpdatedBy,
				CreatedAt:   slug.UpdatedAt,
			})
			if err != nil {
				easyzap.Errorf("fail to registry history slug status from %s to %s for slugId %s: %v", slugDb.Status, slug.Status, slug.Id.String(), err)

				return err
			}
		}

		if slugDb.Cost != slug.Cost {
			err := ss.slugHistoryDao.Create(ctx, model.SlugHistory{
				Id:          uuid.New(),
				SlugId:      slug.Id,
				Status:      slug.Status,
				Description: fmt.Sprintf(model.SlugUpdateCost, slugDb.Cost, slug.Cost),
				CreatedBy:   slug.UpdatedBy,
				CreatedAt:   slug.UpdatedAt,
			})
			if err != nil {
				easyzap.Errorf("fail to registry history slug cost from %.2f to %.2f for slugId %s: %v", slugDb.Cost, slug.Cost, slug.Id.String(), err)

				return err
			}
		}
	}

	return nil
}
