package event

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

var (
	// TODO change the ports to be internal when dockerized app is ready
	brokers = []string{"localhost:29092"}
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

	return &Consumer{c: consumer, done: done}, nil
}

func (c *Consumer) Get(topic string, received chan Composite) error {
	defer close(received)

	ps, err := c.c.Partitions(topic)
	if err != nil {
		return fmt.Errorf("error getting partitions with error: %w", err)
	}

	for p := range ps {
		consumer, err := c.c.ConsumePartition(topic, int32(p), sarama.OffsetOldest)
		if err != nil {
			return fmt.Errorf("error consuming kafka partition %v, with error: %w", p, err)
		}

		go func(done <-chan interface{},
			topic string, pc sarama.PartitionConsumer) {

			for {
				select {
				case msg := <-consumer.Messages():
					val := msg.Value

					var event Event
					if err := json.Unmarshal(val, &event); err != nil {
						received <- Composite{
							Err: fmt.Errorf("error unmarshalling event from '%v' topic, with error: %w", topic, err),
						}
					}

					received <- Composite{
						Event: event,
						Err:   nil,
					}
				case err := <-consumer.Errors():
					received <- Composite{
						Err: fmt.Errorf("error unmarshalling event from '%v' topic, with error: %w", topic, err),
					}
				case <-done:
					return
				}
			}
		}(c.done, topic, consumer)

	}

	return nil
}
