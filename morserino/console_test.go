package morserino

//TODO: create proper documentation and attribution (https://pkg.go.dev/testing/iotest)

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestConsole(t *testing.T){


	dev, err := Open("simulator")
	if err != nil {
		t.Errorf("Unable to open morserino device. err=%v", err)
	}

	defer dev.Close()

	ctx := context.Background()

	var (
		// ctx context.Context

		// An errgroup is an abstraction that allows to wait for the completion of a group of goroutines.
		// see https://pkg.go.dev/golang.org/x/sync/errgroup.
		gr, gctx = errgroup.WithContext(ctx)
		console  = NewConsole()
		listener = NewListener(dev, console)
	)

	// Let'snow start our listener goroutine, that will emit events as soon as it gets a meaningful message.
	gr.Go(func() error {
		return listener.Listen(gctx)
	})

	// To finish, let's start our console goroutine that will handle events coming from the user and from the listener.
	gr.Go(func() error {
		return console.Run(gctx)
	})

	// Let's now wait for completion of all of those goroutines, and complain in case of error.
	if err := gr.Wait(); err != nil && !errors.Is(err, ErrEmitterDone) {
		t.Errorf("Received an unexpected error, exiting. err=%v", err)
	}

}
