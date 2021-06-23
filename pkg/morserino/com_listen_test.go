package morserino

import (
	"fmt"
	// "io/ioutil"
	"os"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/on4kjm/morserino_display/pkg/safebuffer"
	"github.com/rs/zerolog"
	// "github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bug.st/serial/enumerator"
)

func init() {

	zerolog.TimeFieldFormat = "15:04:05.000"
	AppLogger = zerolog.New(os.Stdout).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.TraceLevel)

}

//
// Testing DetectDevice()
//
func TestDetectDevice_HappyCase(t *testing.T) {
	//Given
	mockPortEnum := mockPortEnumerator{ports: []*enumerator.PortDetails{
		{Name: "COM Port Name", IsUSB: true, VID: "XXXX", PID: "YYYY", SerialNumber: "001"},
		{Name: "Morserino Port", IsUSB: true, VID: "10C4", PID: "EA60", SerialNumber: "002"}},
	}

	//when
	got, err := DetectDevice(mockPortEnum)

	//Then
	require.NoError(t, err)
	assert.Equal(t, "Morserino Port", got)
}

func TestDetectDevice_NoMoresrinoDetected(t *testing.T) {
	//Given
	mockPortEnum := mockPortEnumerator{ports: []*enumerator.PortDetails{
		{Name: "COM Port Name", IsUSB: true, VID: "XXXX", PID: "YYYY", SerialNumber: "001"}},
	}

	//when
	got, err := DetectDevice(mockPortEnum)

	//Then
	require.EqualError(t, err, "Did not find a usable port.")
	assert.Equal(t, "", got)
}

func TestDetectDevice_TooManyMorserinosDetected(t *testing.T) {
	//Given
	mockPortEnum := mockPortEnumerator{ports: []*enumerator.PortDetails{
		{Name: "COM Port Name", IsUSB: true, VID: "XXXX", PID: "YYYY", SerialNumber: "001"},
		{Name: "Morserino Port", IsUSB: true, VID: "10C4", PID: "EA60", SerialNumber: "002"},
		{Name: "Morserino Port2", IsUSB: true, VID: "10C4", PID: "EA60", SerialNumber: "003"}},
	}

	//when
	got, err := DetectDevice(mockPortEnum)

	//Then
	require.EqualError(t, err, "ERROR: Multiple Morserino devices found.")
	assert.Equal(t, "", got)
}

func TestDetectDevice_PortEnumerationWentWrong(t *testing.T) {
	//Given
	mockPortEnum := mockPortEnumerator{ports: nil, err: fmt.Errorf("An error occured")}

	//when
	got, err := DetectDevice(mockPortEnum)

	//Then
	require.EqualError(t, err, "An error occured")
	assert.Equal(t, "", got)
}

//
// Test Listen()
//

func TestListen_HappyCase(t *testing.T) {
	AppLogger.Info().Msg("==> " + t.Name())

	// Given
	testMsg := "Test = test <sk> e e"
	mock := iotest.OneByteReader(strings.NewReader(testMsg))
	var testBuffer safebuffer.Buffer
	mc := &MorserinoChannels{}
	mc.Init()

	// When
	// Starts listener function so that we can check what has been actually received
	AppLogger.Debug().Msg("Starting the mock listener in the background")
	go MockMorserinoDisplayer(mc, &testBuffer)
	AppLogger.Debug().Msg("Starting the Listen under test (with the mock port reader)")
	//FIXME: How to receive errors
	// err := Listen(mock, mc)
	go Listen(mock, mc)

	AppLogger.Debug().Msg("Waiting for the done signal")
	<-mc.Done
	AppLogger.Debug().Msg("All is completed. Go Routines exited")

	// Then
	// require.NoError(t, err)
	assert.Equal(t, testBuffer.String(), testMsg+exitString)

	AppLogger.Info().Msg("<== " + t.Name())
}

func TestListen_missedEndMarker(t *testing.T) {
	AppLogger.Info().Msg("==> " + t.Name())

	// Given
	testMsg := "Test = test <skaaa <sk> e e"
	mock := iotest.OneByteReader(strings.NewReader(testMsg))
	var testBuffer safebuffer.Buffer
	mc := &MorserinoChannels{}
	mc.Init()

	// When
	// Starts listener function so that we can check what has been actually received
	AppLogger.Debug().Msg("Starting the mock listener in the background")
	go MockMorserinoDisplayer(mc, &testBuffer)
	AppLogger.Debug().Msg("Starting the Listen under test (with the mock port reader)")
	go Listen(mock, mc)

	AppLogger.Debug().Msg("Waiting for the done signal")
	<-mc.Done
	AppLogger.Debug().Msg("All is completed. Go Routines exited")

	// Then
	assert.Equal(t, testMsg+exitString, testBuffer.String())

	AppLogger.Info().Msg("<== " + t.Name())
}

//EOF error (no error but no data returned)
func TestListen_EOF(t *testing.T) {
	AppLogger.Info().Msg("==> " + t.Name())

	// Given
	testMsg := "\nEOF"
	mock := iotest.ErrReader(nil)
	var testBuffer safebuffer.Buffer
	mc := &MorserinoChannels{}
	mc.Init()

	// When
	AppLogger.Debug().Msg("Starting the mock listener in the background")
	go MockMorserinoDisplayer(mc, &testBuffer)
	AppLogger.Debug().Msg("Starting the Listen under test (with the mock port reader)")
	go Listen(mock, mc)

	AppLogger.Debug().Msg("Waiting for the done signal")
	<-mc.Done
	AppLogger.Debug().Msg("All is completed. Go Routines exited")

	// Then
	assert.Equal(t, testMsg+exitString, testBuffer.String())

	AppLogger.Info().Msg("<== " + t.Name())
}

//
// Test Listen()
//
func TestListen_withSimulator(t *testing.T) {
	AppLogger.Info().Msg("==> " + t.Name())

	// Given
	testMsg := "cq cq de on4kjm on4kjm = tks fer call om = ur rst 599 = hw? \n73 de on4kjm = <sk> e e"
	var testBuffer safebuffer.Buffer
	mc := &MorserinoChannels{}
	mc.Init()

	// When
	// Starts listener function so that we can check what has been actually received
	AppLogger.Debug().Msg("Starting the Listen under test (with the mock port reader)")
	go MockMorserinoDisplayer(mc, &testBuffer)
	AppLogger.Debug().Msg("Starting the OpenAndListen() under test (with simulator)")
	go OpenAndListen("simul", nil, mc)

	AppLogger.Debug().Msg("Waiting for the done signal")
	<-mc.Done
	AppLogger.Debug().Msg("All is completed. Go Routines exited")

	// Then
	assert.Equal(t, testBuffer.String(), testMsg+exitString)

	AppLogger.Info().Msg("<== " + t.Name())
}

func MockMorserinoDisplayer(mc *MorserinoChannels, workBuffer *safebuffer.Buffer) {

	for {
		var output string
		output = <-mc.MessageBuffer
		_, err := workBuffer.Write([]byte(output))
		if err != nil {
			AppLogger.Error().Err(err).Msg("Error writing to safebuffer")
		}

		if strings.Contains(output, "\nExiting...\n") {
			AppLogger.Trace().Msg("MockMorserinoDisplayer: Signaling that the display processing is complete")
			mc.DisplayCompleted <- true
			return
		}
	}
}
