package main

import (
	"encoding/binary"
	"encoding/json"
)

// SpeedKind measures speed
var SpeedKind SensorKind = "csc_speed"

// CadenceKind measures speed
var CadenceKind SensorKind = "csc_cadence"

// CSCMeasuremetCharID is Speed and Cadence Measurement
const CSCMeasuremetCharID = "2a5b"

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

func (sensor *Sensor) handleCSC(data []byte) {
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
		sensor.previous = sensor.current
		sensor.current = cscData
	} else {
		sensor.previous = cscData
		sensor.current = cscData
		return
	}
	var time uint16
	if sensor.current.EventTime >= sensor.previous.EventTime {
		time = sensor.current.EventTime - sensor.previous.EventTime
	} else {
		time = 65535 - sensor.previous.EventTime + sensor.current.EventTime + 1
	}
	sensor.Logger.Printf("[%s] Rev: %d, Time: %d\n", sensor.Kind, sensor.current.Revolutions, time)
	msgCSC := CSCMessage{
		ID:           sensor.Address,
		RecognizedAs: sensor.Kind,
		Revolutions:  sensor.current.Revolutions - sensor.previous.Revolutions,
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
