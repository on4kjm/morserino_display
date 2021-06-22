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
	"io"
	"strings"
	"testing/iotest"

	// "github.com/rs/zerolog"
	"go.bug.st/serial"
)


const exitString string = "\nExiting...\n"

// Main listen function with display to the console
func OpenAndListen(morserinoPortName string, genericEnumPorts comPortEnumerator, channels *MorserinoChannels) error {

	//If requested, use the simulator instead of a real Morserino
	if strings.HasPrefix("SIMULATOR", strings.ToUpper(morserinoPortName)) {
		AppLogger.Debug().Msg("Simulator mode listener")
		TestMessage := "cq cq de on4kjm on4kjm = tks fer call om = ur rst 599 = hw? \n73 de on4kjm = <sk> e e"
		return Listen(iotest.OneByteReader(strings.NewReader(TestMessage)), channels)
	}

	//If portname "auto" was specified, we scan for the Morserino port
	if strings.ToUpper(morserinoPortName) == "AUTO" {
		AppLogger.Debug().Msg("Tying to detect the morsorino port")
		portName, err := DetectDevice(genericEnumPorts)
		if err != nil {
			return err
		}
		morserinoPortName = portName
	}

	AppLogger.Info().Msg("Listening to port \"" + morserinoPortName + "\"")

	//Port parameters for a Morserino
	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	AppLogger.Debug().Msg("Trying to open " + morserinoPortName)
	p, err := serial.Open(morserinoPortName, mode)
	if err != nil {
		AppLogger.Error().Err(err).Msg("Error opening port")
		return err
	}
	defer p.Close()

	return Listen(p, channels)
}

// Main receive loop
func Listen(port io.Reader, channels *MorserinoChannels) error {

	// //TODO: needs to be moved as a goroutine
	// consoleDisplay := morserino_console.ConsoleDisplay{}

	// variables for tracking the exit pattern
	var (
		closeKey            string
		possibleExitRequest bool
		closeRequested      bool
	)

	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := port.Read(buff)
		if err != nil {
			AppLogger.Error().Err(err).Msg("Error reading on port")
		}

		// Check whether the "end of transmission" was sent
		// TODO: move this in a seperate structure/function
		if string(buff[0:1]) == "<" {
			closeKey = "<"
			possibleExitRequest = true
		} else {
			if possibleExitRequest {
				closeKey = closeKey + string(buff[:n])
				if !strings.HasPrefix("<sk> e e", closeKey) {
					possibleExitRequest = false
				} else {
					if closeKey == "<sk> e e" {
						closeRequested = true
					}
				}
			}
		}

		if n == 0 {
			AppLogger.Debug().Msg("EOF detectd")
			AppLogger.Trace().Msg("Sending EOF to the console displayer")
			channels.MessageBuffer <- "\nEOF"
			//sending the exit marker to the diplay goroutine
			AppLogger.Trace().Msg("Sending exit marker to the console displayer")
			channels.MessageBuffer <- exitString
			//waiting for it to complete (blocking read)
			AppLogger.Debug().Msg("Waiting for the signal that the display processing was completed")
			<-channels.DisplayCompleted
			AppLogger.Debug().Msg("Display processing completed (received signal)")
			break
		}

		channels.MessageBuffer <- string(buff[:n])

		if closeRequested {
			//sending the exit marker to the diplay goroutine
			AppLogger.Trace().Msg("Sending exit marker to the console displayer as we received the exit sequence")
			channels.MessageBuffer <- exitString
			//waiting for it to complete (blocking read)
			AppLogger.Debug().Msg("Waiting for the signal that the display processing was completed")
			<-channels.DisplayCompleted
			AppLogger.Debug().Msg("Display processing completed (received signal)")
			break
		}
	}
	AppLogger.Debug().Msg("Sending signal that all processing is done")
	channels.Done <- true

	AppLogger.Debug().Msg("Exiting Listen")
	return nil
}

// Tries to auto detect the Morserino port based on the USB Vendor and Device ID
func DetectDevice(genericEnumPorts comPortEnumerator) (string, error) {
	theComPortList, err := Get_com_list(genericEnumPorts)
	if err != nil {
		return "", err
	}

	if theComPortList.nbrOfMorserinoPorts == 0 {
		fmt.Println("Didn't find a connected Morserino!")
		return "", fmt.Errorf("Did not find a usable port.")
	}

	if theComPortList.nbrOfMorserinoPorts > 1 {
		fmt.Println("ERROR: Multiple Morserino devices found.")
		return "", fmt.Errorf("ERROR: Multiple Morserino devices found.")
	}

	morserinoPortName := theComPortList.morserinoPortName
	fmt.Println("Automatically detected Morserino port: " + morserinoPortName + "\n")

	return morserinoPortName, nil
}
