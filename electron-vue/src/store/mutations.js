let defaultAvailableDevice = {
    name: null,
    id: null,
    connecting: null,
};

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
                state.devices.hr.id = data.id;
                state.devices.hr.connected = true;
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
                state.devices.hr.id = data.id;
                state.devices.hr.connected = false;
                break;
        }
    },
    MEASUREMENT_HR(state, data) {
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
        state.race.currentRevolutions = data.revolutions;
        state.race.totalRevolutions += data.revolutions;
        state.race.currentRevPerSec = data.rev_per_sec;
    },
    VIDEOFILE_URL(state, urlObj) {
        state.race.videoFile = urlObj;
    },
};

export default mutations;
