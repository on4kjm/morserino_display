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
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.bug.st/serial/enumerator"
)

func Test_Get_com_list_HappyCase(t *testing.T) {

	// Setup mock
	portEnumerator := mockPortEnumerator{
		ports: []*enumerator.PortDetails{
			{
				Name:         "COM Port Name",
				IsUSB:        true,
				VID:          "XXXX",
				PID:          "YYYY",
				SerialNumber: "001",
			},
			{
				Name:         "Morserino Port",
				IsUSB:        true,
				VID:          "10C4",
				PID:          "EA60",
				SerialNumber: "002",
			},
		},
	}

	// System under test
	comList, err := Get_com_list(portEnumerator)

	//Expected Port List
	targetPortList := []comPortItem{
		{
			portName:        "COM Port Name",
			usbVendorID:     "XXXX",
			usbProductID:    "YYYY",
			serialNumber:    "001",
			isMorserinoPort: false,
		},
		{
			portName:        "Morserino Port",
			usbVendorID:     "10C4",
			usbProductID:    "EA60",
			serialNumber:    "002",
			isMorserinoPort: true,
		},
	}

	//validate results
	require.NoError(t, err)

	assert.Equal(t, comList.nbrOfPorts, 2)
	assert.Equal(t, comList.nbrOfMorserinoPorts, 1)
	assert.Equal(t, comList.morserinoPortName, "Morserino Port")
	assert.Equal(t, comList.portList, targetPortList)
}

func Test_Get_com_list_Error(t *testing.T) {
	//Mock that returns an error
	portEnumerator := mockPortEnumerator{
		ports: nil,
		err:   fmt.Errorf("An error occured"),
	}

	// System under test
	_, err := Get_com_list(portEnumerator)

	//validate results
	assert.Error(t, fmt.Errorf("An error occured"), err)
}

func Test_Get_com_list_noPort(t *testing.T) {
	//Mock that found no ports
	portEnumerator := mockPortEnumerator{
		ports: []*enumerator.PortDetails{},
		err:   nil,
	}

	// System under test
	comList, err := Get_com_list(portEnumerator)

	//validate results
	require.NoError(t, err)
	assert.Equal(t, 0, comList.nbrOfPorts)
}

// ===============================

func Test_List_com_happyCase(t *testing.T) {
	// Happy case mock
	portEnumerator := mockPortEnumerator{
		ports: []*enumerator.PortDetails{
			{
				Name:         "COM Port Name",
				IsUSB:        true,
				VID:          "XXXX",
				PID:          "YYYY",
				SerialNumber: "001",
			},
			{
				Name:         "Morserino Port",
				IsUSB:        true,
				VID:          "10C4",
				PID:          "EA60",
				SerialNumber: "002",
			},
		},
	}

	// System under test with mock
	output, err := List_com(portEnumerator)

	// Validating results
	if assert.NoError(t, err) {
		expectedOutput := "\n    COM Port Name                  (USB ID: XXXX:YYYY, USB Serial: 001)\n=>  Morserino Port                 (USB ID: 10C4:EA60, USB Serial: 002)"
		assert.Equal(t, expectedOutput, output)
	}
}

func Test_List_com_error(t *testing.T) {
	//Mock that returns an error
	portEnumerator := mockPortEnumerator{
		ports: nil,
		err:   fmt.Errorf("An error occured"),
	}

	// System under test with mock
	output, err := List_com(portEnumerator)

	// Validating results
	assert.Error(t, fmt.Errorf("An error occured"), err)
	assert.Equal(t, "", output)
}

// ===============================

func Test_prettyPrint_comList(t *testing.T) {
	type args struct {
		portList comPortList
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"Happy case",
			args{portList: comPortList{
				nbrOfPorts:          2,
				nbrOfMorserinoPorts: 1,
				morserinoPortName:   "blaaah",
				portList: []comPortItem{
					{
						portName:        "COM Port Name",
						usbVendorID:     "XXXX",
						usbProductID:    "YYYY",
						serialNumber:    "001",
						isMorserinoPort: false,
					},
					{
						portName:        "Morserino Port",
						usbVendorID:     "10C4",
						usbProductID:    "EA60",
						serialNumber:    "002",
						isMorserinoPort: true,
					},
				},
			},
			},
			[]string{
				"",
				"    COM Port Name                  (USB ID: XXXX:YYYY, USB Serial: 001)",
				"=>  Morserino Port                 (USB ID: 10C4:EA60, USB Serial: 002)",
			},
		},
		{
			"too many Morserinos detected",
			args{portList: comPortList{
				nbrOfPorts:          2,
				nbrOfMorserinoPorts: 2,
				morserinoPortName:   "blaaah",
				portList: []comPortItem{
					{
						portName:        "COM Port Name",
						usbVendorID:     "XXXX",
						usbProductID:    "YYYY",
						serialNumber:    "001",
						isMorserinoPort: true,
					},
					{
						portName:        "Morserino Port",
						usbVendorID:     "10C4",
						usbProductID:    "EA60",
						serialNumber:    "002",
						isMorserinoPort: true,
					},
				},
			},
			},
			[]string{
				"WARNING: Multiple multiple Morserino devices detected",
				"",
				"=>  COM Port Name                  (USB ID: XXXX:YYYY, USB Serial: 001)",
				"=>  Morserino Port                 (USB ID: 10C4:EA60, USB Serial: 002)",
			},
		},
		{
			"No ports were detected",
			args{portList: comPortList{
				nbrOfPorts:          0,
				nbrOfMorserinoPorts: 0,
			},
			},
			[]string{
				"No ports found !",
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prettyPrint_comList(tt.args.portList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prettyPrint_comList() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func Test_format_com_item(t *testing.T) {
	type args struct {
		item comPortItem
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"Format item: happy case",
			args{item: comPortItem{
				portName:        "Port Name",
				usbVendorID:     "XXXX",
				usbProductID:    "YYYY",
				serialNumber:    "001",
				isMorserinoPort: false}},
			"    Port Name                      (USB ID: XXXX:YYYY, USB Serial: 001)",
		},
		{
			"Format item: happy case with Morserino detected",
			args{item: comPortItem{
				portName:        "Port Name",
				usbVendorID:     "XXXX",
				usbProductID:    "YYYY",
				serialNumber:    "001",
				isMorserinoPort: true}},
			"=>  Port Name                      (USB ID: XXXX:YYYY, USB Serial: 001)",
		},
		{
			"Format item: shorter VID:PID",
			args{item: comPortItem{
				portName:        "Port Name",
				usbVendorID:     "X",
				usbProductID:    "Y",
				serialNumber:    "01",
				isMorserinoPort: false}},
			"    Port Name                      (USB ID:    X:   Y, USB Serial: 01)",
		},
		{
			"Format item: non-USB device",
			args{item: comPortItem{
				portName:        "Port Name",
				isMorserinoPort: false}},
			"    Port Name                     ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format_com_item(tt.args.item); got != tt.want {
				t.Errorf("format_com_item() = %v, want %v", got, tt.want)
			}
		})
	}
}
