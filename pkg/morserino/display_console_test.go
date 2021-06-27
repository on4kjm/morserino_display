package morserino

/*
Copyright © 2021 Jean-Marc Meessen, ON4KJM <on4kjm@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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
