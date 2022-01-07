package morserino

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormattedConsole(t *testing.T) {
	var (
		output bytes.Buffer

		ctx         = context.Background()
		inputEvents = []Event{
			{
				Kind:    KindMessage,
				Payload: []byte{'c'},
			},
			{Kind: KindEOF},
		}

		subject = NewFormattedConsole(&output)
	)

	for _, event := range inputEvents {
		err := subject.Handle(event)
		require.NoError(t, err)
	}

	err := subject.Run(ctx)
	require.NoError(t, err)

	assert.Equal(t, "c", string(output.Bytes()))
}
