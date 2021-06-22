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
	"bufio"
	"fmt"

	// "os"

	// "io"
	"strings"

)


func ConsoleDisplayListener(mc *MorserinoChannels, outputStream *bufio.Writer) {
	display := &ConsoleDisplay{}
	display.w = outputStream

	AppLogger.Debug().Msg("Entering Console Display Listener")

	for {
		var output string
		output = <-mc.MessageBuffer
		display.Add(output)

		if strings.Contains(output, "\nExiting...\n") {
			mc.DisplayCompleted <- true
			AppLogger.Debug().Msg("Exiting Console Display Listener")
			return
		}
	}

}

//FIXME: Add comment
type ConsoleDisplay struct {
	currentLine strings.Builder
	newLine     string
	// output writer
	w *bufio.Writer
}

func (cd *ConsoleDisplay) String() string {
	//FIXME: add something useful here
	return cd.currentLine.String()
}

func (cd *ConsoleDisplay) Add(buff string) {
	// log.Println("ConsoleDisplay output ", cd.w)
	if strings.Contains(buff, "=") {
		//FIXME: is the buffer one char long? It is generally followed by a space
		fmt.Fprintln(cd.w, "=")
		//FIXME: better string accumulation
		cd.currentLine.WriteString("=\n")
	} else {
		_, err := fmt.Fprintf(cd.w, "%s", buff)
		if err != nil {
			// log.Fatal("Error writing to file: ", err)
			AppLogger.Error().Err(err).Msg("Error writing to file")
		}
		cd.w.Flush()
		cd.currentLine.WriteString(buff)
	}
}

//TODO: add break on column
//TODO: Add tests
