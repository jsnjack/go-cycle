let defaultAvailableDevice = {
    name: null,
    id: null,
    connecting: null,
};

import utils from "../utils/gpx";
import real from "../utils/real";
import trainer from "../trainers/cycleops_fluid";


const mutations = {
    WS_CONNECTED(state, connection) {
        state.ws.obj = connection;
        state.ws.connected = true;
    },
    WS_DISCONNECTED(state) {
        state.ws.connected = false;
    },
    DEVICE_DISCOVERED(state, device) {
        for (let i = 0; i<state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === device.id) {
                console.info("Device already registered", device);
                return;
            }
        }
        let newDevice = Object.assign({}, defaultAvailableDevice, device);
        state.devices.availableDevices.push(newDevice);
    },
    DEVICE_CONNECTING(state, id) {
        for (let i = 0; i<state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === id) {
                state.devices.availableDevices[i].connecting = true;
                return;
            }
        }
        console.warn(`Device ${id} not found in availableDevices`);
    },
    DEVICE_CONNECTED(state, data) {
        for (let i = 0; i<state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === data.id) {
                state.devices.availableDevices[i].connecting = false;
                break;
            }
        }
        switch (data.recognizedAs) {
            case "hr":
                state.devices.hr.id = data.id;
                state.devices.hr.connected = true;
                break;
            case "csc":
                state.devices.csc.id = data.id;
                state.devices.csc.connected = true;
                break;
        }
    },
    DEVICE_DISCONNECTED(state, data) {
        for (let i = 0; i<state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === data.id) {
                state.devices.availableDevices[i].connecting = false;
                break;
            }
        }
        switch (data.recognizedAs) {
            case "hr":
                state.devices.hr.id = data.id;
                state.devices.hr.connected = false;
                break;
            case "csc":
                state.devices.csc.id = data.id;
                state.devices.csc.connected = false;
                break;
        }
    },
    MEASUREMENT_HR(state, data) {
        if (!state.devices.hr.connected) {
            this.commit("DEVICE_CONNECTED", data);
        }
        state.race.currentBPM = data.bpm;
    },
    MEASUREMENT_CSC(state, data) {
        let started = performance.now();
        let grade = 0;
        if (!state.devices.csc.connected) {
            this.commit("DEVICE_CONNECTED", data);
        }
        state.race.csc.distance = data.revolutions * state.user.wheelSize / 1000;
        state.race.csc.speed = state.race.csc.distance / (data.time / 1000);
        state.race.currentPower = trainer.getYfromX(state.race.csc.speed);
        if (state.race.startedAt && !state.race.finishedAt) {
            // Race is in progress
            state.race.csc.points++;
            let estDistance = state.race.distance + state.race.csc.distance;
            for (let i=state.race.currentGPXID; i<state.race.gpxData.length - 1; i++) {
                state.race.currentGPXID = i;
                grade = state.race.gpxData[i].grade;
                if (state.race.gpxData[i].distance <= estDistance && state.race.gpxData[i+1].distance > estDistance) {
                    break;
                }
            }
            state.race.grade = grade;
            state.race.speed = real.getRealSpeed(state.race.csc.speed, grade, state.user.weight);
            if (state.race.speed > state.race.maxSpeed) {
                state.race.maxSpeed = state.race.speed;
            }
            state.race.distance += real.toMS(state.race.speed) * data.time / 1000;
            state.race.opponents[0].distance = state.race.distance;
            let point = {
                time: new Date().toISOString(),
                distance: state.race.distance,
                hr: state.race.currentBPM,
            };
            localStorage.setItem("trkpt_" + state.race.points, JSON.stringify(point));
            state.race.points++;
            state.race.avgPower = (state.race.avgPower * (state.race.points - 1) + state.race.currentPower) / state.race.points;
            if (state.race.distance > state.race.totalDistance) {
                this.commit("FINISH_RACE");
            }
        }
        let finished = performance.now();
        console.debug(
            `CSC data: took ${finished - started}, grade, ${grade}, speed ${state.race.speed} (${state.race.csc.speed})`
        );
    },
    VIDEOFILE_URL(state, urlObj) {
        state.race.videoFile = urlObj;
    },
    SET_GPX_DOC(state, doc) {
        let data = utils.extractDataFromGPX(doc);
        state.race.totalDistance = data[data.length -1].distance;
        if (data[0].time === 0) {
            state.race.opponents.push({name: "Joe", distance: 0});
        }
        state.race.gpxData = data;
    },
    NEW_RACE(state) {
        state.race.startedAt = null;
        state.race.finishedAt = null;
    },
    START_RACE(state) {
        state.race.currentBPM = 0;
        state.race.startedAt = new Date();
        state.race.finishedAt = null;
        state.race.csc.time = 0;
        state.race.csc.revolutions = 0;
        state.race.csc.points = 0;
        localStorage.clear();
    },
    FINISH_RACE(state) {
        state.race.finishedAt = new Date();
    },
    UPDATE_USER_WEIGHT(state, value) {
        state.user.weight = parseInt(value, 10);
    },
    UPDATE_WARM_UP_DURATION(state, value) {
        state.user.warmUpDuration = parseInt(value, 10);
    },
    UPDATE_USER_GENDER(state, value) {
        state.user.gender = value;
    },
    UPDATE_USER_AGE(state, value) {
        state.user.age = parseInt(value, 10);
    },
    UPDATE_USER_WHEEL_SIZE(state, value) {
        state.user.wheelSize = value;
    },
    UPDATE_SIMPLE_ROUTE_DISTANCE(state, value) {
        state.race.totalDistance = parseInt(value, 10);
    },
    SET_OPPONENT_DISTANCE(state, obj) {
        state.race.opponents[obj.id].distance = obj.distance;
    },
};

export default mutations;
