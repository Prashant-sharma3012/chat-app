package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var supportedCommands = map[string]string{
	"Quit":    "_Q_",
	"Connect": "Connect <connection_id>",
}

var usersChatMap = map[string]string{}
var userConnMap = map[string]*websocket.Conn{}

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

		// its a connect message
		if parts := strings.Split(string(msg), " "); parts[0] == "_C_" {
			userid := parts[1]
			usersChatMap[userid] = c.userid
			usersChatMap[c.userid] = userid
			userConnMap[c.userid] = c.conn

			for _, conn := range clientsConnected.conns {
				if conn.userid == userid {
					userConnMap[userid] = conn.conn
				}
			}

			c.conn.WriteMessage(2, []byte("Connection successful"))

		} else {
			// check if user has an entry in userchatmap
			// route messages tothat user
			if sendTo, ok := usersChatMap[c.userid]; ok {
				userConnMap[sendTo].WriteMessage(2, msg)
			} else {
				message := c.userid + ": " + string(msg)
				c.conn.WriteMessage(2, []byte(message))
			}
		}
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
