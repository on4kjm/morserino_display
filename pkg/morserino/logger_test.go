package morserino

import (
	"bufio"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupLogger(t *testing.T) {

}

func Test_getLoggerFileHandle_default(t *testing.T) {

	// ** When
	logfile, err := getLoggerFileHandle("")

	// ** Then
	require.NoError(t, err)
	assert.Regexp(t, regexp.MustCompile("morserinoTrace_20"+getTimetampRegExp()+".log"), logfile.Name())

	//cleanup
	err = os.Remove(logfile.Name())
	require.NoError(t, err)
}

func Test_getLoggerFileHandle_stdout(t *testing.T) {

	// ** When
	logfile, err := getLoggerFileHandle("StdOut")

	// ** Then
	require.NoError(t, err)
	assert.Equal(t, os.Stdout.Name(), logfile.Name())
}

func Test_getLoggerFileHandle_create(t *testing.T) {
	// ** Given
	testLogName := "test.log"
	marker := "Killroy was here"

	// ** When
	logfile, err := getLoggerFileHandle(testLogName)
	require.NoError(t, err)
	defer logfile.Close()

	w := bufio.NewWriter(logfile)
	i, err := w.WriteString(marker)
	require.NoError(t, err)

	err = w.Flush()
	require.NoError(t, err)


	// ** Then
	assert.Equal(t, testLogName, logfile.Name())

	//length should be zero as we just created it
	fi, err := logfile.Stat()
	require.NoError(t, err)
	assert.Equal(t,int64(i),fi.Size())


	//cleanup
	err = os.Remove(logfile.Name())
	require.NoError(t, err)
}

func Test_createUniqueFilename(t *testing.T) {

	// ** When
	result := createUniqueFilename()

	// ** Then
	assert.Regexp(t, regexp.MustCompile("morserinoTrace_20"+getTimetampRegExp()+".log"), result)
}

func getTimetampRegExp() string {
	year := "[0-9][0-9]"
	month := "[0-1][0-9]"
	day := "[0-3][0-9]"
	time := "[0-5][0-9][0-5][0-9][0-5][0-9]"
	return (year + month + day + time)
}
