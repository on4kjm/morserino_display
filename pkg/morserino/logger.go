package morserino

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

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

var AppLogger zerolog.Logger

// Configures the logging subsystem following the CLI parameter
func SetupLogger(morserinoDebugLevel string, morserinoDebugFilename string) {
	//is it set?
	//if debugleve is set to trace, add the code line number

}

//Creates or opens the logger file
func getLoggerFileHandle(morserinoDebugFilename string) (*os.File, error) {
	//if "stdout", direct the logger output to it
	if strings.ToLower(morserinoDebugFilename) == "stdout" {
		// tempFile, err := ioutil.TempFile(os.TempDir(),"deleteme")
		return os.Stdout, nil
	}

	//if "", create a filename
	if morserinoDebugFilename == "" {
		morserinoDebugFilename = createUniqueFilename()
	}

	//if a filename is specified create or append to the file
	return os.OpenFile(morserinoDebugFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

}

//Generates a file name on the format "morserinoTrace_yyyymmddhhmmss.log"
func createUniqueFilename() string {
	//get current time
	dt := time.Now()

	time := dt.Format("20060102150405")
	return "morserinoTrace_" + time + ".log"
}
