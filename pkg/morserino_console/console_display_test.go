package morserino_console

import (
	"bufio"
	// "log"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/on4kjm/morserino_display/pkg/morserino_channels"
	"github.com/on4kjm/morserino_display/pkg/morserino_com"
	"github.com/stretchr/testify/assert"
)

func TestConsoleDisplayListener_happyCase(t *testing.T) {
	// ** Given
	testMsg := "Test = test <skaaa <sk> e e"
	mc := &morserino_channels.MorserinoChannels{}
	mc.Init()

	f, err := os.Create("testfile.txt")
	assert.NoError(t, err)

	defer f.Close()

	w := bufio.NewWriter(f)

	// ** When
	go serialListenerMock(testMsg, mc)
	go ConsoleDisplayListener(mc, w)

	<-mc.Done

	// ** Then
	// fmt.Println(buff.String())
	w.Flush()
	fi, err := f.Stat()
	assert.NoError(t, err)
	var zero int64
	zero = 0
	assert.Greater(t, fi.Size(), zero)

}

// A mock to simulaate the serial port listener goroutine
func serialListenerMock(TestString string, mc *morserino_channels.MorserinoChannels) {
	morserino_com.Listen(iotest.OneByteReader(strings.NewReader(TestString)), mc)
	return
}
