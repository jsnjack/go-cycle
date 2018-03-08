let defaultAvailableDevice = {
    name: null,
    id: null,
    connecting: null,
};

import utils from "../utils/gpx";


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
        let energyPerMin;
        let now = new Date().getTime();
        let period = now - state.race.lastHREvent; // ms
        state.race.currentBPM = data.bpm;
        // https://community.fitbit.com/t5/Charge-HR/How-Charge-HR-calculates-calories-burned/td-p/1021859
        if (state.race.lastHREvent !== 0) {
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
    MEASUREMENT_CSC(state, data) {
        if (!state.devices.csc.connected) {
            this.commit("DEVICE_CONNECTED", data);
        }
        state.race.currentRevolutions = data.revolutions;
        state.race.totalRevolutions += data.revolutions;
        state.race.currentRevPerSec = data.rev_per_sec;
    },
    VIDEOFILE_URL(state, urlObj) {
        state.race.videoFile = urlObj;
    },
    SET_GPX_DOC(state, doc) {
        let started = performance.now();
        state.race.gpxDoc = doc;
        let total = 0;
        let points = doc.getElementsByTagName("trkpt");
        let container = new Array(points.length);
        container[0] = {distance: 0, elevation: points[0].getElementsByTagName("ele")[0].textContent};
        for (let i=0; i<points.length - 1; i++) {
            total += utils.getDistance(
                points[i].getAttribute("lat"),
                points[i].getAttribute("lon"),
                points[i+1].getAttribute("lat"),
                points[i+1].getAttribute("lon")
            );
            container[i+1] = {distance: total, elevation: points[i+1].getElementsByTagName("ele")[0].textContent};
        }
        state.race.gpxDistToElev = container;
        let finished = performance.now();
        console.debug("gpxDistToElev construction took, ms:", finished - started);
    },
    START_RACE(state) {
        state.race.currentBPM = 0;
        state.race.calories = 0;
        state.race.lastHREvent = 0;
        state.race.startedAt = new Date();
        state.race.currentRevPerSec = 0;
        state.race.currentRevolutions = 0;
        state.race.totalRevolutions = 0;
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
    UPDATE_USER_WHEEL_SIZE(state, value) {
        state.user.wheelSize = value;
    },
    UPDATE_SIMPLE_ROUTE_DISTANCE(state, value) {
        state.race.simpleRouteDistance = parseInt(value, 10);
    },
};

export default mutations;
