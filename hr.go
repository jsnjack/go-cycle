package main

import (
	"encoding/binary"
	"encoding/json"
)

// HRKind measures heart rate
const HRKind SensorKind = "hr"

// HRMeasuremetCharID is Heart Rate Measurement
const HRMeasuremetCharID = "2a37"

// HRMessage is a message from the HR sensor
type HRMessage struct {
	ID           string     `json:"id"` // Device id
	RecognizedAs SensorKind `json:"recognizedAs"`
	BPM          uint16     `json:"bpm"`
}

func (sensor *Sensor) handleHR(data []byte) {
	heartRate := binary.LittleEndian.Uint16(append([]byte(data[1:2]), []byte{0}...))
	sensor.Logger.Printf("BPM: %d\n", heartRate)
	msgHR := HRMessage{
		ID:           sensor.Address,
		RecognizedAs: sensor.Kind,
		BPM:          heartRate,
	}
	msgWS := WSMessage{Type: "ws.device:measurement", Data: msgHR}
	msgB, err := json.Marshal(msgWS)
	if err != nil {
		sensor.Logger.Println(err)
	} else {
		BroadcastChannel <- msgB
	}
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
