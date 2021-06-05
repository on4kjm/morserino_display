package morserino_com

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
	"testing"

	"go.bug.st/serial/enumerator"
)

func TestListen_console(t *testing.T) {

}

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
