package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Panic(err)
		}

		fmt.Println(messageType)
		fmt.Println(string(msg))

		conn.WriteMessage(2, msg)
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	reader(conn)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All Good"))
}

func main() {

	http.HandleFunc("/", healthCheck)
	http.HandleFunc("/ws", socketHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
