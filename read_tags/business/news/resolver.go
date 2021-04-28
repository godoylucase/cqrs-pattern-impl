package news

import (
	"github.com/godoylucase/read_tags/business/dto"
	"github.com/godoylucase/read_tags/external/event"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type articleRepository interface {
	SaveArticleByGlobalTags(dto dto.ArticleDTO) error
}

type articleResolver struct {
	ar articleRepository
}

func NewArticleResolver(ar articleRepository) *articleResolver {
	return &articleResolver{ar: ar}
}

func (r *articleResolver) Run(ec event.Composite) {
	logrus.Infof("received event with payload %v", ec)

	var adto dto.ArticleDTO
	if err := mapstructure.Decode(ec.Event.Data, &adto); err != nil {
		logrus.Errorf("error converting event data into ArticleDTO with error: %v", err)
		return
	}

	if err := r.ar.SaveArticleByGlobalTags(adto); err != nil {
		logrus.Errorf("saving the ArticleDTO data into storage with error %v", err)
	}
}

// TODO this will be replaced by a real DB operation
type ARMock struct {
}

func (arm *ARMock) SaveArticleByGlobalTags(articleDto dto.ArticleDTO) error {
	logrus.Infof("saving article dto with data %v", articleDto)
	return nil
}
