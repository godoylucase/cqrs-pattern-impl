package news

import (
	"strconv"
	"testing"

	"github.com/godoylucase/read_tags/external/event"
	"github.com/stretchr/testify/assert"
)

type resolverMock struct {
	runMockFn func(eventComposite event.Composite)
}

func (r *resolverMock) Run(eventComposite event.Composite) {
	r.runMockFn(eventComposite)
}

type consumerMock struct {
	getMockFn func(topic string) (<-chan event.Composite, error)
}

func (c *consumerMock) Get(topic string) (<-chan event.Composite, error) {
	return c.getMockFn(topic)
}

func Test_processor_Run(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	eventQty := 5
	slide := getEventSlide(eventQty)

	rm := &resolverMock{
		func(ec event.Composite) {
			assert.EqualValues(t, slide[ec.Event.Data.(int)], ec)
		}}

	cm := &consumerMock{
		func(topic string) (<-chan event.Composite, error) {
			received := make(chan event.Composite)

			go func() {
				defer close(received)
				for i := 0; i < eventQty; i++ {
					received <- slide[i]
				}
			}()

			return received, nil
		}}

	assert.Nil(t, NewProcessor(done, event.ARTICLE, rm, cm).Run())
}

func getEventSlide(eventQty int) []event.Composite {
	var ecs []event.Composite
	for i := 0; i < eventQty; i++ {
		ecs = append(ecs, event.Composite{
			Event: event.Event{
				Key:          strconv.Itoa(i),
				FromResource: event.ARTICLE.String(),
				WithAction:   event.CREATE.String(),
				Data:         i,
			},
		})
	}
	return ecs
}
