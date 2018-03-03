const state = {
    ws: {
        obj: null,
        url: "ws://localhost:8000/ws",
        connected: false,
        reconnectTimeout: 1000, // ms
    },
    devices: {
        availableDevices: [], // [{name: "MOOV", id: "45:fg:56", connecting: false}]
        hr: {
            connected: false,
        },
        csc: {
            connected: false,
        },
    },
};

export default state;
