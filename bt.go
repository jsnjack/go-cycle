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

// ActiveDevices ...
type ActiveDevices struct {
	HR  HRSensor
	CSC CSCSensor
}

// ConnectedDevices ...
var ConnectedDevices ActiveDevices

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
	if ConnectedDevices.HR.Initialized && ConnectedDevices.HR.GetPeripheral().ID() == id {
		return "hr"
	} else if ConnectedDevices.CSC.Initialized && ConnectedDevices.CSC.GetPeripheral().ID() == id {
		return "csc"
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
		hrsensor := HRSensor{Peripheral: p, Initialized: true}
		ConnectedDevices.HR = hrsensor
		go hrsensor.Listen()
	case CSCPeripheral:
		cscsensor := CSCSensor{Peripheral: p, Initialized: true}
		ConnectedDevices.CSC = cscsensor
		go cscsensor.Listen()
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
	if ConnectedDevices.HR.Initialized && ConnectedDevices.HR.GetPeripheral().ID() == p.ID() ||
		ConnectedDevices.CSC.Initialized && ConnectedDevices.CSC.GetPeripheral().ID() == p.ID() {
		Logger.Println("Reconnecting", p.ID())
		p.Device().Connect(p)
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
