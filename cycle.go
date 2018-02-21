package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

// List of services https://www.bluetooth.com/specifications/gatt/services

// MOOVHR is the ID of a pecific MOOH HR device
const MOOVHR = "CC:78:AB:26:B2:73"

// SPEEDSENSOR is the ID of powertap spd 53292
const SPEEDSENSOR = "E7:C6:C3:FC:FC:97"

var connected int

var wg sync.WaitGroup

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("Scanning...")
		d.Scan([]gatt.UUID{}, false)
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	// Move HR
	// https://www.bluetooth.com/specifications/gatt/viewer?attributeXmlFile=org.bluetooth.service.heart_rate.xml
	switch p.ID() {
	case MOOVHR, SPEEDSENSOR:
		fmt.Printf("Found %s\n", p.Name())
		p.Device().Connect(p)
		connected = connected + 1
	default:
		return
	}

	// Stop scanning once we've got the peripheral we're looking for.
	if connected == 2 {
		p.Device().StopScanning()
	}
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	fmt.Printf("Connected %s\n", p.Name())

	if err := p.SetMTU(500); err != nil {
		fmt.Printf("Failed to set MTU, err: %s\n", err)
	}

	switch p.ID() {
	case MOOVHR:
		go getHRData(p)
	case SPEEDSENSOR:
		go getSpeedData(p)
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
		logger.Printf("BPM: %d: data %v\n", heartRate, data)
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
	p.SetNotifyValue(ch, func(ch *gatt.Characteristic, data []byte, err error) {
		if err != nil {
			resultCh <- err.Error()
		}
		offset := 1

		revolutions := binary.LittleEndian.Uint32(append([]byte(data[offset:])))
		offset += 4
		eventTime := binary.LittleEndian.Uint16(append([]byte(data[offset:])))
		logger.Printf("Revolutions: %d, Time: %fs: data %v\n", revolutions, float32(eventTime)/1024, data)
	})
	<-resultCh
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	wg.Done()
	fmt.Printf("Disconnected %s\n", p.Name())
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
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers.
	d.Handle(
		gatt.PeripheralDiscovered(onPeriphDiscovered),
		gatt.PeripheralConnected(onPeriphConnected),
		gatt.PeripheralDisconnected(onPeriphDisconnected),
	)

	wg.Add(2)
	d.Init(onStateChanged)
	wg.Wait()
	fmt.Println("Done")
}
