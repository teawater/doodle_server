package main

import (
	"log"
	"net/http"
	"sync"
	"time"
	"encoding/json"

	"golang.org/x/net/websocket"
)

const ssl_crt = "server.crt"
const ssl_key = "server.key"

const color_x = 750
const color_y = 750
var color_lock sync.RWMutex
var color [color_x][color_y]string

func onConnected(ws *websocket.Conn) {
	var err error
	var reply string

	//check the client
	log.Println("Client:", ws.RemoteAddr(), ws.RemoteAddr().Network(), ws.RemoteAddr().String())

	//Handle first pack
	ws.SetDeadline(time.Now().Add(time.Second * 60))
	if err = websocket.Message.Receive(ws, &reply); err != nil {
		log.Println("Get first packet error:", err)
		return
	}
	err = json.Unmarshal([]byte(reply), &new)

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
	for x := 0; x < color_x; x++ {
		for y := 0; y < color_x; y++ {
			color[x][y] = "#ffffff"
		}
	}

	http.Handle("/", websocket.Handler(onConnected))

	log.Fatal(http.ListenAndServeTLS(":443", ssl_crt, ssl_key, nil))
}
