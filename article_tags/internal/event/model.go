package event

type Action uint
type Resource uint
type Key string

const (
	CREATE Action = iota + 100
	UPDATE
	DELETE
)

const (
	ARTICLE Resource = iota * 100
	TAG
)

type Partitionable interface {
	GetKey() string
}

type Event struct {
	Key          Key         `json:"key"`
	FromResource Resource    `json:"resource"`
	WithAction   Action      `json:"action"`
	Data         interface{} `json:"data"`
}

func New(r Resource, a Action, data Partitionable) *Event {
	return &Event{
		Key:          Key(data.GetKey()),
		FromResource: r,
		WithAction:   a,
		Data:         data,
	}
}
