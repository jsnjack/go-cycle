package main

import (
	"encoding/binary"
	"encoding/json"

	"github.com/paypal/gatt"
)

// HRServiceUUID is UUID for heart_rate service
var HRServiceUUID = gatt.UUID16(0x180d)

// HRMessage is a message from the HR sensor
type HRMessage struct {
	ID           string `json:"id"` // Device id
	RecognizedAs string `json:"recognizedAs"`
	BPM          uint16 `json:"bpm"`
}

// HRSensor ...
type HRSensor struct {
	Peripheral  gatt.Peripheral
	Initialized bool
}

// Listen ...
func (sensor *HRSensor) Listen() {
	Logger.Println("Setting up HR sensor")
	defer func() {
		sensor.Peripheral.Device().CancelConnection(sensor.Peripheral)
	}()

	if err := sensor.Peripheral.SetMTU(500); err != nil {
		Logger.Printf("Failed to set MTU, err: %s\n", err)
	}

	service, err := GetService(sensor.Peripheral, gatt.UUID16(0x180d))
	if err != nil {
		Logger.Println(err)
		return
	}

	ch, err := GetCharacteristic(sensor.Peripheral, service, gatt.UUID16(0x2a37))
	if err != nil {
		Logger.Println(err)
		return
	}

	_, err = sensor.Peripheral.DiscoverDescriptors(nil, ch)
	if err != nil {
		Logger.Println(err)
		return
	}

	resultCh := make(chan string)
	sensor.Peripheral.SetNotifyValue(ch, func(ch *gatt.Characteristic, data []byte, err error) {
		if err != nil {
			resultCh <- err.Error()
		}
		sensor.decode(data)
	})
	<-resultCh
}

func (sensor *HRSensor) decode(data []byte) {
	heartRate := binary.LittleEndian.Uint16(append([]byte(data[1:2]), []byte{0}...))
	Logger.Printf("BPM: %d\n", heartRate)
	msgHR := HRMessage{
		ID:           sensor.Peripheral.ID(),
		RecognizedAs: GetActiveDeviceType(sensor.Peripheral.ID()),
		BPM:          heartRate,
	}
	msgWS := WSMessage{Type: "ws.device:measurement", Data: msgHR}
	msgB, err := json.Marshal(msgWS)
	if err != nil {
		Logger.Println(err)
	} else {
		BroadcastChannel <- msgB
	}
}

// GetType returns type of the sensor
func (sensor *HRSensor) GetType() PeripheralType {
	return HRPeripheral
}

// GetPeripheral ...
func (sensor *HRSensor) GetPeripheral() gatt.Peripheral {
	return sensor.Peripheral
}
