package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

// List of connected clients
var clients []*websocket.Conn

// BroadcastChannel contains messages to send to all connected clients
var BroadcastChannel = make(chan []byte)

// IncomingMessage ...
type IncomingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// BroadcastMessages broadcasts messages to all connected clients
func BroadcastMessages() {
	for {
		msg := <-BroadcastChannel
		for _, client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				Logger.Printf("error: %v\n", err)
			}
		}
	}
}

func handleWSConnection(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	Logger.Println("New ws connection")
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		Logger.Fatal(err)
	}
	// Register our new client
	clients = append(clients, ws)

	defer func() {
		ws.Close()
		removeConnection(ws)
	}()

	for {
		var msg IncomingMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			Logger.Println(err)
			removeConnection(ws)
			return
		}
		switch msg.Type {
		case "app.bt:scan":
			for _, p := range DiscoveredDevices {
				data := DeviceDiscoveredData{Name: p.Name(), ID: p.ID()}
				msg := WSMessage{Type: "ws.device:discovered", Data: data}
				msgByte, _ := json.Marshal(&msg)
				BroadcastChannel <- msgByte
			}
		case "app.device:connect":
			var data ConnectToDeviceData
			err = json.Unmarshal(msg.Data, &data)
			if err != nil {
				Logger.Println(err)
			} else {
				ConnectToDevice(data.ID)
			}
		case "app.bt:scan_stop":
			if len(DiscoveredDevices) > 0 {
				DiscoveredDevices[0].Device().StopScanning()
				Logger.Println("Stop scanning")
			} else {
				Logger.Println("No discovered devices")
			}
		default:
			Logger.Println("Unhandled message", msg)
		}
	}
}

func removeConnection(conn *websocket.Conn) {
	var updatedClients []*websocket.Conn
	for _, item := range clients {
		if item != conn {
			updatedClients = append(updatedClients, item)
		}
	}
	clients = updatedClients
}
