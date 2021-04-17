package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	// TODO change the ports to be internal when dockerized app is ready
	brokers      = []string{"localhost:29092"}
	actionTopics = map[Resource]string{
		ARTICLE: "articles",
	}
)

type Dispatcher struct {
	p    sarama.AsyncProducer
	done <-chan interface{}
}

func NewDispatcher(done <-chan interface{}) *Dispatcher {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil
	}

	return &Dispatcher{p: producer, done: done}
}

func (d *Dispatcher) Send(r Resource, a Action, p Partitionable) error {
	topic, ok := actionTopics[r]
	if !ok {
		return fmt.Errorf("topic not supported: %v", r)
	}

	event := New(r, a, p)
	if err := d.doSend(event, topic); err != nil {
		return fmt.Errorf("error when sending event %v, with error: %w", event, err)
	}

	for {
		select {
		case s := <-d.p.Successes():
			logrus.Infof("event sent to topic %v with payload: %v", s.Topic, p)
			return nil
		case err := <-d.p.Errors():
			logrus.Errorf("event could not be sent to topic %v with err: %v", topic, err)
			return err
		case <-d.done:
			return nil
		}
	}
}

func (d *Dispatcher) doSend(e *Event, topic string) error {
	encodedMsg, err := json.Marshal(e)
	if err != nil {
		return err
	}

	d.p.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(e.Key),
		Value:     sarama.ByteEncoder(encodedMsg),
		Timestamp: time.Time{},
	}

	return nil
}
