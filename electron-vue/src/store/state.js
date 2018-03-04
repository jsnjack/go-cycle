const state = {
    ws: {
        obj: {sendMessage: function() {}},
        url: "ws://localhost:8000/ws",
        connected: false,
        reconnectTimeout: 1000, // ms
    },
    devices: {
        availableDevices: [], // [{name: "MOOV", id: "45:fg:56", connecting: false}]
        hr: {
            id: null,
            connected: false,
        },
        csc: {
            id: null,
            connected: false,
        },
    },
    user: {
        weight: 65, // kg
        gender: "m", // male
        age: 30,
        wheelSize: 622 + 28*2, // mm
    },
    race: {
        videoFile: null,

        currentBPM: 0,
        calories: 0, // kCals
        lastHREvent: 0,

        currentRevPerSec: 0,
        currentRevolutions: 0,
        totalRevolutions: 0,
    },
};

export default state;
