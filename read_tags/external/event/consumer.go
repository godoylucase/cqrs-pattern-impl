package event

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	// TODO change the ports to be internal when dockerized app is ready
	brokers = []string{"kafka:9092"}
)

type Consumer struct {
	c    sarama.Consumer
	done <-chan interface{}
}

func NewConsumer(done <-chan interface{}) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("error creating a new cons, with error %w", err)
	}

	return newConsumer(done, consumer)
}

// newConsumer valid constructor to inject mocks for the consumer during unit testing
func newConsumer(done <-chan interface{}, consumer sarama.Consumer) (*Consumer, error) {
	return &Consumer{c: consumer, done: done}, nil
}

func (c *Consumer) Get(topic string) (<-chan Composite, error) {
	eventComposites := make(chan Composite)
	//defer close(eventComposites)

	partitions, err := c.c.Partitions(topic)
	if err != nil {
		return nil, fmt.Errorf("error getting partitions with error: %w", err)
	}

	for p := range partitions {
		pc, err := c.c.ConsumePartition(topic, int32(p), sarama.OffsetOldest)
		if err != nil {
			return nil, fmt.Errorf("error consuming kafka partition %v, with error: %w", p, err)
		}

		go func(pc sarama.PartitionConsumer) {
			<-c.done
			pc.AsyncClose()
		}(pc)

		go func(done <-chan interface{},
			topic string, pc sarama.PartitionConsumer) {
			defer close(eventComposites)

			for {
				select {
				case msg, ok := <-pc.Messages():
					if !ok {
						return
					}
					val := msg.Value

					var event Event
					if err := json.Unmarshal(val, &event); err != nil {
						eventComposites <- Composite{
							Err: fmt.Errorf("error unmarshalling event from '%v' topic, with error: %w", topic, err),
						}
					}

					eventComposites <- Composite{
						Event: event,
						Err:   nil,
					}
				case err, ok := <-pc.Errors():
					if !ok {
						return
					}

					eventComposites <- Composite{
						Err: fmt.Errorf("error unmarshalling event from '%v' topic, with error: %w", topic, err),
					}
				}
			}
		}(c.done, topic, pc)
	}

	return eventComposites, nil
}
