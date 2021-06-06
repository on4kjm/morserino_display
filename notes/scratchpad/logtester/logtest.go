package main

import (
	"fmt"
	"io/ioutil"

	"os"

	// "strings"
	// "time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	fmt.Println("Killroy was here")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    // create a temp file
    tempFile, err := ioutil.TempFile(os.TempDir(),"deleteme")
    if err != nil {
        // Can we log an error before we have our logger? :)
        log.Error().Err(err).Msg("there was an error creating a temporary file four our log")
    }
    fileLogger := zerolog.New(tempFile).With().Timestamp().Logger()
    fileLogger.Info().Msg("This is an entry from my log")

    fmt.Printf("The log file is allocated at %s\n", tempFile.Name())
	
	// output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05", NoColor: false}
	// output.FormatLevel = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	// }
	// // output.FormatMessage = func(i interface{}) string {
	// // 	return fmt.Sprintf("***%s****", i)
	// // }
	// output.FormatFieldName = func(i interface{}) string {
	// 	return fmt.Sprintf("%s:", i)
	// }
	// output.FormatFieldValue = func(i interface{}) string {
	// 	return strings.ToUpper(fmt.Sprintf("%s", i))
	// }
	
	// log := zerolog.New(output).With().Timestamp().Logger()
	
	fileLogger.Info().Str("foo", "bar").Msg("Hello World")
	fileLogger.Debug().Msg("This is a debug message")

}