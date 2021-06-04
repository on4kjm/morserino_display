package morserino_com

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListen(t *testing.T) {
	f, err := os.Open("./example-file.txt")
	require.NoError(t, err)

	defer f.Close()

	err = Listen(f)
	assert.NoError(t, err)
}
