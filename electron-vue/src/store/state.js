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
    race: {
        videoFile: null,
        wheelSize: 622 + 28*2,
        currentBPM: 0,
        currentRevPerSec: 0,
        currentRevolutions: 0,
        totalRevolutions: 0,
        bodyWeight: 65, // kg
    },
};

export default state;
