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
        }
    },
    DEVICE_DISCONNECTED(state, data) {
        for (let i = 0; i < state.devices.availableDevices.length; i++) {
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
        case "csc_speed":
            state.devices.csc_speed.id = data.id;
            state.devices.csc_speed.connected = false;
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
        state.race.csc_speed.speed = real.toKMH(
            state.race.csc_speed.distance / (data.time / 1000) || 0
        );
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
                state.race.currentGPXID = i;
                grade = state.race.gpxData[i].grade;
                if (grade === Infinity) {
                    grade = 0;
                }
                toLog.grade = state.race.gpxData[i].grade;
                toLog.ele = state.race.gpxData[i].elevation;
                if (
                    state.race.gpxData[i].distance <= estDistance &&
                    state.race.gpxData[i + 1].distance > estDistance
                ) {
                    break;
                }
            }
            state.race.grade = grade;
            state.race.speed = real.getRealSpeed(
                state.race.csc_speed.speed,
                grade,
                state.user.weight
            );
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
    UPDATE_USER_WEIGHT(state, value) {
        state.user.weight = parseInt(value, 10);
    },
    UPDATE_USER_GENDER(state, value) {
        state.user.gender = value;
    },
    UPDATE_USER_AGE(state, value) {
        state.user.age = parseInt(value, 10);
    },
    UPDATE_WARM_UP_DURATION(state, value) {
        state.user.warmUpDuration = parseInt(value, 10);
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
