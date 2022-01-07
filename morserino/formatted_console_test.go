package morserino

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestFormattedConsole(t *testing.T) {
	testCases := []struct {
		desc       string
		events     []Event
		wantOutput string
	}{
		{
			desc: "writes to output",
			events: []Event{
				{
					Kind:    KindMessage,
					Payload: []byte{'c'},
				},
				{Kind: KindEOF},
			},
			wantOutput: "c",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			var (
				output bytes.Buffer
				group  errgroup.Group
				err    error
			)

			// Initialize our console.
			subject := NewFormattedConsole(&output)

			// Run a goroutine to begin processing events while we're sending them.
			// This is done to avoid blocking the eventQueue.
			group.Go(func() error {
				return subject.Run(context.Background())
			})

			// Send our events from the test case.
			for _, event := range test.events {
				err = subject.Handle(event)
				require.NoError(t, err)
			}

			// Wait for our formatted console to process all the sent events.
			err = group.Wait()
			require.NoError(t, err)

			// Assert that the output is what we actually want"
			assert.Equal(t, test.wantOutput, string(output.Bytes()))
		})
	}
}
