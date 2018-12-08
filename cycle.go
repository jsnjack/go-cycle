package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/emitter"
)

var debugFlag *bool

// IgnoredDevices is a list of not interesting devices
var IgnoredDevices []string

// Logger is the main logger
var Logger *log.Logger

const adapterID = "hci0"

func init() {
	debugFlag = flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	Logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
}

func main() {
	go func() {
		// Websocket section
		http.HandleFunc("/ws", handleWSConnection)

		go BroadcastMessages()

		Logger.Println("Starting ws server on port 8000...")
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	// Bluetooth section

	// Start discovery on default device
	api.StartDiscovery()

	Logger.Println("Scanning...")
	api.On("discovery", emitter.NewCallback(func(ev emitter.Event) {
		Logger.Println(ev)
		discoveryEvent := ev.GetData().(api.DiscoveredDeviceEvent)
		Logger.Println(discoveryEvent, discoveryEvent.Status)
		device := discoveryEvent.Device
		DeviceDiscoveredHandler(device)
	}))

	// Cached devices
	cached, err := api.GetDevices()
	if err != nil {
		Logger.Println(err)
	}
	for _, device := range cached {
		DeviceDiscoveredHandler(&device)
	}

	select {}
}
