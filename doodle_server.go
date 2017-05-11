package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

const ssl_crt = "server.crt"
const ssl_key = "server.key"

func onConnected(ws *websocket.Conn) {
	var err error

	log.Println("Client:", ws)

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(onConnected))

	log.Fatal(http.ListenAndServeTLS(":8080", ssl_crt, ssl_key, nil))
}
