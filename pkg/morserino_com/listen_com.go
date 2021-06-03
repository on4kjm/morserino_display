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

// Channel used to forward data received on the serial port to the display modules
var MessageBuffer = make(chan string, 10)

//
// Interfaces for COM related stuff
// (inspired by https://stackoverflow.com/questions/41053280/how-to-write-mock-for-structs-in-go)
//

// Morserino Port interface
type IMorserinoPort interface {
	MorserinoRead(p []byte)  (n int, err error)
	MorserinoClose() error
}

//Real implementation
type MorserinoPort struct {
	Port serial.Port
}

func (mp MorserinoPort) MorserinoRead(p []byte) (n int, err error) {
	return (mp.Port.Read(p))
}

func (mp MorserinoPort) MorserinoClose() error {
	return (mp.Port.Close())
}


// Open MorserinoPort interface
type IMorserinoPortOpen interface {
	MorserinoOpen(portName string, mode *serial.Mode) (IMorserinoPort, error)
}

// type MorserinoOpenWraper struct {
// 	OpenMorserinoPort
// }

// func (w MorserinoOpenWraper) MorserinoOpenW(portName string, mode *serial.Mode) (IMorserinoPort, error) {
// 	return w.OpenMorserinoPort.MorserinoOpen(portName, mode)

// Real implementation
type OpenMorserinoPort struct{}

func (r OpenMorserinoPort) MorserinoOpen(portName string, mode *serial.Mode) (serial.Port, error) {
	return serial.Open(portName, mode)
}

//
// Business logic
//

// Main listen function with display to the console
func Listen_console(morserinoPortName string, genericEnumPorts comPortEnumerator, genericOpenComPort IMorserinoPortOpen, genericComPort IMorserinoPort) error {

	//Port parameters for a Morserino
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	//TODO: implement the simulator

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

	log.Println("Listening to port \"" + morserinoPortName + "\"")


	myPort, err := genericOpenComPort.MorserinoOpen(morserinoPortName, mode)
	if err != nil {
		log.Fatal(err)
	}
	genericComPort.Port = myPort

	//We want to make sure that we close the port
	defer genericComPort.MorserinoClose()

	//TODO: needs to be moved as a goroutine
	consoleDisplay := morserino_console.ConsoleDisplay{}

	var closeKey string
	possibleExitRequest := false
	closeRequested := false
	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := genericComPort.MorserinoRead(buff)
		if err != nil {
			log.Fatal(err)
		}

		// Check whether the "end of transmission" was sent
		// TODO: move this in a seperate structure/function
		if string(buff[0:1]) == "<" {
			closeKey = "<"
			possibleExitRequest = true
		} else {
			if possibleExitRequest {
				closeKey = closeKey + string(buff[:n])
				if(!strings.HasPrefix("<sk> e e", closeKey)) {
					possibleExitRequest = false
				} else {
					if(closeKey == "<sk> e e") {
						closeRequested = true
					}
				}
			}
		}

		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		// TODO: move this out and use a channel for that
		consoleDisplay.Add(string(buff[:n]))

		if closeRequested {
			consoleDisplay.Add("\nExiting...\n")
			break
		}
	}
	return nil
}
