package morserino

import (
	"context"
	"io"
)

type FormattedConsole struct {
	eventQueue chan Event
	output     io.Writer
}

func NewFormattedConsole(output io.Writer) *FormattedConsole {
	return &FormattedConsole{
		eventQueue: make(chan Event, eventQueueSize),
		output:     output,
	}
}

func (c *FormattedConsole) Handle(evt Event) error {
	c.eventQueue <- evt
	return nil
}

func (c *FormattedConsole) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case evt, ok := <-c.eventQueue:
			if !ok {
				return nil
			}

			if evt.Kind == KindEOF {
				return nil
			}

			// Event handling logic
			c.output.Write(evt.Payload)
			// c.output.Write(evt.
		}
	}
}
