package morserino

import (
	"context"

	"github.com/rs/zerolog/log"
)

const eventQueueSize = 10

// Console implements EventHandler and prints received events to the console.
type Console struct {
	eventQueue chan Event
}

func NewConsole() *Console {
	return &Console{
		eventQueue: make(chan Event, eventQueueSize),
	}
}

// Handle pushes an event into the events queue so it is processed by the console event loop.
// This call needs to be non blocking as this method is called from the listener loop.
func (c *Console) Handle(evt Event) error {
	c.eventQueue <- evt
	return nil
}

// Run executes the main console event loop.
// - It watches for context cancellation.
// - It processes the events from the eventQueue and prints their content to the console screen.
func (c *Console) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case evt, ok := <-c.eventQueue:
			if !ok {
				return nil
			}

			log.Info().Int("Kind", int(evt.Kind)).Msg(string(evt.Payload))
		}
	}
}
