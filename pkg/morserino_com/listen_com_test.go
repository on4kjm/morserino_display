package morserino_com

import (
	"fmt"
	"io"
	"strings"
	"testing"
	"testing/iotest"

	"go.bug.st/serial/enumerator"
)

func TestGet_MorserinoPort_automatically(t *testing.T) {
	type args struct {
		genericEnumPorts comPortEnumerator
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"Happy case",
			args{
				mockPortEnumerator{ports: []*enumerator.PortDetails{
					{Name: "COM Port Name", IsUSB: true, VID: "XXXX", PID: "YYYY", SerialNumber: "001"},
					{Name: "Morserino Port", IsUSB: true, VID: "10C4", PID: "EA60", SerialNumber: "002"}},
				},
			},
			"Morserino Port", false,
		},
		{
			"No morserino detected",
			args{
				mockPortEnumerator{ports: []*enumerator.PortDetails{
					{Name: "COM Port Name", IsUSB: true, VID: "XXXX", PID: "YYYY", SerialNumber: "001"}},
				},
			},
			"", true,
		},
		{
			"Too many Morserinos detected",
			args{
				mockPortEnumerator{ports: []*enumerator.PortDetails{
					{Name: "COM Port Name", IsUSB: true, VID: "XXXX", PID: "YYYY", SerialNumber: "001"},
					{Name: "Morserino Port", IsUSB: true, VID: "10C4", PID: "EA60", SerialNumber: "002"},
					{Name: "Morserino Port2", IsUSB: true, VID: "10C4", PID: "EA60", SerialNumber: "003"}},
				},
			},
			"", true,
		},
		{
			"Port enumeration went wrong",
			args{mockPortEnumerator{ports: nil, err: fmt.Errorf("An error occured")}},
			"", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get_MorserinoPort_automatically(tt.args.genericEnumPorts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get_MorserinoPort_automatically() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get_MorserinoPort_automatically() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListen(t *testing.T) {
	type args struct {
		port io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Happy case",
			args{iotest.OneByteReader(strings.NewReader("Test = test <sk> e e"))},
			false,
		},
		{
			"missed end marker",
			args{iotest.OneByteReader(strings.NewReader("Test = test <skaaa <sk> e e"))},
			false,
		},
		{
			"EOF error (no error but no dat returned",
			args{iotest.ErrReader(nil)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Listen(tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Listen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListen_console(t *testing.T) {
	type args struct {
		morserinoPortName string
		genericEnumPorts  comPortEnumerator
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Test simulator",
			args{morserinoPortName: "simul", genericEnumPorts: nil},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Listen_console(tt.args.morserinoPortName, tt.args.genericEnumPorts); (err != nil) != tt.wantErr {
				t.Errorf("Listen_console() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
