package news

import (
	"github.com/godoylucase/read_tags/business/dto"
	"github.com/godoylucase/read_tags/external/event"
	"github.com/sirupsen/logrus"
)

type articleRepository interface {
	SaveArticleByGlobalTags(dto dto.ArticleDTO) error
}

type articleResolver struct {
	ar       articleRepository
}

func NewArticleResolver(ar articleRepository) *articleResolver {
	return &articleResolver{ar: ar}
}

func (r *articleResolver) Run(done <-chan interface{}, events <-chan event.Composite) {
	for {
		select {
		case ec := <-events:
			logrus.Infof("received event with payload %v", ec)

			aDto, ok := ec.Event.Data.(dto.ArticleDTO)
			if !ok {
				logrus.Errorf("converting event data into ArticleDTO, event data value %v", ec.Event.Data)
			}

			if err := r.ar.SaveArticleByGlobalTags(aDto); err != nil {
				logrus.Errorf("saving the ArticleDTO data into storage with error %v", err)
			}
		case <-done:
			return
		}
	}
}

// TODO this will be replaced by a real DB operation
type ARMock struct {
}

func (arm *ARMock) SaveArticleByGlobalTags(articleDto dto.ArticleDTO) error {
	logrus.Infof("saving article dto with data %v", articleDto)
	return nil
}
