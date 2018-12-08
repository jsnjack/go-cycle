package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/muka/go-bluetooth/api"
)

//ConnectedDevices list of connected devices
var ConnectedDevices []string

// DiscoveredDevice ...
type DiscoveredDevice struct {
	Name    string
	Address string // aka ID
}

// DiscoveredDevices is a list of discovered devices
var DiscoveredDevices []DiscoveredDevice

// DeviceDiscoveredHandler handles discovery of the device
func DeviceDiscoveredHandler(device *api.Device) {
	props, err := device.GetProperties()
	if err != nil {
		Logger.Println(err)
		return
	}
	Logger.Printf("Discovered %s %s\n", props.Name, props.Address)
	if props.Name != "" {
		DiscoveredDevices = append(
			DiscoveredDevices,
			DiscoveredDevice{Name: props.Name, Address: props.Address},
		)
		data := DeviceDiscoveredData{Name: props.Name, ID: props.Address}
		msg := WSMessage{Type: "ws.device:discovered", Data: data}
		msgByte, _ := json.Marshal(&msg)
		BroadcastChannel <- msgByte
	}
}

// DeviceConnectedHandler handles discovery of the device
func DeviceConnectedHandler(device *api.Device, attempt int) {
	address, err := device.GetProperty("Address")
	if err != nil {
		Logger.Println(err)
		return
	}

	addressStr := fmt.Sprintf("%s", address)

	deviceLogger := log.New(os.Stdout, "["+addressStr+"] ", log.Lshortfile)

	deviceLogger.Printf("Device connected")

	// Get all characteristics
	charList, err := device.GetCharsList()
	if err != nil {
		deviceLogger.Println(err)
		return
	}
	deviceLogger.Printf("Discovered %d characteristics", len(charList))
	if len(charList) == 0 {
		if attempt < 10 {
			deviceLogger.Println("Device is not ready. Get characteristics later")
			go Reconnect(device, attempt)
			return
		}
	}

	// Look for known characteristic
	for _, charPath := range charList {
		char, err := device.GetChar(fmt.Sprintf("%s", charPath))
		if err != nil {
			deviceLogger.Println(err)
			continue
		}
		uuid, err := char.GetProperty("UUID")
		if err != nil {
			deviceLogger.Println(err)
			continue
		}
		deviceLogger.Printf("  char: %s", uuid)
		switch fmt.Sprintf("%s", uuid)[4:8] {
		case HRMeasuremetCharID:
			deviceLogger.Println("is HR device")
			sensor := Sensor{
				Char:    char,
				Kind:    HRKind,
				Logger:  deviceLogger,
				Address: addressStr,
				Device:  device,
			}
			ConnectedDevices = append(ConnectedDevices, addressStr)
			go sensor.Listen()
			return
		case CSCMeasuremetCharID:
			deviceLogger.Println("is CSC device")
			sensor := Sensor{
				Char:    char,
				Logger:  deviceLogger,
				Address: addressStr,
				Device:  device,
			}
			ConnectedDevices = append(ConnectedDevices, addressStr)
			go sensor.Listen()
			return

		}
	}
	deviceLogger.Println("Device with unknown characteristics")
	err = device.Disconnect()
	if err != nil {
		deviceLogger.Println(err)
	}
}

// Reconnect to device (normally when unable to get characteristics)
func Reconnect(device *api.Device, attempt int) {
	attempt = attempt + 1
	device.Disconnect()
	time.Sleep(5 * time.Second)
	device.Connect()
	DeviceConnectedHandler(device, attempt)
}

// ConnectToDevice connects to the device with specified ID
func ConnectToDevice(address string) {
	for _, cd := range ConnectedDevices {
		if cd == address {
			Logger.Printf("Device %s already connected\n", address)
			return
		}
	}
	device, err := api.GetDeviceByAddress(address)
	if err != nil {
		Logger.Println(err)
		return
	}
	if device.IsConnected() {
		Logger.Printf("Device %s already connected\n", address)
	} else {
		err = device.Connect()
		if err != nil {
			Logger.Println(err)
		}
	}
	DeviceConnectedHandler(device, 0)
}
