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
    DEVICE_CONNECTED(state, id, as) {
        for (let i = 0; i<state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === id) {
                state.devices.availableDevices[i].connecting = false;
                return;
            }
        }
    },
    DEVICE_DISCONNECTED(state, id, as) {
        for (let i = 0; i<state.devices.availableDevices.length; i++) {
            if (state.devices.availableDevices[i].id === id) {
                state.devices.availableDevices[i].connecting = false;
                return;
            }
        }
    },
};

export default mutations;
