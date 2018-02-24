package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// List of connected clients
var clients []*websocket.Conn

// BroadcastChannel contains messages to send to all connected clients
var BroadcastChannel = make(chan []byte)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// BroadcastMessages broadcasts messages to all connected clients
func BroadcastMessages() {
	for {
		msg := <-BroadcastChannel
		Logger.Println("Message arrived")
		for _, client := range clients {
			Logger.Println("Delivering message")
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				Logger.Printf("error: %v\n", err)
				client.Close()
				removeConnection(client)
			}
		}
	}
}

func handleWSConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	Logger.Println("Establishing connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	Logger.Println("Establishing connection 2")

	if err != nil {
		Logger.Fatal(err)
	}
	// Register our new client
	clients = append(clients, ws)

}

func removeConnection(conn *websocket.Conn) {
	Logger.Println("Closing connection")
	var updatedClients []*websocket.Conn
	for _, item := range clients {
		if item != conn {
			updatedClients = append(updatedClients, item)
		}
	}
	clients = updatedClients
}
