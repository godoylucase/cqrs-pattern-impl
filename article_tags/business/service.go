package business

import (
	"context"
	"fmt"

	"github.com/godoylucase/articles_tags/internal/event"
	"github.com/sirupsen/logrus"
)

type repository interface {
	GetByUserIDAndSourceURL(ctx context.Context, userID string, url URL) (*BaseArticle, error)
	Create(ctx context.Context, article *BaseArticle) (string, error)
}

type eventBroker interface {
	ArticleCreation(article event.Partitionable) error
}

type service struct {
	ar repository
	eb eventBroker
}

func NewArticleService(ar repository, eb eventBroker) *service {
	return &service{
		ar: ar,
		eb: eb,
	}
}

func (s *service) Create(ctx context.Context, ba *BaseArticle) (string, error) {
	article, err := s.ar.GetByUserIDAndSourceURL(ctx, string(ba.UserID), ba.SourceURL)
	if err != nil {
		return "", fmt.Errorf("error when fetching from base article repository: %w", err)
	} else if article != nil {
		return article.ID.Hex(), nil
	}

	id, err := s.ar.Create(ctx, ba)
	if err != nil {
		return "", fmt.Errorf("error when creating base article: %w", err)
	}

	if err = s.eb.ArticleCreation(ba.ToDTO()); err != nil {
		// TODO enhance this log line
		logrus.Warn("event could not be sent")
	}

	return id, nil
}
