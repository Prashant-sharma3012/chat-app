package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type connection struct {
	userid string
	conn   *websocket.Conn
	send   chan string
}

type connections struct {
	numOfConnections int64
	conns            []*connection
}

var clientsConnected connections

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(c *connection) {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Panic(err)
		}

		message := c.userid + " : " + string(msg)
		c.conn.WriteMessage(2, []byte(message))
	}
}

func newSocketConnection(w http.ResponseWriter, r *http.Request) {
	userid := uuid.New().String()

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	conn.WriteMessage(5, []byte(userid))

	c := &connection{
		userid: userid,
		conn:   conn,
		send:   make(chan string),
	}

	clientsConnected.numOfConnections++
	clientsConnected.conns = append(clientsConnected.conns, c)

	fmt.Println("New user connected: " + userid)

	go reader(c)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All Good"))
}

func main() {
	http.HandleFunc("/", healthCheck)
	http.HandleFunc("/ws", newSocketConnection)

	fmt.Println("Server is Up and running on 3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
