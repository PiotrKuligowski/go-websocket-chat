package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var wsClients = make(map[*websocket.Conn]bool)

func handleWebsocketClients(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer ws.Close()

	wsClients[ws] = true

	for {
		_, receivedMessage, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(wsClients, ws)
			break
		}

		msg := string(receivedMessage)
		fmt.Println(msg)

		for client := range wsClients {
			if err := client.WriteMessage(websocket.TextMessage, receivedMessage); err != nil {
				fmt.Println("Failed to broadcast message to one of connected clients")
				client.Close()
				delete(wsClients, client)
			}
		}

		fmt.Println("Iteration complete")
	}
}

func main() {
    fmt.Println("Starting websocket chat...")

	http.HandleFunc("/", handleWebsocketClients)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err)
	}
}