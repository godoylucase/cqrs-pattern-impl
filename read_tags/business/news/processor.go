package news

import (
	"fmt"
	"sync"

	"github.com/godoylucase/read_tags/external/event"
	"github.com/sirupsen/logrus"
)

var (
	resourceTopics = map[event.Resource]string{
		event.ARTICLE: "articles",
	}
)

type Resolver interface {
	Run(eventComposite event.Composite) error
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
		return fmt.Errorf("consumer error with value: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func(wg *sync.WaitGroup, ecs <-chan event.Composite) {
		defer wg.Done()
		for ec := range ecs {
			if err := p.resolver.Run(ec); err != nil {
				logrus.Errorf("error when processing incoming events from topic with error: %v", err)
			}
		}
	}(&wg, received)
	wg.Wait()

	return nil
}
