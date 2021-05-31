package business

import "github.com/godoylucase/read_tags/business/dto"

type repository interface {
	ArticleByGlobalTags(globalHashTags []string) (dto.ArticleByGlobalHashTagRead, error)
	UserArticlesBySourceURL(articleID string) (dto.UserArticlesBySourceURLRead, error)
}

type service struct {
	repository repository
}

func NewQueryService(repository repository) *service {
	return &service{repository: repository}
}

func (s *service) GetArticleByGlobalTags(globalHashTags []string) (dto.ArticleByGlobalHashTagRead, error) {
	articles, err := s.repository.ArticleByGlobalTags(globalHashTags)
	if err != nil {
		return dto.ArticleByGlobalHashTagRead{}, err
	}

	return articles, nil
}

func (s *service) GetUserArticlesBySourceURL(sourceURL string) (dto.UserArticlesBySourceURLRead, error) {
	ua, err := s.repository.UserArticlesBySourceURL(sourceURL)
	if err != nil {
		return dto.UserArticlesBySourceURLRead{}, err
	}

	return ua, nil
}
