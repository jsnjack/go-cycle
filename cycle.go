package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
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
		IgnoredDevices = append(IgnoredDevices, p.ID())
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
			go getHRData(p)
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
			go getSpeedData(p)
		} else {
			logger.Println("Speed sensor already connected")
		}
	default:
		p.Device().CancelConnection(p)
		IgnoredDevices = append(IgnoredDevices, p.ID())
		logger.Printf("Ignoring device %s", p.Name())
	}

	if Connected.AllConnected() {
		logger.Println("All devices connected. Stop scanning")
		p.Device().StopScanning()
	}
}

func getHRData(p gatt.Peripheral) {
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

func getSpeedData(p gatt.Peripheral) {
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
		rps := float32(values[1].Revolutions-values[0].Revolutions) / (float32(time) * 1024)
		speed := rps * math.Pi * (622 + 28*2) * 1000 * 3.6
		fmt.Printf("Speed: %f km/h\n", speed)
	})
	<-resultCh
}

// SpeedSensorData is data from the sensor
type SpeedSensorData struct {
	Revolutions uint32
	EventTime   uint16
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
		return
	}
	Logger.Println("Scanning for device to reconnect...")
	p.Device().Scan([]gatt.UUID{}, false)
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

func main() {
	Logger.Println("Starting...")
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		Logger.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)

	mainLoop := make(chan bool)
	d.Init(onStateChanged)
	<-mainLoop
	Logger.Println("Done")
}
