package morserino_com

import (
	"reflect"
	"testing"
	// "go.bug.st/serial/enumerator"
	//"github.com/stretchr/testify/mock"
)

func Test_Get_com_list_HappyCase(t *testing.T) {
	//Prepare mock

}


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
			args{portList: 
				comPortList{
					nbrOfPorts: 2, 
					nbrOfMorserinoPorts: 1, 
					morserinoPortName: "blaaah",
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
			args{portList: 
				comPortList{
					nbrOfPorts: 2, 
					nbrOfMorserinoPorts: 2, 
					morserinoPortName: "blaaah",
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
			args{portList: 
				comPortList{
					nbrOfPorts: 0, 
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
