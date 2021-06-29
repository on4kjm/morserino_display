package morserino

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/rs/zerolog/log"
	"go.bug.st/serial"
)

const defaultReadSize = 100

var (
	ErrListenDone = errors.New("listener exited")
)

// Listener reads on the given device, and pushes events to the given event Handler as soon as needed.
// It raises ErrListenDone as soon as it terminates.
type Listener struct {
	com io.ReadCloser

	hdl EventHandler
}

func NewListener(com io.ReadCloser, handler EventHandler) *Listener {
	return &Listener{com: com, hdl: handler}
}

func (c *Listener) Listen(ctx context.Context) error {
	buf := make([]byte, defaultReadSize)

	// A first goroutine is being used to track for context cancellation.
	// A context carries the deadline of an intent. In other words, as soon as a context is not valid
	// anymore, then any goroutine depending on that context needs to exit.
	// More about context => https://www.youtube.com/watch?v=LSzR0VEraWw
	go func() {
		// Wait and block until context is cancelled.
		<-ctx.Done()

		// As soon as the context is expired, we'll close the device.
		// When the device is closed, any read // write will return a specific error,
		// that could be handled by the caller code.
		//
		// We'll use this to signal termination for our read loop.
		if err := c.com.Close(); err != nil {
			log.Err(err).Msg("Can't close the com device")
		}
	}()

	// Read indefinitely until we face an error...
	for {
		// Reads up to 100 bytes.
		n, err := c.com.Read(buf)

		// Let's check the **type** of our error.
		switch perr := err.(type) {
		case *serial.PortError:
			// This error is an error coming from the serial package.

			if perr.Code() == serial.PortClosed {
				// Someone closed the port, so we need to exit.
				// Though, that's not an error.
				// It's part of our termination mechanism as stated above.
				// Let's signal termination by returning ErrListenDone.
				return ErrListenDone
			}

			// Returns this error if it's not a PortClosed error.
			return err
		case nil:
			// No error, let's continue execution.
		default:
			// Let's not complain about EOF. It's means that our job is done but not an error.
			if errors.Is(err, io.EOF) {
				return ErrListenDone
			}

			// Some unknown error happened, let's exit.
			return err
		}

		if n == 0 {
			log.Debug().Msg("Detected an empty read, EOF")

			if err := c.hdl.Handle(Event{Kind: KindEOF}); err != nil {
				log.Err(err).Msg("Failed to notify EOF")
			}

			return ErrListenDone
		}

		if bytes.HasPrefix(buf, []byte("<sk> e e")) {
			log.Trace().
				Msg("Sending exit marker to the console displayer as we received the exit sequence")

			if err := c.hdl.Handle(Event{Kind: KindEOF}); err != nil {
				log.Err(err).Msg("Failed to notify EOF")
				return err
			}

			log.Trace().Msg("Exiting listen loop")

			return ErrListenDone
		}

		if err := c.hdl.Handle(Event{Kind: KindMessage, Payload: buf[:n]}); err != nil {
			log.Err(err).Msg("Unable to Handle a message")
			return err
		}
	}
}
