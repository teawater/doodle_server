package main

import (
	"log"
	"net/http"
	"sync"
	"time"
	"fmt"
	"container/list"
	"encoding/json"

	"golang.org/x/net/websocket"
)

const ssl_crt = "server.crt"
const ssl_key = "server.key"

const color_x = 750
const color_y = 750

type color_s struct {
	Id uint64
	Fmt int
	Data [color_x][color_y]string
}
var color color_s
var color_lock sync.RWMutex

var clients *list.List
var clients_lock sync.Mutex

func sync_color_to_client(ws *websocket.Conn, id *uint64) (err error) {
	err = nil

	color_lock.RLock()
	defer color_lock.RUnlock()

	if color.Id == *id {
		return
	}
	if color.Id < *id {
		err = fmt.Errorf("get wrong id %d that it should small than %d", *id, color.Id)
		return
	}

	*id = color.Id
	
	err = websocket.JSON.Send(ws, color)
	// b, err := json.Marshal(color)
	// if (err != nil) {
		// log.Println(err)
	// }
	// log.Println(string(b))
	// ws.SetWriteDeadline(time.Now().Add(time.Second * 10))
	// if err = websocket.Message.Send(ws, b); err != nil {
		// log.Println("Can't send")
	// }

	return
}

func handle_receive_pack(ws *websocket.Conn, quit chan bool) {
	
}

func onConnected(ws *websocket.Conn) {
	var err error
	type reply_s struct {
		fmt int
		id uint64
	}
	var reply reply_s

	//check the client
	log.Println("Client:", ws.RemoteAddr(), ws.RemoteAddr().Network(), ws.RemoteAddr().String())

	//Handle first pack
	ws.SetReadDeadline(time.Now().Add(time.Second * 60))
	err = websocket.JSON.Receive(ws, &reply)
	if err != nil {
		log.Println("Get first packet error:", err)
		return
	}
	if reply.fmt != 0 {
		log.Println("Get first packet error:", reply)
		return
	}

	//First sync the color to the client
	id := reply.id
	err = sync_color_to_client(ws, &id)
	if err != nil {
		log.Println("First sync fail:", err)
		return
	}

	/* Add client to clients.  */
	sync_ch := make(chan bool, 1)
	clients_lock.Lock()
	clients.PushBack(sync_ch)
	clients_lock.Unlock()

	/* Creat new goroutine to receive value.  */
	quit_ch := make(chan bool, 1)
	go handle_receive_pack(ws, quit_ch)

	loop := true
	for loop {
		select {
			case <-sync_ch:
				err = sync_color_to_client(ws, &id)
				if err != nil {
					log.Println("Sync fail:", err)
					loop = false
				}
			case <-quit_ch:
				loop = false
		}
	}

	ws.Close()
}

func main() {
	color.Id = 1
	for x := 0; x < color_x; x++ {
		for y := 0; y < color_x; y++ {
			color.Data[x][y] = "#000000"
		}
	}
	clients = list.New()

	http.Handle("/", websocket.Handler(onConnected))

	log.Fatal(http.ListenAndServeTLS(":443", ssl_crt, ssl_key, nil))
}
