package morserino_core

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
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/on4kjm/morserino_display/pkg/morserino_channels"
	"github.com/on4kjm/morserino_display/pkg/morserino_com"
	"github.com/on4kjm/morserino_display/pkg/morserino_console"
)

// Main entry point for console output
func Morserino_console(morserinoPortName string) {

	// initialize the structure containing all the channels we are going to use
	channels := &morserino_channels.MorserinoChannels{}
	channels.Init()

	// Setting up the EnumPorts to the "real life" implementation
	var realEnumPorts morserino_com.EnumeratePorts

	go morserino_com.OpenAndListen(morserinoPortName, realEnumPorts, channels)
	go morserino_console.ConsoleDisplayListener(channels, bufio.NewWriter(os.Stdout))

	<-channels.Done //Waiting here for everything to be orderly completed
}

//Main entry point for listing ports
func Morserino_list() {
	//We are going to use the real function to enumerate ports
	var realEnumPorts morserino_com.EnumeratePorts

	//Get the pretty printed list of devices
	output, err := morserino_com.List_com(realEnumPorts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
