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
	"Quit":      "_Q_",
	"Connect":   "Connect <connection_id>",
	"Quit Chat": "_QC_ <connection_id>",
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

func removeFromClientsConnected(pos int) []*connection {
	clientsConnected.conns[pos] = clientsConnected.conns[len(clientsConnected.conns)-1]
	return clientsConnected.conns[:len(clientsConnected.conns)-1]
}

func deleteAllEntries(c *connection) {
	// remove from connMap
	delete(userConnMap, c.userid)

	// remove from chat map
	for k, v := range usersChatMap {
		if k == c.userid || v == c.userid {
			delete(usersChatMap, k)
		}
	}

	// remove from connected clients
	indexToRemove := 0
	for indx, v := range clientsConnected.conns {
		if v.userid == c.userid {
			indexToRemove = indx
		}
	}

	clientsConnected.conns = removeFromClientsConnected(indexToRemove)
	clientsConnected.numOfConnections--
}

func reader(c *connection) {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				fmt.Println("User Disconnected: " + c.userid)
				// remove the conn from map and also form connections array
				deleteAllEntries(c)
			}
			fmt.Println(err)
			return
		}

		parts := strings.Split(string(msg), " ")

		if string(msg) == "_ID_" {
			c.conn.WriteMessage(2, []byte(c.userid))
		} else if parts[0] == "_C_" {
			userid := parts[1]
			usersChatMap[userid] = c.userid
			usersChatMap[c.userid] = userid
			c.conn.WriteMessage(2, []byte("Connection successful"))
		} else if parts[0] == "_QC_" {
			userToNotify := usersChatMap[c.userid]
			closeMessage := c.userid + " Closed the current chat"
			delete(usersChatMap, c.userid)
			delete(usersChatMap, userToNotify)
			userConnMap[userToNotify].WriteMessage(2, []byte(closeMessage))
		} else if sendTo, ok := usersChatMap[c.userid]; ok {
			userConnMap[sendTo].WriteMessage(2, msg)
		} else {
			message := c.userid + ": " + string(msg)
			c.conn.WriteMessage(2, []byte(message))
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
	userConnMap[userid] = conn

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
