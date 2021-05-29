package business

import "github.com/godoylucase/read_tags/business/dto"

type repository interface {
	GetArticleByGlobalTags(globalHashTags []string) (dto.ArticleByGlobalHashTagRead, error)
	GetUsersByArticle(articleID string) ([]dto.UserByArticle, error)
}

type service struct {
	repository repository
}

func NewQueryService(repository repository) *service {
	return &service{repository: repository}
}

func (s *service) GetArticleByGlobalTags(globalHashTags []string) (dto.ArticleByGlobalHashTagRead, error) {
	articles, err := s.repository.GetArticleByGlobalTags(globalHashTags)
	if err != nil {
		return dto.ArticleByGlobalHashTagRead{}, err
	}

	return articles, nil
}

func (s *service) GetUsersByArticle(articleID string) ([]dto.UserByArticle, error) {
	users, err := s.repository.GetUsersByArticle(articleID)
	if err != nil {
		return nil, err
	}

	return users, nil
}
