package morserino

//TODO: create proper documentation and attribution (https://pkg.go.dev/testing/iotest)

import (
	"io"
	"time"
)

// DelayedOneByteReader returns a Reader that implements
// each non-empty Read by reading one byte from r.
func DelayedOneByteReader(r io.Reader) io.Reader { return &delayedOneByteReader{r} }

type delayedOneByteReader struct {
	r io.Reader
}

func (r *delayedOneByteReader) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	time.Sleep(1 * time.Second)
	return r.r.Read(p[0:1])
}