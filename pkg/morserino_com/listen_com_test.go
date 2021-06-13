package morserino_com

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/on4kjm/morserino_display/pkg/safebuffer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bug.st/serial/enumerator"
)

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
	// Given
	testMsg := "Test = test <sk> e e"
	mock := iotest.OneByteReader(strings.NewReader(testMsg))
	var MessageBuffer = make(chan string, 10)
	var testBuffer safebuffer.Buffer

	// When
	// Starts listener function so that we can check what has been actually received
	go MockListener(MessageBuffer, &testBuffer)
	err := Listen(mock, MessageBuffer)

	// Then
	require.NoError(t, err)
	assert.Equal(t, testBuffer.String(), testMsg + exitString)
}

func TestListen_missedEndMarker(t *testing.T) {
	// Given
	testMsg := "Test = test <skaaa <sk> e e"
	mock := iotest.OneByteReader(strings.NewReader(testMsg))
	var MessageBuffer = make(chan string, 10)
	var testBuffer safebuffer.Buffer


	// When
	// Starts listener function so that we can check what has been actually received
	go MockListener(MessageBuffer, &testBuffer)
	err := Listen(mock, MessageBuffer)

	// Then
	require.NoError(t, err)
	assert.Equal(t, testMsg + exitString, testBuffer.String())
}

//EOF error (no error but no data returned)
func TestListen_EOF(t *testing.T) {
	// Given
	mock := iotest.ErrReader(nil)
	var MessageBuffer = make(chan string, 10)

	// When
	err := Listen(mock, MessageBuffer)

	// Then
	require.NoError(t, err)
}

//
// Test Listen()
//
func TestListen_withSimulator(t *testing.T) {
	// Given
	var MessageBuffer = make(chan string, 10)

	// When
	err := OpenAndListen("simul", nil, MessageBuffer)

	// Then
	require.NoError(t, err)
}

func MockListener(MessageBuffer chan string, workBuffer *safebuffer.Buffer) {

	for {
		var output string
		output = <-MessageBuffer
		_, err := workBuffer.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(output, "\nExiting...\n") {
			return
		}
	}

}
