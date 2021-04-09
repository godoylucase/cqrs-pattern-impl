package business

import (
	"context"
	"fmt"
)

type Repository interface {
	FindByURL(ctx context.Context, url URL) (*BaseArticle, error)
	Create(ctx context.Context, article *BaseArticle) (string, error)
}

type service struct {
	ar Repository
}

func NewArticleService(ar Repository) *service {
	return &service{
		ar: ar,
	}
}

func (s *service) Create(ctx context.Context, ba *BaseArticle) (string, error) {
	article, err := s.ar.FindByURL(ctx, ba.SourceURL)
	if err != nil {
		return "", fmt.Errorf("error when fetching from base article repository: %w", err)
	} else if article != nil {
		return "", nil
	}

	id, err := s.ar.Create(ctx, ba)
	if err != nil {
		return "", fmt.Errorf("error when creating base article: %w", err)
	}

	return id, nil
}
