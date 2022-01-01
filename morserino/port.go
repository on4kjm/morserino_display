package morserino

import (
	"io"
	"strings"

	"go.bug.st/serial"
)

var (
	defaultMode = serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
)

const simulatorMessage = "cq cq de on4kjm on4kjm = tks fer call om = ur rst 599 = hw? \n73 de on4kjm = <sk> e e"

func Open(portName string) (io.ReadWriteCloser, error) {
	switch strings.ToLower(portName) {
	case "simulator":
		return nopCloseWriter{DelayedOneByteReader(strings.NewReader(simulatorMessage))}, nil
	default:
		return serial.Open(portName, &defaultMode)
	}
}

type nopCloseWriter struct {
	io.Reader
}

func (nopCloseWriter) Write([]byte) (int, error) { return 0, nil }
func (nopCloseWriter) Close() error              { return nil }
