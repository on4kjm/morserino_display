package morserino_com

import (
	"reflect"
	"testing"
	// "go.bug.st/serial/enumerator"
	//"github.com/stretchr/testify/mock"
)

// This helps in assigning mock at the runtime instead of compile time
// var userExistsMock func(email string) bool

// type preCheckMock struct{}

// func (u preCheckMock) userExists(email string) bool {
// 	return userExistsMock(email)
// }

// type enumeratePortsMock struct {}

// func (e enumeratePortsMock) GetDetailedPortsList() ([]*enumerator.PortDetails, error) {
// 	return mockGetDetailedPortsList()
// }

// func TestGet_com_list(t *testing.T) {

// 	genericEnumPorts = enumeratePortsMock{}
// 	mockGetDetailedPortsList = func() ([]*enumerator.PortDetails, error) {
// 		return false
// 	}

// }

// func TestGet_com_list(t *testing.T) {
// 	mockEnumPorts := new(enumeratePortsInterface)

// 	mockEnumPorts.

// }

func Test_prettyPrint_comList(t *testing.T) {
	type args struct {
		portList comPortList
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prettyPrint_comList(tt.args.portList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("prettyPrint_comList() = %v, want %v", got, tt.want)
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
