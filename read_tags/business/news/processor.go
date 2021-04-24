package news

import (
	"fmt"

	"github.com/godoylucase/read_tags/external/event"
)

var (
	resourceTopics = map[event.Resource]string{
		event.ARTICLE: "articles",
	}
)

type Resolver interface {
	Run(done <-chan interface{}, received <-chan event.Composite)
}

type consumer interface {
	Get(topic string) (<-chan event.Composite, error)
}

type processor struct {
	done     <-chan interface{}
	resource event.Resource
	consumer consumer
	resolver Resolver
}

func NewProcessor(done <-chan interface{}, resource event.Resource, resolver Resolver, consumer consumer) *processor {
	return &processor{
		done:     done,
		resource: resource,
		consumer: consumer,
		resolver: resolver,
	}
}

func (p *processor) Run() error {
	topic, ok := resourceTopics[p.resource]
	if !ok {
		return fmt.Errorf("the %v is a not supported topic", topic)
	}

	received, err := p.consumer.Get(topic)
	if err != nil {
		return err
	}

	p.resolver.Run(nil, received)

	return nil
}
