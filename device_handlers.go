package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/muka/go-bluetooth/api"
)

var mutex = &sync.Mutex{}

// ConnectedDevices list of connected devices
var ConnectedDevices []string

// GetCharsAttempts is amount of attempts to get available characteristics
const GetCharsAttempts = 10

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
		if attempt < GetCharsAttempts {
			deviceLogger.Println("Device is not ready. Get characteristics later")
			go RetryForChars(device, attempt)
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

// RetryForChars try to get all chars again
func RetryForChars(device *api.Device, attempt int) {
	attempt = attempt + 1
	device.Disconnect()
	time.Sleep(5 * time.Second)
	device.Connect()
	DeviceConnectedHandler(device, attempt)
}

// ConnectToDevice connects to the device with specified ID
func ConnectToDevice(address string) {
	if IsConnectedDevice(address) {
		Logger.Printf("Device %s already connected\n", address)
		return
	}
	device, err := api.GetDeviceByAddress(address)
	if err != nil {
		Logger.Println(err)
		return
	}
	if device.IsConnected() {
		Logger.Printf("Device %s already connected\n", address)
	} else {
		Logger.Printf("Connecting...%s\n", address)
		err = device.Connect()
		if err != nil {
			Logger.Println(err)
		}
	}
	DeviceConnectedHandler(device, 0)
}

// Reconnect to device
func Reconnect(address string) {
	Logger.Printf("Reconnecting to device in 3 %s", address)
	time.Sleep(3 * time.Second)
	manager, err := api.GetManager()
	if err != nil {
		Logger.Println(err)
	}
	Logger.Println("Refreshing state...")
	err = manager.RefreshState()
	if err != nil {
		Logger.Println(err)
	}
	ConnectToDevice(address)
	// Keep trying until succeded
	go func() {
		time.Sleep(5 * time.Second)
		if !IsConnectedDevice(address) {
			Reconnect(address)
		}
	}()
}

//IsConnectedDevice returns if device is connected
func IsConnectedDevice(address string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	for _, cd := range ConnectedDevices {
		if cd == address {
			return true
		}
	}
	return false
}
