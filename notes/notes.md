# Random notes

* Example for manipulating a com port: https://stackoverflow.com/questions/62659313/unable-to-parse-gps-information-via-serial-port
* Testing a serial port https://terinstock.com/post/2018/07/Black-box-Serial-Testing/

## Steps to create the project
see [How to create a CLI in golang with cobra](https://towardsdatascience.com/how-to-create-a-cli-in-golang-with-cobra-d729641c7177), Nov 2019
See also [Create a CLI in golang with Cobra](https://codesource.io/create-a-cli-in-golang-with-cobra/), March 2020 

* `go get -u github.com/spf13/cobra/cobra`
* From the `~/work` directory: `cobra init morserino_display --pkg-name morserino_display -l MIT -a ON4KJM`
* `cd morserino_display`
* `go mod init morserino_display`
* `go mod tidy`
* `go build`
* `./morserino_display`
