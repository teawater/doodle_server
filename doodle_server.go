package main

import (
	"log"
	"net/http"
	"sync"
	"time"
	"encoding/json"
	"container/list"

	"golang.org/x/net/websocket"
)

const ssl_crt = "server.crt"
const ssl_key = "server.key"

const color_x = 750
const color_y = 750

type color_s struct {
	lock sync.RWMutex
	id uint64_t
	data [color_x][color_y]string
}
var color color_s

var clients *list.List
var clients_lock sync.Mutex

func onConnected(ws *websocket.Conn) {
	var err error
	var reply string
	type reply_pkg_s struct {
		fmt int
		id uint64
	}
	var reply_pkg reply_pkg_s

	//check the client
	log.Println("Client:", ws.RemoteAddr(), ws.RemoteAddr().Network(), ws.RemoteAddr().String())

	//Handle first pack
	ws.SetDeadline(time.Now().Add(time.Second * 60))
	if err = websocket.Message.Receive(ws, &reply); err != nil {
		log.Println("Get first packet error:", err)
		return
	}
	err = json.Unmarshal([]byte(reply), &reply_pkg)
	if reply_pkg.fmt != 0 {
		log.Println("Get first packet error:", reply_pkg)
		return
	}

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
