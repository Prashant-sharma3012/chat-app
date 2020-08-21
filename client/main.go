package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// !=

func reader(conn *websocket.Conn) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			log.Panic(err)
		}

		fmt.Println(messageType)
		fmt.Println(string(msg))
	}
}

func sender(conn *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		conn.WriteMessage(2, []byte(scanner.Text()))
	}

	if scanner.Err() != nil {
		// handle error.
		log.Fatal(scanner.Err())
	}
}

func main() {

	dialer := &websocket.Dialer{}

	conn, res, err := dialer.Dial("ws://localhost:3000/ws", http.Header{})
	if err != nil {
		fmt.Println(err)
		log.Panic()
	}

	go sender(conn)

	fmt.Println(res)
	conn.WriteMessage(2, []byte("Hello"))
	reader(conn)
}
