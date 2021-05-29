package morserino_com

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.bug.st/serial/enumerator"
)

//
// ============ Mock definitions ==================

//Happy case mock
type mockEnumeratePorts_OK struct{}

//Happy case mocked method
func (e mockEnumeratePorts_OK) GetDetailedPortsList() ([]*enumerator.PortDetails, error) {
	happyCaseList := []*enumerator.PortDetails{}

	item1 := enumerator.PortDetails{
		Name:         "COM Port Name",
		IsUSB:        true,
		VID:          "XXXX",
		PID:          "YYYY",
		SerialNumber: "001",
	}
	item2 := enumerator.PortDetails{
		Name:         "Morserino Port",
		IsUSB:        true,
		VID:          "10C4",
		PID:          "EA60",
		SerialNumber: "002",
	}
	happyCaseList = append(happyCaseList, &item1)
	happyCaseList = append(happyCaseList, &item2)

	return happyCaseList, nil
}

// Mock that returns an error
type mockEnumeratePorts_error struct{}

// Mocked method that returns an error
func (e mockEnumeratePorts_error) GetDetailedPortsList() ([]*enumerator.PortDetails, error) {
	return nil, fmt.Errorf("An error occured")
}

// Mock that found no ports
type mockEnumeratePorts_noPort struct{}

// Mocked method that found no ports
func (e mockEnumeratePorts_noPort) GetDetailedPortsList() ([]*enumerator.PortDetails, error) {
	emptyList := []*enumerator.PortDetails{}

	return emptyList, nil
}

//
// =================  Tests ====================
//

func Test_Get_com_list_HappyCase(t *testing.T) {

	//Happy case mock
	var mockEnumPorts mockEnumeratePorts_OK

	// System under test
	comList, err := Get_com_list(mockEnumPorts)

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
	if assert.NoError(t, err) {
		assert.Equal(t, comList.nbrOfPorts, 2)
		assert.Equal(t, comList.nbrOfMorserinoPorts, 1)
		assert.Equal(t, comList.morserinoPortName, "Morserino Port")
		assert.Equal(t, comList.portList, targetPortList)
	}
}

func Test_Get_com_list_Error(t *testing.T) {
	//Mock that returns an error
	var mockEnumPorts mockEnumeratePorts_error

	// System under test
	_, err := Get_com_list(mockEnumPorts)

	//validate results
	assert.Error(t, fmt.Errorf("An error occured"), err)
}

func Test_Get_com_list_noPort(t *testing.T) {
	//Mock that found no ports
	var mockEnumPorts mockEnumeratePorts_noPort

	// System under test
	comList, err := Get_com_list(mockEnumPorts)

	//validate results
	if assert.NoError(t, err) {
		assert.Equal(t, 0, comList.nbrOfPorts)
	}
}

// ===============================

func Test_List_com_happyCase(t *testing.T) {
	// Happy case mock
	var mockEnumPorts mockEnumeratePorts_OK

	// System under test with mock
	output, err := List_com(mockEnumPorts)

	// Validating results
	if assert.NoError(t, err) {
		expectedOutput := "\n    COM Port Name                  (USB ID: XXXX:YYYY, USB Serial: 001)\n=>  Morserino Port                 (USB ID: 10C4:EA60, USB Serial: 002)"
		assert.Equal(t, expectedOutput, output)
	}
}

func Test_List_com_error(t *testing.T) {
	// Mock that returns an error
	var mockEnumPorts mockEnumeratePorts_error

	// System under test with mock
	output, err := List_com(mockEnumPorts)

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := format_com_item(tt.args.item); got != tt.want {
				t.Errorf("format_com_item() = %v, want %v", got, tt.want)
			}
		})
	}
}
