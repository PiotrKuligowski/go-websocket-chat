package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatMessage struct {
	ChatId   string `json:"ChatId"`
	PlayerId string `json:"PlayerId"`
	Data     string `json:"Data"`
}

var rooms = make(map[string]map[string]*websocket.Conn)

func handleWebsocketClients(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		_, receivedMessage, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}

		var message ChatMessage
		json.Unmarshal(receivedMessage, &message)

		if rooms[message.ChatId] == nil {
			rooms[message.ChatId] = make(map[string]*websocket.Conn)
		}

		rooms[message.ChatId][message.PlayerId] = conn

		for playerId, client := range rooms[message.ChatId] {
			if err := client.WriteMessage(websocket.TextMessage, []byte(message.Data)); err != nil {
				fmt.Println("Failed to broadcast message to player: ", message.PlayerId)
				client.Close()
				delete(rooms[message.ChatId], playerId)
			}
			fmt.Println("Message from chat:", message.ChatId, "broadcasted to client:", playerId)
		}
	}
}

func main() {
	fmt.Println("Starting websocket chat...")

	http.HandleFunc("/ws", handleWebsocketClients)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err)
	}
}
