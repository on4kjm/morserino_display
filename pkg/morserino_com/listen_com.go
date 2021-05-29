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
	"log"
	"strings"

	"github.com/on4kjm/morserino_display/pkg/morserino_console"

	"go.bug.st/serial"
)

// Main listen function with display to the console
func Listen_console(morserinoPortName string, genericEnumPorts comPortEnumerator) error {

	//Port parameters for a Morserino
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	//If portname "auto" was specified, we scan for the Morserino port
	if strings.ToUpper(morserinoPortName) == "AUTO" {
		theComPortList, err := Get_com_list(genericEnumPorts)
		if err != nil {
			return err
		}

		if theComPortList.nbrOfMorserinoPorts == 0 {
			fmt.Println("Didn't find a connected Morserino!")
			return fmt.Errorf("Did not find a usable port.")
		}

		if theComPortList.nbrOfMorserinoPorts > 1 {
			fmt.Println("ERROR: Multiple Morserino devices found.")
			return fmt.Errorf("ERROR: Multiple Morserino devices found.")
		}

		morserinoPortName = theComPortList.morserinoPortName
		fmt.Println("Automatically detected Morserino port: " + morserinoPortName + "\n")
	}

	fmt.Println("Listening to port \"" + morserinoPortName + "\"")
	myPort, err := genericOpenComPort.Open(morserinoPortName, mode)
	if err != nil {
		log.Fatal(err)
	}

	consoleDisplay := morserino_console.ConsoleDisplay{}

	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := myPort.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		consoleDisplay.Add(string(buff[:n]))
	}
	return nil
}
