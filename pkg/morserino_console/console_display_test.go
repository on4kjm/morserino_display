package morserino_console

import (
	"strings"
	"testing"
	"testing/iotest"

	"github.com/on4kjm/morserino_display/pkg/morserino_channels"
	"github.com/on4kjm/morserino_display/pkg/morserino_com"
)

func TestConsoleDisplayListener_happyCase(t *testing.T) {
	// ** Given
	testMsg := "Test = test <skaaa <sk> e e"
	var mc morserino_channels.MorserinoChannels
	mc.Init()

	// ** When
	go serialListenerMock(testMsg, mc)
	go ConsoleDisplayListener(mc)

	// ** Then
}

// A mock to simulaate the serial port listener goroutine
func serialListenerMock(TestString string, mc morserino_channels.MorserinoChannels) {
	morserino_com.Listen(iotest.OneByteReader(strings.NewReader(TestString)), mc)
}
