package news

import (
	"fmt"

	"github.com/godoylucase/read_tags/internal/event"
	"github.com/sirupsen/logrus"
)

var (
	resourceTopics = map[event.Resource]string{
		event.ARTICLE: "articles",
	}
)

type Resolver interface {
	Run(received chan event.Composite) error
}

type consumer interface {
	Get(topic string, received chan event.Composite) error
}

type processor struct {
	resolver Resolver
	received chan event.Composite
	resource event.Resource
	consumer consumer
}

func New(resource event.Resource, consumer consumer, resolver Resolver, received chan event.Composite) *processor {
	return &processor{
		resolver: resolver,
		received: received,
		resource: resource,
		consumer: consumer,
	}
}

func (p *processor) Run() {
	topic, ok := resourceTopics[p.resource]
	if !ok {
		panic(fmt.Errorf("the %v is a not supported topic", topic))
	}

	if err := p.consumer.Get(topic, p.received); err != nil {
		panic(err)
	}

	if err := p.resolver.Run(p.received); err != nil {
		logrus.Errorf("error consuming events from topic %v with error %v", topic, err)
	}
}
