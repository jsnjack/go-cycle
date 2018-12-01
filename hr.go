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
	ID           string     `json:"id"` // Device id
	RecognizedAs SensorKind `json:"recognizedAs"`
	BPM          uint16     `json:"bpm"`
}

// HRSensor ...
type HRSensor struct {
	Peripheral gatt.Peripheral
	Kind       SensorKind
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
		ID:           sensor.GetID(),
		RecognizedAs: sensor.GetKind(),
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

// GetID ...
func (sensor *HRSensor) GetID() string {
	return sensor.Peripheral.ID()
}

// GetKind ...
func (sensor *HRSensor) GetKind() SensorKind {
	return sensor.Kind
}

// SendSynthHREvent sends synthetic HR event
func SendSynthHREvent() {
	msgHR := HRMessage{
		ID:           "fake-hr",
		RecognizedAs: HRKind,
		BPM:          uint16(Random(60, 100)),
	}
	msgWS := WSMessage{Type: "ws.device:measurement", Data: msgHR}
	msgB, _ := json.Marshal(msgWS)
	BroadcastChannel <- msgB
}
