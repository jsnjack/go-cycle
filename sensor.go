package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/godbus/dbus"
	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile"
	"github.com/muka/go-bluetooth/emitter"
)

// SensorKind kind of the sensor, depends on returned measurements
type SensorKind string

// Sensor is a common struct for all sensors
type Sensor struct {
	Address  string
	Kind     SensorKind
	Char     *profile.GattCharacteristic1
	Device   *api.Device
	Logger   *log.Logger
	previous SpeedSensorData
	current  SpeedSensorData
}

// Listen subscribes to the events
func (sensor *Sensor) Listen() {
	sensor.Logger.Println("Subscribing for notifications")
	sensor.ListenChanges()

	dataCh, err := sensor.Char.Register()
	if err != nil {
		sensor.Logger.Println("Failed to register.", err)
		return
	}
	go func() {
		for event := range dataCh {
			if sensor.Char.Path == fmt.Sprintf("%s", event.Path) {
				switch event.Body[0].(type) {
				case dbus.ObjectPath:
					continue
				case string:
				}

				if event.Body[0] != bluez.GattCharacteristic1Interface {
					continue
				}
				props := event.Body[1].(map[string]dbus.Variant)
				if _, ok := props["Value"]; !ok {
					continue
				}
				value := props["Value"].Value().([]byte)
				switch sensor.Kind {
				case HRKind:
					sensor.handleHR(value)
				default:
					sensor.handleCSC(value)
				}
			}
		}
	}()
	err = sensor.Char.StartNotify()
	if err != nil {
		sensor.Logger.Println("Failed to start notifications.", err)
		return
	}
}

// ListenChanges listens for changes from the device
func (sensor *Sensor) ListenChanges() {
	sensor.Logger.Println("Listening for changes...")
	sensor.Device.On("changed", emitter.NewCallback(func(ev emitter.Event) {
		evData := ev.GetData().(api.PropertyChangedEvent)
		sensor.Logger.Println("Change:", evData.Field, evData.Value)
		switch evData.Field {
		case "Connected":
			if !evData.Value.(bool) {
				sensor.Logger.Println("Disconnected.")

				msgStatus := DeviceStatusData{ID: sensor.Address, Status: "disconnected"}
				wsMsgStatus := WSMessage{Type: "ws.device:status", Data: msgStatus}
				msgB, err := json.Marshal(&wsMsgStatus)
				if err != nil {
					Logger.Println(err)
				} else {
					BroadcastChannel <- msgB
				}

				for i, address := range ConnectedDevices {
					if address == sensor.Address {
						ConnectedDevices = append(ConnectedDevices[:i], ConnectedDevices[i+1:]...)
					}
					break
				}
				manager, err := api.GetManager()
				if err != nil {
					sensor.Logger.Println(err)
				}
				sensor.Logger.Println("Refreshing state...")
				err = manager.RefreshState()
				if err != nil {
					sensor.Logger.Println(err)
				}
				ConnectToDevice(sensor.Address)
			}
			break
		}
	}))
}

func (sensor *Sensor) hasPrevious() bool {
	if sensor.previous.EventTime != 0 || sensor.previous.Revolutions != 0 {
		return true
	}
	return false
}

// Random generates random integer number within threshold
func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
