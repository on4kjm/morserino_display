package morserino_http

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
	"net/http"
)

// We set our Read and Write buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

// our new serveWs function 
func serveWs(w http.ResponseWriter, r *http.Request) {
	// we "upgrade" the connection
	// from a standard HTTP connection to a websocket one
	ws, err := morserino_websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}
	// we then call our Writer function
	// which continually polls the serial port and writes the results
	// to this websocket connection
	go websocket.Writer(ws)
}


//TODO: specify serial port
func Start_web() {
	fmt.Println("Starting web server")

	//go serial_mock.Receive_loop()

	http.HandleFunc("/", serveHome)
	http.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir("./fonts/"))))
	http.HandleFunc("/morserino", serveWs)
	log.Fatal(http.ListenAndServe(":7373", nil))
}