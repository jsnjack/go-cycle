package main

import (
	"encoding/json"
	"fmt"

	"github.com/paypal/gatt"
)

// List of services https://www.bluetooth.com/specifications/gatt/services

// PeripheralType contains type of the Peripheral
type PeripheralType int

// HRPeripheral is a Heart Rate monitor device
// https://www.bluetooth.com/specifications/gatt/viewer?attributeXmlFile=org.bluetooth.service.heart_rate.xml
var HRPeripheral PeripheralType = 1

// CSCPeripheral is a Speed and cadence device
// https://www.bluetooth.com/specifications/gatt/viewer?attributeXmlFile=org.bluetooth.service.cycling_speed_and_cadence.xml
var CSCPeripheral PeripheralType = 2

// DiscoveredDevices is a list of discovered Peripheral
var DiscoveredDevices []gatt.Peripheral

// ISensors an interface for sensors
type ISensors interface {
	Listen()
	GetType() PeripheralType
}

// RequestedDevices contains information about connected devices
type RequestedDevices struct {
	Peripherals []ISensors
}

// IsHRSelected ...
func (r *RequestedDevices) IsHRSelected() bool {
	for _, item := range r.Peripherals {
		if item.GetType() == HRPeripheral {
			return true
		}
	}
	return false
}

// IsCSCSelected ...
func (r *RequestedDevices) IsCSCSelected() bool {
	for _, item := range r.Peripherals {
		if item.GetType() == CSCPeripheral {
			return true
		}
	}
	return false
}

// ActiveDevices ...
var ActiveDevices RequestedDevices

// ConnectToDevice connects to the device with specified ID
func ConnectToDevice(id string) {
	for _, p := range DiscoveredDevices {
		if p.ID() == id {
			Logger.Println("Connecting to device", id)
			p.Device().Connect(p)
		}
	}
}

// GetService returns service with specified name
func GetService(p gatt.Peripheral, uuid gatt.UUID) (*gatt.Service, error) {
	services, err := p.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return nil, err
	}
	for _, item := range services {
		if item.UUID().Equal(uuid) {
			return item, nil
		}
	}
	return nil, fmt.Errorf("Service %s not found", uuid.String())
}

// GetCharacteristic returns characteristics with specified name
func GetCharacteristic(p gatt.Peripheral, service *gatt.Service, uuid gatt.UUID) (*gatt.Characteristic, error) {
	chs, err := p.DiscoverCharacteristics(nil, service)
	if err != nil {
		fmt.Printf("Failed to discover characteristics, err: %s\n", err)
		return nil, err
	}
	for _, item := range chs {
		if item.UUID().Equal(uuid) {
			return item, nil
		}
	}
	return nil, fmt.Errorf("Characteristic %s not found", uuid.String())
}

// GetActiveDeviceType returns type of active device for status message
func GetActiveDeviceType(id string) string {
	for _, item := range ActiveDevices.Peripherals {
		switch item.GetType() {
		case HRPeripheral:
			return "hr"
		case CSCPeripheral:
			return "csc"
		}
	}
	return ""
}

// GetPeripheralType returns type of the device - HR or CSC
func GetPeripheralType(p gatt.Peripheral) (PeripheralType, error) {
	Logger.Println("Discovering services")
	services, err := p.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Peripheral %s: Failed to discover services, err: %s\n", p.Name(), err)
		return 0, err
	}
	for _, item := range services {
		if item.UUID().Equal(HRServiceUUID) {
			return HRPeripheral, nil
		} else if item.UUID().Equal(SpeedServiceUUID) {
			return CSCPeripheral, nil
		}
	}
	return 0, fmt.Errorf("Unknown device")
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	Logger.Printf("Discovered %s %s\n", p.Name(), p.ID())
	if p.Name() != "" {
		DiscoveredDevices = append(DiscoveredDevices, p)
		data := DeviceDiscoveredData{Name: p.Name(), ID: p.ID()}
		msg := WSMessage{Type: "ws.device:discovered", Data: data}
		msgByte, _ := json.Marshal(&msg)
		BroadcastChannel <- msgByte
	}
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	Logger.Printf("Connected %s\n", p.ID())

	pType, err := GetPeripheralType(p)
	if err != nil {
		Logger.Println(err.Error())
		p.Device().CancelConnection(p)
		return
	}
	switch pType {
	case HRPeripheral:
		if !ActiveDevices.IsHRSelected() {
			hrsensor := HRSensor{Peripheral: p}
			ActiveDevices.Peripherals = append(ActiveDevices.Peripherals, &hrsensor)
			go hrsensor.Listen()
		} else {
			p.Device().CancelConnection(p)
			// TODO: Remove old sensor
			Logger.Println("HR sensor already connected")
		}
	case CSCPeripheral:
		if !ActiveDevices.IsCSCSelected() {
			Logger.Printf("Found %s\n", p.Name())
			cscsensor := CSCSensor{Peripheral: p}
			ActiveDevices.Peripherals = append(ActiveDevices.Peripherals, &cscsensor)
			go cscsensor.Listen()
		} else {
			p.Device().CancelConnection(p)
			// TODO: Remove old sensor
			Logger.Println("Speed sensor already connected")
		}
	default:
		p.Device().CancelConnection(p)
		return
	}

	msgStatus := DeviceStatusData{ID: p.ID(), Status: "connected", RecognizedAs: GetActiveDeviceType(p.ID())}
	wsMsgStatus := WSMessage{Type: "ws.device:status", Data: msgStatus}
	msgB, err := json.Marshal(&wsMsgStatus)
	if err != nil {
		Logger.Println(err)
	} else {
		BroadcastChannel <- msgB
	}
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	fmt.Printf("Disconnected %s\n", p.Name())

	msgStatus := DeviceStatusData{ID: p.ID(), Status: "disconnected", RecognizedAs: GetActiveDeviceType(p.ID())}
	wsMsgStatus := WSMessage{Type: "ws.device:status", Data: msgStatus}
	msgB, err := json.Marshal(&wsMsgStatus)
	if err != nil {
		Logger.Println(err)
	} else {
		BroadcastChannel <- msgB
	}
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
