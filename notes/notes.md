# Random notes

## COM port manipulations

* Example for manipulating a com port: https://stackoverflow.com/questions/62659313/unable-to-parse-gps-information-via-serial-port
* https://stackoverflow.com/questions/50318058/reading-from-arduinos-serial-port-using-go-serial/50330709
* https://stackoverflow.com/questions/17599232/reading-from-serial-port-with-while-loop
* https://stackoverflow.com/questions/50649256/go-parsed-serial-port-input-using-goroutines-not-printing
* https://stackoverflow.com/questions/50088669/golang-reading-from-serial
* https://github.com/chrissnell/tnc-server 
* https://github.com/tv42/topic
* https://golang.hotexamples.com/examples/github.com.tarm.serial/Port/Read/golang-port-read-method-examples.html
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
