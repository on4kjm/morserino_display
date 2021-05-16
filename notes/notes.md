# Random notes

## COM port manipulations

* Example for manipulating a com port: https://stackoverflow.com/questions/62659313/unable-to-parse-gps-information-via-serial-port
* https://stackoverflow.com/questions/50318058/reading-from-arduinos-serial-port-using-go-serial/50330709
* https://stackoverflow.com/questions/17599232/reading-from-serial-port-with-while-loop
* https://stackoverflow.com/questions/50649256/go-parsed-serial-port-input-using-goroutines-not-printing
* https://stackoverflow.com/questions/50088669/golang-reading-from-serial
* https://github.com/chrissnell/tnc-server (outdated serial library)
* https://github.com/tv42/topic
* https://golang.hotexamples.com/examples/github.com.tarm.serial/Port/Read/golang-port-read-method-examples.html
* https://pkg.go.dev/go.bug.st/serial
* Testing a serial port https://terinstock.com/post/2018/07/Black-box-Serial-Testing/
* `screen /dev/tty.usbserial-0001 115200` also  Parity.None, 8, StopBits.One

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

* What is CGO ?
  ```
  This library tries to avoid the use of the "C" package (and consequently the need of cgo) to simplify cross compiling. Unfortunately the USB enumeration package for darwin (MacOSX) requires cgo to access the IOKit framework. This means that if you need USB enumeration on darwin you're forced to use cgo.
  ```

* https://github.com/goreleaser/goreleaser-action/issues/233

## Mocking articles
* [How to mock? Go Way.](https://medium.com/@ankur_anand/how-to-mock-in-your-go-golang-tests-b9eee7d7c266)
* [GoMock vs. Testify: Mocking frameworks for Go](https://blog.codecentric.de/2019/07/gomock-vs-testify/)
* [Writing Go Code with Visual Studio Code](https://medium.com/@mekilis/writing-go-code-with-visual-studio-code-6daa326cb6b8)
* [How to mock an external service for tests in GO](https://wawand.co/blog/posts/how-to-mock-an-external-service-for-test-in-go-73251a7a/)
* [Mocking dependencies in Go](https://dev.to/jonfriesen/mocking-dependencies-in-go-1h4d)
* [How do I mock a function from another package without using dependency injection?](https://stackoverflow.com/questions/51428617/how-do-i-mock-a-function-from-another-package-without-using-dependency-injection)
* [How to write mock for structs in Go](https://stackoverflow.com/questions/41053280/how-to-write-mock-for-structs-in-go)
* [Improving Your Go Tests and Mocks With Testify](https://tutorialedge.net/golang/improving-your-tests-with-testify-go/)
* [Mockery ](https://github.com/vektra/mockery)
* [Mocking Techniques for Go](https://www.myhatchpad.com/insight/mocking-techniques-for-go/)

## Miscelaneous
* [Buy Me a Coffe badge](https://gist.github.com/gbraad/216f8162d9b382d14b8a097a37cc2c72#file-readme-md)
* https://github.com/troian/golang-cross

## Web interface
* conversation sur Slack
  * "je dirais "vanilla js" ce qui veut dire aucun framework, et requête ajax oui"
  * si tu veux pas te prendre la tête je dirais ajax aussi
  * https://www.w3schools.com/xml/ajax_intro.asp
  * example
    ```
    const interval = window.setInterval(() => {callServer()}, 1000)
    window.clearInterval(interval); quand il veut stop
    function callServer(){
    fetch('flowers.jpg')
    .then(function(response) {
      return response.blob();
    })
    .then(function(myBlob) {
      const objectURL = URL.createObjectURL(myBlob);
      myImage.src = objectURL;
    });
    }
    ```
  * https://developer.mozilla.org/fr/docs/Web/API/Fetch_API/Using_Fetch
  * Lecture série -> go Channel -> websocket <- js
* [WebSocket simple example / server in go, client in JS](https://gist.github.com/owulveryck/57d8c2469fd1f8a840747b064c50ff4e)
* [Go: Asynchronous Real-Time Broadcasting using Channels and WebSockets](https://dev.to/danielkun/go-asynchronous-and-safe-real-time-broadcasting-using-channels-and-websockets-4g5d), 2018
* [Simple Go-based HTTP streaming via HTTP and websockets.](https://gist.github.com/vmarmol/b967b29917a34d9307ce)
* [HTTP Streaming in Golang](https://dmathieu.com/articles/development/golang-streaming/)
* [Create a Single Page App With Go, Echo and Vue](https://scotch.io/tutorials/create-a-single-page-app-with-go-echo-and-vue)
* [How to update the textbox value without page refresh using ajax in mvc](https://stackoverflow.com/questions/42396886/how-to-update-the-textbox-value-without-page-refresh-using-ajax-in-mvc)
* [Using websockets in golang with gorilla-websocket](https://rogerwelin.github.io/golang/websockets/gorilla/2018/03/13/golang-websockets.html)
* see Gorilla-websocket examples: https://github.com/gorilla/websocket/blob/master/examples/filewatch/main.go
