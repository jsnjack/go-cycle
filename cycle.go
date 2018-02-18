package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

// MOOVHR is the ID of a pecific MOOH HR device
const MOOVHR = "CC:78:AB:26:B2:73"

// SPEEDSENSOR is the ID of powertap spd 53292
const SPEEDSENSOR = "E7:C6:C3:FC:FC:97"

var done = make(chan struct{})

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
	if strings.ToUpper(p.ID()) != MOOVHR {
		return
	}

	// Stop scanning once we've got the peripheral we're looking for.
	p.Device().StopScanning()

	fmt.Printf("Found %s\n", p.Name())
	p.Device().Connect(p)
}

func onPeriphConnected(p gatt.Peripheral, err error) {
	fmt.Println("Connected")
	defer p.Device().CancelConnection(p)

	if err := p.SetMTU(500); err != nil {
		fmt.Printf("Failed to set MTU, err: %s\n", err)
	}

	service, err := GetService(p, "Heart Rate")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Service %s found\n", service.Name())

	ch, err := GetCharacteristic(p, service, "Heart Rate Measurement")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Characteristic %s found\n", ch.Name())

	p.DiscoverDescriptors(nil, ch)

	resultCh := make(chan string)
	p.SetNotifyValue(ch, func(ch *gatt.Characteristic, data []byte, err error) {
		if err != nil {
			resultCh <- err.Error()
		}
		heartRate := binary.LittleEndian.Uint16(append([]byte(data[1:2]), []byte{0}...))
		fmt.Printf("HR: %d\n", heartRate)
	})
	<-resultCh
}

func onPeriphDisconnected(p gatt.Peripheral, err error) {
	fmt.Println("Disconnected")
	close(done)
}

// GetService returns service with specified name
func GetService(p gatt.Peripheral, name string) (*gatt.Service, error) {
	services, err := p.DiscoverServices(nil)
	if err != nil {
		fmt.Printf("Failed to discover services, err: %s\n", err)
		return nil, err
	}
	for _, item := range services {
		if item.Name() == name {
			return item, nil
		}
	}
	return nil, fmt.Errorf("Service %s not found", name)
}

// GetCharacteristic returns characteristics with specified name
func GetCharacteristic(p gatt.Peripheral, service *gatt.Service, name string) (*gatt.Characteristic, error) {
	chs, err := p.DiscoverCharacteristics(nil, service)
	if err != nil {
		fmt.Printf("Failed to discover characteristics, err: %s\n", err)
		return nil, err
	}
	for _, item := range chs {
		if item.Name() == name {
			return item, nil
		}
	}
	return nil, fmt.Errorf("Characteristic %s not found", name)
}

// HRHandler handles heart rate data
func HRHandler(ch *gatt.Characteristic, data []byte, err error) {
	heartRate := binary.LittleEndian.Uint16(append([]byte(data[1:2]), []byte{0}...))
	fmt.Printf("HR: %d\n", heartRate)
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

	d.Init(onStateChanged)
	<-done
	fmt.Println("Done")
}
