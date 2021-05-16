package business

import (
	"context"
	"fmt"

	"github.com/godoylucase/articles_tags/internal"
	"github.com/godoylucase/articles_tags/internal/event"
	"github.com/sirupsen/logrus"
)

type repository interface {
	GetByUserIDAndSourceURL(ctx context.Context, userID string, url URL) (*BaseArticle, error)
	Create(ctx context.Context, article *BaseArticle) (string, error)
	Get(ctx context.Context, id string) (*BaseArticle, error)
	Update(ctx context.Context, id string, ba *BaseArticle) error
}

type eventBroker interface {
	ArticleCreation(article event.Partitionable) error
	ArticleUpdate(article event.Partitionable) error
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

func (s *service) Get(ctx context.Context, id string) (*BaseArticle, error) {
	got, err := s.ar.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error when getting article with ID %v and error: %w", id, err)
	} else if got == nil {
		return nil, internal.ErrResourceNotFound
	}
	return got, nil
}

func (s *service) Update(ctx context.Context, id string, newSnapshot *BaseArticle) error {
	oldSnapshot, err := s.ar.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("error when fetching from base article repository with ID %v and error: %w", id, err)
	} else if oldSnapshot == nil {
		return internal.ErrResourceNotFound
	}

	oldSnapshot.GlobalHashTags = newSnapshot.GlobalHashTags
	oldSnapshot.Paragraphs = newSnapshot.Paragraphs
	oldSnapshot.Title = newSnapshot.Title

	if err := s.ar.Update(ctx, oldSnapshot.ID.Hex(), newSnapshot); err != nil {
		return err
	}

	if err = s.eb.ArticleUpdate(newSnapshot.ToDTO()); err != nil {
		// TODO enhance this log line
		logrus.Warn("event could not be sent")
	}

	return nil
}
