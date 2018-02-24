package main

import (
	"log"
	"os"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

// Connected contains information about connected devices
var Connected ConnectedDevices

// IgnoredDevices is a list of not interesting devices
var IgnoredDevices []string

// Logger is the main logger
var Logger *log.Logger

func init() {
	Logger = log.New(os.Stdout, "", log.Lmicroseconds|log.Lshortfile)
}

func main() {
	Logger.Println("Starting...")

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
