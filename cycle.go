package main

import (
	"fmt"
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

func onStateChanged(d gatt.Device, s gatt.State) {
	Logger.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		Logger.Println("Scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		Logger.Println("Stop scanning")
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	Logger.Printf("Discovered %s %s\n", p.Name(), p.ID())
	if IsInterestingPeripheral(p.ID()) {
		if p.Name() != "" {
			Logger.Printf("Connecting to %s...\n", p.ID())
			p.Device().Connect(p)
		} else {
			Logger.Printf("Ignoring %s\n", p.ID())
			IgnoredDevices = append(IgnoredDevices, p.ID())
		}
	}
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	logger := log.New(os.Stdout, fmt.Sprintf("%s ", p.ID()), log.Lmicroseconds|log.Lshortfile)
	logger.Printf("Connected %s\n", p.ID())

	pType, err := GetPeripheralType(p)
	if err != nil {
		logger.Println(err.Error())
		p.Device().CancelConnection(p)
		return
	}
	switch pType {
	case HRPeripheral:
		if !Connected.HRSensor {
			logger.Printf("Found %s\n", p.Name())
			Connected.HRSensorID = p.ID()
			Connected.HRSensor = true
			if err := p.SetMTU(500); err != nil {
				logger.Printf("Failed to set MTU, err: %s\n", err)
			}
			go HandleHRData(p)
		} else {
			logger.Println("HR sensor already connected")
		}
	case SpeedPeripheral:
		if !Connected.SpeedSensor {
			logger.Printf("Found %s\n", p.Name())
			Connected.SpeedSensorID = p.ID()
			Connected.SpeedSensor = true
			if err := p.SetMTU(500); err != nil {
				logger.Printf("Failed to set MTU, err: %s\n", err)
			}
			go HandleSpeedData(p)
		} else {
			logger.Println("Speed sensor already connected")
		}
	default:
		p.Device().CancelConnection(p)
		logger.Printf("Ignoring device %s", p.Name())
	}

	if Connected.AllConnected() {
		logger.Println("All devices connected. Stop scanning")
		p.Device().StopScanning()
	}
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	fmt.Printf("Disconnected %s\n", p.Name())
	switch p.ID() {
	case Connected.HRSensorID:
		Connected.HRSensor = false
	case Connected.SpeedSensorID:
		Connected.SpeedSensor = false
	default:
		Logger.Printf("Unsupported device %s", p.Name())
		IgnoredDevices = append(IgnoredDevices, p.ID())
		return
	}
	Logger.Println("Scanning for device to reconnect...")
	p.Device().Scan([]gatt.UUID{}, false)
}

func main() {
	Logger.Println("Starting...")
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
