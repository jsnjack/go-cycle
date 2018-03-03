package main

// WSMessage is a message to deliver over websocket to the electorn app
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// DeviceDiscoveredData ...
type DeviceDiscoveredData struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// ConnectToDeviceData ...
type ConnectToDeviceData struct {
	ID string `json:"id"`
}

// DeviceStatusData ...
type DeviceStatusData struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	RecognizedAs string `json:"recognizedAs"`
}
