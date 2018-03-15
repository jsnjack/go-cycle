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
        wheelSize: "2136", // mm
        stravaAccessToken: window.process.env.GO_CYCLE_STRAVA_TOKEN,
    },
    race: {
        simpleRouteDistance: 0, // m, route distance provided manually or calculated from gpx

        videoFile: null,
        gpxData: [],
        opponents: [{
            name: "You",
            distance: 1,
        }],

        currentBPM: 0,
        calories: 0, // kCals
        lastHREvent: 0,

        startedAt: null,
        finishedAt: null,

        currentRevPerSec: 0,
        maxRevPerSec: 0,
        currentRevolutions: 0,
        totalRevolutions: 0,

        point: 0, // Amount of recieved datapoints from csc sensor
    },
};

export default state;
