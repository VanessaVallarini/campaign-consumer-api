package service

import (
	"context"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
)

type SlugDao interface {
	Upsert(context.Context, model.Slug) error
	Fetch(context.Context, uuid.UUID) (model.Slug, error)
}

type SlugService struct {
	slugDao SlugDao
}

func NewSlugService(slugDao SlugDao) SlugService {

	return SlugService{
		slugDao: slugDao,
	}
}

func (ss SlugService) Upsert(ctx context.Context, slug model.Slug) error {

	return ss.slugDao.Upsert(ctx, model.Slug{
		Id:        slug.Id,
		Name:      strings.ToUpper(slug.Name),
		Status:    slug.Status,
		Cost:      slug.Cost,
		CreatedBy: slug.CreatedBy,
		UpdatedBy: slug.UpdatedBy,
		CreatedAt: slug.CreatedAt,
		UpdatedAt: slug.UpdatedAt,
	})
}

func (ss SlugService) Fetch(ctx context.Context, id uuid.UUID) (model.Slug, error) {

	return ss.slugDao.Fetch(ctx, id)
}
