package morserino

// A Kind qualifies what an Event actualy represents.
type Kind int

const (
	KindUnknown Kind = iota
	KindEOF
	KindMessage
)

// Event represents something meaningful read from the underlying device.
type Event struct {
	Kind    Kind
	Payload []byte
}

// EventHandler represents any type that can handle an event.
type EventHandler interface {
	Handle(evt Event) error
}
