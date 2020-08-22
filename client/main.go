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

var userid string
var messagePrefixOut = "> you: "
var messagePrefixIn = "> "
var quitMessage = "Good Bye :)"

func reader(conn *websocket.Conn, errChan chan<- string) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			errChan <- "Stop"
		}

		fmt.Print(messagePrefixIn)
		fmt.Println(string(msg))
		fmt.Print(messagePrefixOut)
	}
}

func sender(conn *websocket.Conn, errChan chan<- string) {
	fmt.Print(messagePrefixOut)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()

		if msg == "Quit" {
			fmt.Print(messagePrefixIn)
			fmt.Println(quitMessage)
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

func successMessage() {
	fmt.Println("Client is Up and running, Type: 'Quit' to stop messaging")
}

func main() {

	dialer := &websocket.Dialer{}

	conn, _, err := dialer.Dial("ws://localhost:3000/ws", http.Header{})
	if err != nil {
		fmt.Println(err)
		log.Panic()
	}

	// on successful connect write options and meesage to console
	successMessage()
	senderOrReceiverError := make(chan string)

	go sender(conn, senderOrReceiverError)
	go reader(conn, senderOrReceiverError)

	<-senderOrReceiverError
}
