package news

import (
	"fmt"

	"github.com/godoylucase/read_tags/business/dto"
	"github.com/godoylucase/read_tags/external/event"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
)

type articleRepository interface {
	SaveArticleByGlobalTags(dto dto.ArticleByGlobalHashTag) error
	SaveUserByArticle(dto dto.UserByArticle) error
}

type articleResolver struct {
	ar articleRepository
}

func NewArticleResolver(ar articleRepository) *articleResolver {
	return &articleResolver{ar: ar}
}

func (r *articleResolver) Run(ec event.Composite) error {
	logrus.Infof("received event with payload %v", ec)

	var adto dto.Article
	if err := mapstructure.Decode(ec.Event.Data, &adto); err != nil {
		return fmt.Errorf("error converting event data into Article with error: %w", err)
	}

	if err := r.convertAndSave(adto); err != nil {
		return err
	}

	return nil
}

func (r *articleResolver) convertAndSave(adto dto.Article) error {
	abght := adto.ToArticleByGlobalHashTag()

	for _, a := range abght {
		if err := r.ar.SaveArticleByGlobalTags(a); err != nil {
			// TODO approach error handling better (by appending errors maybe)
			logrus.Errorf("error when saving article by global hash tags with values %v and error %v", a, err)
		}
	}

	uba := adto.ToUserByArticle()
	if err := r.ar.SaveUserByArticle(uba); err != nil {
		// TODO approach error handling better (by appending errors maybe)
		logrus.Errorf("error when saving user by article with values %v and error %v", uba, err)
	}

	return nil
}
