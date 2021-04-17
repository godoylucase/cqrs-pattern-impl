package business

import "github.com/godoylucase/cqrs-pattern-impl/internal/event"

type eb struct {
	dispatcher *event.Dispatcher
}

func NewEventBroker(done <-chan interface{}) *eb {
	d := event.NewDispatcher(done)
	return &eb{dispatcher: d}
}

func (eb *eb) ArticleCreation(article event.Partitionable) error {
	if err := eb.dispatcher.Send(event.ARTICLE, event.CREATE, article); err != nil {
		return err
	}

	return nil
}
