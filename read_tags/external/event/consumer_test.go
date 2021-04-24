package event

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	"github.com/stretchr/testify/assert"
)

func TestConsumer_Get(t *testing.T) {
	fmt.Println(runtime.NumGoroutine())

	consMock := mocks.NewConsumer(t, mocks.NewTestConfig())

	consMock.SetTopicMetadata(map[string][]int32{
		"article": {0},
	})

	createdArticleMock := getNewEvent(ARTICLE, CREATE, 1)
	updatedArticleMock := getNewEvent(ARTICLE, UPDATE, 2)

	consMock.ExpectConsumePartition("article", 0, sarama.OffsetOldest).
		YieldMessage(&sarama.ConsumerMessage{
			Value: getJSONBytes(createdArticleMock)},
		)

	consMock.ExpectConsumePartition("article", 0, sarama.OffsetOldest).
		YieldError(sarama.ErrOutOfBrokers)

	consMock.ExpectConsumePartition("article", 0, sarama.OffsetOldest).
		YieldMessage(&sarama.ConsumerMessage{
			Value: getJSONBytes(updatedArticleMock)},
		)

	done := make(chan interface{})
	defer close(done)

	consumerMock, err := newConsumer(done, consMock)
	if err != nil {
		t.Fatalf("unexpected error")
	}

	composites, err := consumerMock.Get("article")
	if err != nil {
		t.Fatalf("unexpected error")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup, cs <-chan Composite) {
		expected := getNewComposite(createdArticleMock)
		c1 := <-cs
		assert.EqualValues(t, expected, c1)

		expectedErr := sarama.ErrOutOfBrokers
		c2 := <-cs
		assert.Equal(t, expectedErr, c2.Err)

		expected = getNewComposite(updatedArticleMock)
		c3 := <-cs
		assert.EqualValues(t, expected, c3)

		wg.Done()
	}(&wg, composites)
}

func getNewEvent(r Resource, a Action, value int) Event {
	return Event{
		Key:          "someValidKey",
		FromResource: r.String(),
		WithAction:   a.String(),
		Data:         map[string]interface{}{"value": value},
	}
}

func getNewComposite(event Event) Composite {
	return Composite{
		Event: event,
		Err:   nil,
	}
}

func getJSONBytes(value interface{}) []byte {
	bytes, _ := json.Marshal(value)
	return bytes
}
