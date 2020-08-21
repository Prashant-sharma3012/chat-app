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

func reader(conn *websocket.Conn, errChan chan<- string) {
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			errChan <- "Stop"
		}

		fmt.Println(messageType)
		fmt.Println(string(msg))
	}
}

func sender(conn *websocket.Conn, errChan chan<- string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()

		if msg == "Quit" {
			conn.Close()
			errChan <- "Stop"
			return
		}

		conn.WriteMessage(2, []byte(msg))
	}

	if scanner.Err() != nil {
		// handle error.
		fmt.Println(scanner.Err())
		errChan <- "Stop"
	}
}

func main() {

	dialer := &websocket.Dialer{}

	conn, _, err := dialer.Dial("ws://localhost:3000/ws", http.Header{})
	if err != nil {
		fmt.Println(err)
		log.Panic()
	}

	senderOrReceiverError := make(chan string)

	go sender(conn, senderOrReceiverError)
	go reader(conn, senderOrReceiverError)

	fmt.Println("Client is Up and running, Type: 'Quit' to stop messaging ")

	<-senderOrReceiverError
}
