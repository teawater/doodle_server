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

	log.Println("Client:", ws.RemoteAddr(), ws.RemoteAddr().Network(), ws.RemoteAddr().String())

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			log.Println("Can't receive")
			break
		}

		log.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		log.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			log.Println("Can't send")
			break
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(onConnected))

	log.Fatal(http.ListenAndServeTLS(":443", ssl_crt, ssl_key, nil))
}
