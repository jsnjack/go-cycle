let defaultAvailableDevice = {
    name: null,
    id: null,
    connecting: null,
};

import utils from "../utils/gpx";
import real from "../utils/real";
import trainer from "../trainers/cycleops_fluid";
import {saveConfig} from "../utils/config";

const mutations = {
    WS_CONNECTED(state, connection) {
        state.ws.obj = connection;
        state.ws.connected = true;
    },
    WS_DISCONNECTED(state) {
        state.ws.connected = false;
    },
    DEVICE_DISCOVERED(state, device) {
        for (let i = 0; i < state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === device.id) {
                console.info("Device already registered", device);
                return;
            }
        }
        let newDevice = Object.assign({}, defaultAvailableDevice, device);
        state.devices.availableDevices.push(newDevice);
    },
    DEVICE_CONNECTING(state, id) {
        for (let i = 0; i < state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === id) {
                state.devices.availableDevices[i].connecting = true;
                return;
            }
        }
        console.warn(`Device ${id} not found in availableDevices`);
    },
    DEVICE_CONNECTED(state, data) {
        for (let i = 0; i < state.devices.availableDevices.length; i++) {
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
        case "csc_speed":
            state.devices.csc_speed.id = data.id;
            state.devices.csc_speed.connected = true;
            break;
        case "csc_cadence":
            state.devices.csc_cadence.id = data.id;
            state.devices.csc_cadence.connected = true;
            break;
        }
    },
    DEVICE_DISCONNECTED(state, data) {
        for (let i = 0; i < state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === data.id) {
                state.devices.availableDevices[i].connecting = false;
                break;
            }
        }
        switch (data.id) {
        case state.devices.hr.id:
            state.devices.hr.connected = false;
            break;
        case state.devices.csc_speed.id:
            state.devices.csc_speed.connected = false;
            break;
        case state.devices.csc_cadence.id:
            state.devices.csc_cadence.connected = false;
            break;
        }
    },
    MEASUREMENT_HR(state, data) {
        if (!state.devices.hr.connected) {
            this.commit("DEVICE_CONNECTED", data);
        }
        let energyPerMin;
        let now = new Date().getTime();
        let period = now - state.race.lastHREvent; // ms
        state.race.currentBPM = data.bpm;
        // https://community.fitbit.com/t5/Charge-HR/How-Charge-HR-calculates-calories-burned/td-p/1021859
        if (state.race.lastHREvent !== 0 && state.race.startedAt && !state.race.finishedAt) {
            /* eslint-disable max-len */
            if (state.user.gender == "m") {
                energyPerMin = -55.0969 + 0.6309 * state.race.currentBPM + 0.1988 * state.user.weight + 0.2017 * state.user.age;
            } else {
                energyPerMin = -20.4022 + 0.4472 * state.race.currentBPM - 0.1263 * state.user.weight + 0.074 * state.user.age;
            }
            /* eslint-enable max-len */
            state.race.calories += energyPerMin * 0.23 / 1000 / 60 * period; // kCal
        }
        state.race.lastHREvent = now;
    },
    MEASUREMENT_CSC_SPEED(state, data) {
        let started = performance.now();
        let grade = 0;
        let toLog = {};
        if (!state.devices.csc_speed.connected) {
            this.commit("DEVICE_CONNECTED", data);
        }
        state.race.csc_speed.distance =
            (data.revolutions * state.user.wheelSize) / 1000 || 0;
        state.race.csc_speed.speed = real.ensureSane(real.toKMH(
            state.race.csc_speed.distance / (data.time / 1000) || 0
        ), 100);
        try {
            state.race.currentPower = trainer.getYfromX(state.race.csc_speed.speed);
        } catch (e) {
            state.race.currentPower = 0;
        }
        if (state.race.startedAt && !state.race.finishedAt) {
            // Race is in progress
            state.race.csc_speed.points++;
            let estDistance = state.race.distance + state.race.csc_speed.distance;
            for (
                let i = state.race.currentGPXID;
                i < state.race.gpxData.length - 1;
                i++
            ) {
                // Deal with a bad estimation - happens with crazy high grades
                if (estDistance < state.race.gpxData[i].distance) {
                    i = i - 2;
                } else {
                    if (
                        state.race.gpxData[i].distance <= estDistance &&
                        state.race.gpxData[i + 1].distance > estDistance
                    ) {
                        state.race.currentGPXID = i;
                        grade = state.race.gpxData[i].grade;
                        toLog.grade = state.race.gpxData[i].grade;
                        toLog.ele = state.race.gpxData[i].elevation;
                        break;
                    }
                }
            }
            state.race.grade = grade;
            state.race.speed = real.ensureSane(real.getRealSpeed(
                state.race.csc_speed.speed,
                grade,
                state.user.weight
            ), 200);
            toLog.speed = state.race.speed;
            toLog.CSCSpeedSpeed = state.race.csc_speed.speed;
            if (state.race.speed > state.race.maxSpeed) {
                state.race.maxSpeed = state.race.speed;
            }
            state.race.distance +=
                (real.toMS(state.race.speed) * data.time) / 1000;
            state.race.opponents[0].distance = state.race.distance;
            toLog.distance = state.race.distance;
            let point = {
                time: new Date().toISOString(),
                distance: state.race.distance,
                hr: state.race.currentBPM,
                power: state.race.currentPower,
                cadence: state.race.currentCadence,
            };
            localStorage.setItem(
                "trkpt_" + state.race.points,
                JSON.stringify(point)
            );
            state.race.points++;
            if (state.race.distance > state.race.totalDistance) {
                this.commit("FINISH_RACE");
            }
        }
        let finished = performance.now();
        console.debug(
            `CSC Speed data: took ${finished - started}`, toLog
        );
    },
    MEASUREMENT_CSC_CADENCE(state, data) {
        let started = performance.now();
        let toLog = {};
        if (!state.devices.csc_cadence.connected) {
            this.commit("DEVICE_CONNECTED", data);
        }
        let rpm = data.revolutions / (data.time / 1000) * 60;
        rpm = real.ensureSane(rpm, 200);
        toLog.rpm = rpm;
        state.race.recentCadences.push(Math.round(rpm));
        if (state.race.recentCadences.length > 3) {
            state.race.recentCadences.shift();
        }
        let sum = 0;
        let sampleSize = state.race.recentCadences.length;
        for (let i=0; i<state.race.recentCadences.length; i++) {
            let value = state.race.recentCadences[i];
            if (value !== 0) {
                sum = sum + state.race.recentCadences[i];
            } else {
                sampleSize = sampleSize - 1;
            }
        }
        let currentCadence = 0;
        if (sampleSize > 0) {
            currentCadence = Math.round(sum/sampleSize);
            currentCadence = real.ensureSane(currentCadence, 200);
        }
        toLog.currentAvg = currentCadence;
        state.race.currentCadence = currentCadence;
        let finished = performance.now();
        console.debug(
            `CSC Cadence data: took ${finished - started}`, toLog
        );
    },
    VIDEOFILE_URL(state, urlObj) {
        state.race.videoFile = urlObj;
    },
    SET_GPX_DOC(state, doc) {
        let data = utils.extractDataFromGPX(doc);
        state.race.totalDistance = data[data.length - 1].distance;
        if (data[0].time === 0) {
            let hasOpponent = false;
            for (let i=0; state.race.opponents.length -1; i++) {
                if (state.race.opponents[i].name === "Joe") {
                    hasOpponent = true;
                    state.race.opponents[i].distance = 0;
                }
            }
            if (!hasOpponent) {
                state.race.opponents.push({name: "Joe", distance: 0});
            }
        }
        state.race.gpxData = data;
    },
    NEW_RACE(state) {
        state.race.startedAt = null;
        state.race.finishedAt = null;
    },
    START_RACE(state) {
        state.race.currentPower = 0;
        state.race.currentBPM = 0;
        state.race.recentCadences = [];
        state.race.currentCadence = 0;
        state.race.startedAt = new Date();
        state.race.finishedAt = null;
        state.race.csc_speed.time = 0;
        state.race.csc_speed.revolutions = 0;
        state.race.csc_speed.points = 0;
        state.race.calories = 0;
        state.race.lastHREvent = 0;
        state.race.currentGPXID = 0;
        state.race.distance = 0;
        state.race.speed = 0;
        state.race.grade = 0;
        state.race.points = 0;
        state.race.maxSpeed = 0;
        localStorage.clear();
    },
    FINISH_RACE(state) {
        state.race.finishedAt = new Date();
    },
    SET_OPPONENT_DISTANCE(state, obj) {
        state.race.opponents[obj.id].distance = obj.distance;
    },
    UPDATE_USER(state, obj) {
        Object.keys(obj).forEach((key) => {
            switch (key) {
            case "weight":
            case "age":
            case "warmUpDuration":
                state.user[key] = parseInt(obj[key], 10);
                break;
            case "gender":
            case "wheelSize":
            case "stravaAccessToken":
                state.user[key] = obj[key];
                break;
            default:
                console.warn("Unknown state.user key:", key);
            }
        });
        saveConfig(state.user);
    },
};

export default mutations;
