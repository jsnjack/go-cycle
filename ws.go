package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/muka/go-bluetooth/api"

	"github.com/gorilla/websocket"
)

// List of connected clients
var clients []*websocket.Conn

// BroadcastChannel contains messages to send to all connected clients
var BroadcastChannel = make(chan []byte)

// SendingSynth shows if synth data is being sent
var SendingSynth = false

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

	// Send artificial sensor's events
	if *debugFlag && !SendingSynth {
		SendingSynth = true
		go func() {
			for {
				SendSynthCSCEvent()
				SendSynthHREvent()
				time.Sleep(1 * time.Second)
			}
		}()
	}

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
			for _, d := range DiscoveredDevices {
				data := DeviceDiscoveredData{Name: d.Name, ID: d.Address}
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
			err = api.StopDiscovery()
			if err != nil {
				Logger.Println(err)
			} else {
				Logger.Println("Stop scanning.")
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
