package main

import "github.com/paypal/gatt"

// HRServiceUUID is UUID for heart_rate service
var HRServiceUUID = gatt.UUID16(0x180d)

// SpeedServiceUUID is UUID for cycling_speed_and_cadence service
var SpeedServiceUUID = gatt.UUID16(0x1816)

// BatteryServiceUUID is UUID for the battery status service
var BatteryServiceUUID = gatt.UUID16(0x180F)

// BatteryCharacteristicUUID is UUID for the battery level
var BatteryCharacteristicUUID = gatt.UUID16(0x2A19)

// List of services https://www.bluetooth.com/specifications/gatt/services

// PeripheralType contains type of the Peripheral
type PeripheralType int

// HRPeripheral is a Heart Rate monitor device
// https://www.bluetooth.com/specifications/gatt/viewer?attributeXmlFile=org.bluetooth.service.heart_rate.xml
var HRPeripheral PeripheralType = 1

// CSCPeripheral is a Speed and cadence device
// https://www.bluetooth.com/specifications/gatt/viewer?attributeXmlFile=org.bluetooth.service.cycling_speed_and_cadence.xml
var CSCPeripheral PeripheralType = 2

// SensorKind kind of the sensor, depends on returned measurements
type SensorKind string

// HRKind measures heart rate
var HRKind SensorKind = "hr"

// SpeedKind measures speed
var SpeedKind SensorKind = "csc_speed"

// CadenceKind measures speed
var CadenceKind SensorKind = "csc_cadence"

// Sensor is a common struct for all sensors
type Sensor struct {
	Peripheral gatt.Peripheral
	Kind       SensorKind
}

// GetPeripheral ...
func (sensor *Sensor) GetPeripheral() gatt.Peripheral {
	return sensor.Peripheral
}

// GetID ...
func (sensor *Sensor) GetID() string {
	return sensor.Peripheral.ID()
}

// GetKind ...
func (sensor *Sensor) GetKind() SensorKind {
	return sensor.Kind
}

// GetBatteryLevel ...
func (sensor *Sensor) GetBatteryLevel() (int, error) {
	service, err := GetService(sensor.Peripheral, BatteryServiceUUID)
	if err != nil {
		Logger.Println(err)
		return 0, err
	}

	ch, err := GetCharacteristic(sensor.Peripheral, service, BatteryCharacteristicUUID)
	if err != nil {
		Logger.Println(err)
		return 0, err
	}

	datab, err := sensor.Peripheral.ReadCharacteristic(ch)
	return int(datab[0]), nil
}
