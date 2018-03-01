const state = {
    ws: {
        obj: null,
        url: "ws://localhost:8000/ws",
        connected: false,
        reconnectTimeout: 1000, // ms
    },
    devices: {
        availableDevices: [],
        hr: {
            connected: false,
        },
        csc: {
            connected: false,
        },
    },
};

export default state;
