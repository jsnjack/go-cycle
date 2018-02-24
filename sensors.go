package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"

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

// SpeedSensorData is data from the sensor
type SpeedSensorData struct {
	Revolutions uint32
	EventTime   uint16
}

// HandleHRData handles HR data from the HR sensor
func HandleHRData(p gatt.Peripheral) {
	logger := log.New(os.Stdout, "HR ", log.Lmicroseconds|log.Lshortfile)
	defer p.Device().CancelConnection(p)
	service, err := GetService(p, gatt.UUID16(0x180d))
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Printf("Service %s found\n", service.Name())

	ch, err := GetCharacteristic(p, service, gatt.UUID16(0x2a37))
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Printf("Characteristic %s found\n", ch.Name())

	p.DiscoverDescriptors(nil, ch)

	resultCh := make(chan string)
	p.SetNotifyValue(ch, func(ch *gatt.Characteristic, data []byte, err error) {
		if err != nil {
			resultCh <- err.Error()
		}
		heartRate := binary.LittleEndian.Uint16(append([]byte(data[1:2]), []byte{0}...))
		logger.Printf("BPM: %d\n", heartRate)
	})
	<-resultCh
}

// HandleSpeedData handles speed data from the Speed sensor
func HandleSpeedData(p gatt.Peripheral) {
	logger := log.New(os.Stdout, "SP ", log.Lmicroseconds|log.Lshortfile)
	defer p.Device().CancelConnection(p)
	service, err := GetService(p, gatt.UUID16(0x1816))
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Printf("Service %s found\n", service.Name())

	ch, err := GetCharacteristic(p, service, gatt.UUID16(0x2A5B))
	if err != nil {
		logger.Println(err)
		return
	}
	logger.Printf("Characteristic %s found\n", ch.Name())

	p.DiscoverDescriptors(nil, ch)

	resultCh := make(chan string)
	values := make([]SpeedSensorData, 2)
	p.SetNotifyValue(ch, func(ch *gatt.Characteristic, data []byte, err error) {
		if err != nil {
			resultCh <- err.Error()
		}
		offset := 1

		revolutions := binary.LittleEndian.Uint32(append([]byte(data[offset:])))
		offset += 4
		eventTime := binary.LittleEndian.Uint16(append([]byte(data[offset:])))
		values = values[1:]
		values = append(values, SpeedSensorData{Revolutions: revolutions, EventTime: eventTime})
		var time uint16
		if values[1].EventTime >= values[0].EventTime {
			time = values[1].EventTime - values[0].EventTime
		} else {
			time = 65535 - values[0].EventTime + values[1].EventTime + 1
		}
		rps := float64(values[1].Revolutions-values[0].Revolutions) / (float64(time) * 1024)
		speed := rps * math.Pi * (622 + 28*2) * 1000 * 3.6
		if !math.IsInf(speed, 0) {
			if math.IsNaN(speed) {
				speed = 0
			}
			fmt.Printf("Speed: %f km/h\n", speed)
		}
	})
	<-resultCh
}

// GetPeripheralType returns type of the Peripheral
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
			return SpeedPeripheral, nil
		}
	}
	return 0, fmt.Errorf("Unknown device")
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

// IsInterestingPeripheral returns true if peripheral is probably HR or Speed sensor
func IsInterestingPeripheral(id string) bool {
	for _, item := range IgnoredDevices {
		if item == id {
			return false
		}
	}
	return true
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
