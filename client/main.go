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

var supportedCommands = map[string]string{
	"Quit":                         "_Q_",
	"Connect":                      "_C_ <connection_id>",
	"Quit Chat":                    "_QC_",
	"My unique ID to start a chat": "_ID_",
}

var userid string
var messagePrefixOut = "> you: "
var messagePrefixIn = "> "
var quitMessage = "Good Bye :)"
var chatQuitMessage = "Current Chat Connection closed."
var id string

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

		if msg == "_Q_" {
			fmt.Print(messagePrefixIn)
			fmt.Println(quitMessage)
			conn.Close()
			errChan <- "Stop"
			return
		}

		if msg == "_QC_" {
			fmt.Print(messagePrefixIn)
			fmt.Println(chatQuitMessage)
			conn.WriteMessage(2, []byte(msg))
		} else if msg == "_ID_" {
			fmt.Print(messagePrefixIn)
			fmt.Print("ID - ")
			fmt.Println(id)
		} else {
			conn.WriteMessage(2, []byte(msg))
		}
	}

	if scanner.Err() != nil {
		// handle error.
		fmt.Println(scanner.Err())
		errChan <- "Stop"
	}
}

func successMessage() {
	fmt.Println("##################Client is Up and running#################")
	fmt.Println("You can Use following Commands:")
	for k, v := range supportedCommands {
		fmt.Println(k, " : ", v)
	}
	fmt.Println("###########################################################")
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
