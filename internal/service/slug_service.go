package service

import (
	"context"
	"errors"
	"strings"

	"github.com/VanessaVallarini/campaign-consumer-api/internal/model"
	"github.com/google/uuid"
	easyzap "github.com/lockp111/go-easyzap"
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
	if err := s.isValidStatus(slug.Status); err != nil {
		return err
	}

	return s.slugRepository.Upsert(ctx, model.Slug{
		Id:        slug.Id,
		Name:      strings.ToUpper(slug.Name),
		Status:    slug.Status,
		CreatedBy: slug.CreatedBy,
		UpdatedBy: slug.UpdatedBy,
		CreatedAt: slug.CreatedAt,
		UpdatedAt: slug.UpdatedAt,
	})
}

func (s SlugService) Fetch(ctx context.Context, id uuid.UUID) (model.Slug, error) {
	return s.slugRepository.Fetch(ctx, id)
}

func (s SlugService) isValidStatus(status string) error {
	modelStatus := model.SlugStatus(status)
	if modelStatus != model.ActiveSlug && modelStatus != model.InactiveSlug {
		easyzap.Errorf("invalid owner status %s", status)

		return errors.New("Invalid owner status")
	}
	return nil
}
