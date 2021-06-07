package morserino_com

import (
	"fmt"
	"strings"
	"testing"
	"testing/iotest"

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
	mock := iotest.OneByteReader(strings.NewReader("Test = test <sk> e e"))

	// When
	err := Listen(mock)

	// Then
	require.NoError(t, err)
}

func TestListen_missedEndMarker(t *testing.T) {
	// Given
	mock := iotest.OneByteReader(strings.NewReader("Test = test <skaaa <sk> e e"))

	// When
	err := Listen(mock)

	// Then
	require.NoError(t, err)
}

//EOF error (no error but no data returned)
func TestListen_EOF(t *testing.T) {
	// Given
	mock := iotest.ErrReader(nil)

	// When
	err := Listen(mock)

	// Then
	require.NoError(t, err)
}

//
// Test ConsoleListen()
//
func TestConsoleListen_withSimulator(t *testing.T) {
	// Given

	// When
	err := ConsoleListen("simul", nil)

	// Then
	require.NoError(t, err)
}
