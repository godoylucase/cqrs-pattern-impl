package news

import (
	"fmt"

	"github.com/godoylucase/read_tags/external/event"
)

type ReceiverFactory struct {
	ar articleRepository
}

func NewResolverFactory(ar articleRepository) *ReceiverFactory {
	return &ReceiverFactory{ar: ar}
}

func (rf *ReceiverFactory) Get(r event.Resource) (Resolver, error) {
	switch r {
	case event.ARTICLE:
		return NewArticleResolver(rf.ar), nil
	default:
		return nil, fmt.Errorf("there is no resolver available for the resource %v", r)
	}
}
