package event

type Action uint
type Resource uint

const (
	CREATE Action = iota + 100
	UPDATE
	DELETE
)

const (
	ARTICLE Resource = iota * 100
	TAG
)

func (a Action) String() string {
	m := map[Action]string{
		CREATE: "create",
		UPDATE: "update",
		DELETE: "delete",
	}

	action, ok := m[a]
	if !ok {
		panic("action not supported")
	}

	return action
}

func (a Resource) String() string {
	m := map[Resource]string{
		ARTICLE: "article",
		TAG:     "tag",
	}

	res, ok := m[a]
	if !ok {
		panic("resource not supported")
	}

	return res
}

type Partitionable interface {
	GetKey() string
}

type Event struct {
	Key          string      `json:"key"`
	FromResource string      `json:"resource"`
	WithAction   string      `json:"action"`
	Data         interface{} `json:"data"`
}

func New(r Resource, a Action, data Partitionable) *Event {
	return &Event{
		Key:          data.GetKey(),
		FromResource: r.String(),
		WithAction:   a.String(),
		Data:         data,
	}
}

type Composite struct {
	Event Event
	Err   error
}
