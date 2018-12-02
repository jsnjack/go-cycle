package main

import (
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/paypal/gatt"
)

// CSCMessage is a message from the CSC sensor
type CSCMessage struct {
	ID           string     `json:"id"` // Device id
	RecognizedAs SensorKind `json:"recognizedAs"`
	Revolutions  uint32     `json:"revolutions"` // Amount of wheel revolutions since last time, for calculating distance
	Time         uint16     `json:"time"`        // Time since last measurement, ms
}

// SpeedSensorData is data from the sensor
type SpeedSensorData struct {
	Revolutions uint32
	EventTime   uint16
}

// CSCSensor ...
type CSCSensor struct {
	Sensor
	Previous SpeedSensorData
	Current  SpeedSensorData
}

// Listen ...
func (sensor *CSCSensor) Listen() {
	Logger.Println("Setting up CSC sensor")
	defer func() {
		sensor.Peripheral.Device().CancelConnection(sensor.Peripheral)
	}()
	level, err := sensor.GetBatteryLevel()
	if err != nil {
		Logger.Println(err)
	} else {
		Logger.Printf("Battery: %d\n", level)
	}
	service, err := GetService(sensor.Peripheral, gatt.UUID16(0x1816))
	if err != nil {
		Logger.Println(err)
		return
	}

	ch, err := GetCharacteristic(sensor.Peripheral, service, gatt.UUID16(0x2A5B))
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

func (sensor *CSCSensor) decode(data []byte) {
	offset := 1
	var revolutions uint32
	var eventTime uint16
	switch data[0] {
	case 1:
		// Speed sensor
		sensor.Kind = SpeedKind
		revolutions = binary.LittleEndian.Uint32(append([]byte(data[offset:])))
		offset += 4
		eventTime = binary.LittleEndian.Uint16(append([]byte(data[offset:])))
		break
	case 2:
		// Cadence sensor
		sensor.Kind = CadenceKind
		_revolutions := binary.LittleEndian.Uint16(append([]byte(data[offset:])))
		revolutions = uint32(_revolutions)
		offset += 2
		eventTime = binary.LittleEndian.Uint16(append([]byte(data[offset:])))
		break
	}
	cscData := SpeedSensorData{Revolutions: revolutions, EventTime: eventTime}
	if sensor.hasPrevious() {
		sensor.Previous = sensor.Current
		sensor.Current = cscData
	} else {
		sensor.Previous = cscData
		sensor.Current = cscData
		return
	}
	var time uint16
	if sensor.Current.EventTime >= sensor.Previous.EventTime {
		time = sensor.Current.EventTime - sensor.Previous.EventTime
	} else {
		time = 65535 - sensor.Previous.EventTime + sensor.Current.EventTime + 1
	}
	Logger.Printf("[%s] Rev: %d, Time: %d\n", sensor.Kind, sensor.Current.Revolutions, time)
	msgCSC := CSCMessage{
		ID:           sensor.Peripheral.ID(),
		RecognizedAs: sensor.Kind,
		Revolutions:  sensor.Current.Revolutions - sensor.Previous.Revolutions,
		Time:         time,
	}
	msgWS := WSMessage{Type: "ws.device:measurement", Data: msgCSC}
	msgB, err := json.Marshal(msgWS)
	if err != nil {
		Logger.Println(err)
	} else {
		BroadcastChannel <- msgB
	}
}

func (sensor *CSCSensor) hasPrevious() bool {
	if sensor.Previous.EventTime != 0 || sensor.Previous.Revolutions != 0 {
		return true
	}
	return false
}

// GetType returns type of the sensor
func (sensor *CSCSensor) GetType() PeripheralType {
	return CSCPeripheral
}

// SendSynthCSCEvent sends synthetic CSC event
func SendSynthCSCEvent() {
	msgSpeed := CSCMessage{
		ID:           "fake-speed-csc",
		RecognizedAs: SpeedKind,
		Revolutions:  uint32(Random(4, 6)),
		Time:         1000,
	}
	msgWS := WSMessage{Type: "ws.device:measurement", Data: msgSpeed}
	msgB, _ := json.Marshal(msgWS)
	BroadcastChannel <- msgB

	msgCadence := CSCMessage{
		ID:           "fake-cadence-csc",
		RecognizedAs: CadenceKind,
		Revolutions:  uint32(Random(1, 3)),
		Time:         1000,
	}
	msgWS = WSMessage{Type: "ws.device:measurement", Data: msgCadence}
	msgB, _ = json.Marshal(msgWS)
	BroadcastChannel <- msgB
}

// Random generates random integer number within threshold
func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
