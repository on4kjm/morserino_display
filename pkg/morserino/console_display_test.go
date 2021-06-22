package morserino

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"testing/iotest"
	// "time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05.000"}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}

	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	AppLogger = zerolog.New(output).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.TraceLevel)

}

func TestConsoleDisplayListener_happyCase(t *testing.T) {

	// ** Given
	testMsg := "Test = test <skaaa <sk> e e"
	mc := &MorserinoChannels{}
	mc.Init()

	//TODO: Create a temporary file in the temporary directory
	testFileName := "testfile.txt"
	f, err := os.Create(testFileName)
	assert.NoError(t, err)

	defer f.Close()

	w := bufio.NewWriter(f)
	AppLogger.Debug().Msg("Starting test, output to \"" + testFileName + "\"")

	// ** When
	go serialListenerMock(testMsg, mc)
	go ConsoleDisplayListener(mc, w)

	AppLogger.Debug().Msg("Waiting on Done channel")
	<-mc.Done //Waiting here for everything to be orderly completed
	AppLogger.Debug().Msg("All is completed. Go Routines exited")

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
func serialListenerMock(TestString string, mc *MorserinoChannels) {
	Listen(iotest.OneByteReader(strings.NewReader(TestString)), mc)
	return
}
