package morserino

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func init() {
    // create a temp file
    // tempFile, err := ioutil.TempFile(".","testLog")
    // if err != nil {
    //     // Can we log an error before we have our logger? :)
    //     log.Error().Err(err).Msg("there was an error creating a temporary file four our log")
    // }

	zerolog.TimeFieldFormat = "15:04:05.000"
    AppLogger= zerolog.New(os.Stdout).With().Timestamp().Logger()

    // fmt.Printf("The log file is allocated at %s\n", tempFile.Name())

	zerolog.SetGlobalLevel(zerolog.TraceLevel)

}

func TestMorserino_console(t *testing.T) {
	AppLogger.Info().Msg("==> " + t.Name())

	// ** Given

	// ** When
	Morserino_console("simulator")

	// ** Then

	AppLogger.Info().Msg("<== " + t.Name())
	// assert.Equal(t,true,false)
	assert.Fail(t,"Stop")
}
