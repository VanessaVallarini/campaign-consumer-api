package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
)

type SlugRepository interface {
	Upsert(context.Context, model.Slug) error
	Fetch(context.Context, uuid.UUID) (model.Slug, error)
}

type SlugService struct {
	slugRepository SlugRepository
}

func NewSlugService(slugRepository SlugRepository) SlugService {
	return SlugService{
		slugRepository: slugRepository,
	}
}

func (s SlugService) Upsert(ctx context.Context, slug model.Slug) error {
	return s.slugRepository.Upsert(ctx, slug)
}

func (s SlugService) Fetch(ctx context.Context, id uuid.UUID) (model.Slug, error) {
	return s.slugRepository.Fetch(ctx, id)
}
