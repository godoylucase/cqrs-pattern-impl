package news

import (
	"github.com/godoylucase/read_tags/business"
	"github.com/godoylucase/read_tags/internal/event"
)

type articleRepository interface {
	SaveArticleByGlobalTags(dto business.ArticleDTO) error
}

type articleReceiver struct {
	ar       articleRepository
	received <-chan event.Composite
}

func NewArticleReceiver(ar articleRepository) *articleReceiver {
	return &articleReceiver{ar: ar}
}

func (r *articleReceiver) Run(received chan event.Composite) error {

	return nil
}
