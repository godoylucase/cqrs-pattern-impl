package business

import "github.com/godoylucase/articles_tags/internal/event"

type eb struct {
	dispatcher *event.Dispatcher
}

func NewEventBroker(done <-chan interface{}) *eb {
	d := event.NewDispatcher(done)
	return &eb{dispatcher: d}
}

func (eb *eb) Close() {
	eb.dispatcher.Close()
}

func (eb *eb) ArticleCreation(article event.Partitionable) error {
	if err := eb.dispatcher.Send(event.ARTICLE, event.CREATE, article); err != nil {
		return err
	}

	return nil
}

func (eb *eb) ArticleUpdate(article event.Partitionable) error {
	if err := eb.dispatcher.Send(event.ARTICLE, event.UPDATE, article); err != nil {
		return err
	}

	return nil
}
