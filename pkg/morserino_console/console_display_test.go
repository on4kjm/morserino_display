package morserino_console

import (
	"bufio"
	"io/ioutil"
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

	testFileName := "testfile.txt"
	f, err := os.Create(testFileName)
	assert.NoError(t, err)

	defer f.Close()

	w := bufio.NewWriter(f)

	// f := bufio.NewWriter(os.Stdout)

	// ** When
	go serialListenerMock(testMsg, mc)
	go ConsoleDisplayListener(mc, w)

	<-mc.Done //Waiting here for everything to be orderly completed

	w.Flush() //Just to be sure everything is written to disk

	// ** Then
	b, err := ioutil.ReadFile(testFileName)
	assert.NoError(t, err)

	expectedOutput := "Test =\n test <skaaa <sk> e e\nExiting...\n"
	assert.Equal(t, expectedOutput, string(b))

	deleteErr := os.Remove(testFileName)
	assert.NoError(t, deleteErr)
}

// A mock to simulaate the serial port listener goroutine
func serialListenerMock(TestString string, mc *morserino_channels.MorserinoChannels) {
	morserino_com.Listen(iotest.OneByteReader(strings.NewReader(TestString)), mc)
	return
}
