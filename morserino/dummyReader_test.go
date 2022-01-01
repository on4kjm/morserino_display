package morserino

//TODO: create proper documentation and attribution (https://pkg.go.dev/testing/iotest)

import (
	"bytes"
	"io"
	"testing"
)

func TestOneByteReader_nonEmptyReader(t *testing.T) {
	msg := "Hello, World!"
	buf := new(bytes.Buffer)
	buf.WriteString(msg)

	obr := DelayedOneByteReader(buf)
	var b []byte
	n, err := obr.Read(b)
	if err != nil || n != 0 {
		t.Errorf("Empty buffer read returned n=%d err=%v", n, err)
	}

	b = make([]byte, 3)
	// Read from obr until EOF.
	got := new(bytes.Buffer)
	for i := 0; ; i++ {
		n, err = obr.Read(b)
		if err != nil {
			break
		}
		if g, w := n, 1; g != w {
			t.Errorf("Iteration #%d read %d bytes, want %d", i, g, w)
		}
		got.Write(b[:n])
	}

	//TODO: add a check that the delay wasa correct
	if g, w := err, io.EOF; g != w {
		t.Errorf("Unexpected error after reading all bytes\n\tGot:  %v\n\tWant: %v", g, w)
	}
	if g, w := got.String(), "Hello, World!"; g != w {
		t.Errorf("Read mismatch\n\tGot:  %q\n\tWant: %q", g, w)
	}
}

func TestOneByteReader_emptyReader(t *testing.T) {
	r := new(bytes.Buffer)

	obr := DelayedOneByteReader(r)
	var b []byte
	if n, err := obr.Read(b); err != nil || n != 0 {
		t.Errorf("Empty buffer read returned n=%d err=%v", n, err)
	}

	b = make([]byte, 5)
	n, err := obr.Read(b)
	if g, w := err, io.EOF; g != w {
		t.Errorf("Error mismatch\n\tGot:  %v\n\tWant: %v", g, w)
	}
	if g, w := n, 0; g != w {
		t.Errorf("Unexpectedly read %d bytes, wanted %d", g, w)
	}
}
