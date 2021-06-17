package morserino_com

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/on4kjm/morserino_display/pkg/morserino_channels"
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
	var testBuffer safebuffer.Buffer
	var mc morserino_channels.MorserinoChannels
	mc.Init()


	// When
	// Starts listener function so that we can check what has been actually received
	go MockListener(mc, &testBuffer)
	err := Listen(mock, mc)

	// Then
	require.NoError(t, err)
	assert.Equal(t, testBuffer.String(), testMsg+exitString)
}

func TestListen_missedEndMarker(t *testing.T) {
	// Given
	testMsg := "Test = test <skaaa <sk> e e"
	mock := iotest.OneByteReader(strings.NewReader(testMsg))
	var testBuffer safebuffer.Buffer
	var mc morserino_channels.MorserinoChannels
	mc.Init()

	// When
	// Starts listener function so that we can check what has been actually received
	go MockListener(mc, &testBuffer)
	err := Listen(mock, mc)

	// Then
	require.NoError(t, err)
	assert.Equal(t, testMsg+exitString, testBuffer.String())
}

//EOF error (no error but no data returned)
func TestListen_EOF(t *testing.T) {
	// Given
	testMsg := "\nEOF"
	mock := iotest.ErrReader(nil)
	var testBuffer safebuffer.Buffer
	var mc morserino_channels.MorserinoChannels
	mc.Init()


	// When
	go MockListener(mc, &testBuffer)
	err := Listen(mock, mc)

	// Then
	require.NoError(t, err)
	assert.Equal(t, testMsg+exitString, testBuffer.String())
}

//
// Test Listen()
//
func TestListen_withSimulator(t *testing.T) {
	// Given
	testMsg := "cq cq de on4kjm on4kjm = tks fer call om = ur rst 599 = hw? \n73 de on4kjm = <sk> e e"
	var testBuffer safebuffer.Buffer
	var mc morserino_channels.MorserinoChannels
	mc.Init()

	// When
	// Starts listener function so that we can check what has been actually received
	go MockListener(mc, &testBuffer)
	err := OpenAndListen("simul", nil, mc)

	// Then
	require.NoError(t, err)
	assert.Equal(t, testBuffer.String(), testMsg+exitString)
}

func MockListener(mc morserino_channels.MorserinoChannels, workBuffer *safebuffer.Buffer) {

	for {
		var output string
		output = <- mc.MessageBuffer
		_, err := workBuffer.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}

		if strings.Contains(output, "\nExiting...\n") {
			mc.DisplayCompleted <- true
			return
		}
	}

}
