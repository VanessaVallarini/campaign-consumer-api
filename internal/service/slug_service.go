package service

import (
	"context"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
)

type SlugRepository interface {
	Upsert(context.Context, model.Slug) error
}

type SlugService struct {
	slugRepository SlugRepository
}

func NewSlugService(slugRepository SlugRepository) SlugService {
	return SlugService{
		slugRepository: slugRepository,
	}
}

func (s SlugService) CreateOrUpdate(ctx context.Context, slug model.Slug) error {
	return s.slugRepository.Upsert(ctx, slug)
}
