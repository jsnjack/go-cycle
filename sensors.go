package main

import (
	"fmt"

	"github.com/paypal/gatt"
)

// List of services https://www.bluetooth.com/specifications/gatt/services

// MOOVHR is the ID of a pecific MOOH HR device
const MOOVHR = "CC:78:AB:26:B2:73"

// SPEEDSENSOR is the ID of powertap spd 53292
const SPEEDSENSOR = "E7:C6:C3:FC:FC:97"

// HRServiceUUID is UUID for heart_rate service
var HRServiceUUID = gatt.UUID16(0x180d)

// SpeedServiceUUID is UUID for cycling_speed_and_cadence service
var SpeedServiceUUID = gatt.UUID16(0x1816)

// ConnectedDevices contains information about connected devices
type ConnectedDevices struct {
	HRSensor      bool
	SpeedSensor   bool
	HRSensorID    string
	SpeedSensorID string
}

// AllConnected returns if all devices were connected
func (c *ConnectedDevices) AllConnected() bool {
	return c.HRSensor && c.SpeedSensor
}

// PeripheralType contains type of the Peripheral
type PeripheralType int

// HRPeripheral is a Heart Rate monitor device
// https://www.bluetooth.com/specifications/gatt/viewer?attributeXmlFile=org.bluetooth.service.heart_rate.xml
var HRPeripheral PeripheralType = 1

// SpeedPeripheral is a Speed and cadence device
// https://www.bluetooth.com/specifications/gatt/viewer?attributeXmlFile=org.bluetooth.service.cycling_speed_and_cadence.xml
var SpeedPeripheral PeripheralType = 2

// GetPeripheralType returns type of the Peripheral
func GetPeripheralType(p gatt.Peripheral) (PeripheralType, error) {
	Logger.Println("Discovering services")
	services, err := p.DiscoverServices(nil)
	Logger.Println(services)
	if err != nil {
		fmt.Printf("Peripheral %s: Failed to discover services, err: %s\n", p.Name(), err)
		return 0, err
	}
	for _, item := range services {
		if item.UUID().Equal(HRServiceUUID) {
			return HRPeripheral, nil
		} else if item.UUID().Equal(SpeedServiceUUID) {
			return SpeedPeripheral, nil
		}
	}
	return 0, fmt.Errorf("Unknown device")
}

// IsInterestingPeripheral returns true if peripheral is probably HR or Speed sensor
func IsInterestingPeripheral(id string) bool {
	for _, item := range IgnoredDevices {
		if item == id {
			return false
		}
	}
	return true
}
