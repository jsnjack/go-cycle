package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

var debugFlag *bool

// IgnoredDevices is a list of not interesting devices
var IgnoredDevices []string

// Logger is the main logger
var Logger *log.Logger

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
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		Logger.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	// Register bluetooth handlers
	d.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)
	d.Init(onStateChanged)

	select {}
}
