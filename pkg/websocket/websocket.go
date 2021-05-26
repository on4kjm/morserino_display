package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	// "github.com/on4kjm/websocket/serial_mock"
)

// We set our Read and Write buffer sizes
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// The Upgrade function will take in an incoming request and upgrade the request
// into a websocket connection
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	// this line allows other origin hosts to connect to our
	// websocket server
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// creates our websocket connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return ws, err
	}
	// returns our new websocket connection
	return ws, nil
}

func Writer(conn *websocket.Conn) {
	// we want to kick off a for loop that runs for the
	// duration of our websockets connection
	for {

		var output string
		output = <-serial_mock.MessageBuffer
		fmt.Println("Received output " + output)

		// and finally we write this JSON string to our WebSocket
		// connection and record any errors if there has been any
		if err := conn.WriteMessage(websocket.TextMessage, []byte(output)); err != nil {
			fmt.Println(err)
			return
		}

	}
}
