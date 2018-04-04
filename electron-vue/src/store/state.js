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
        warmUpDuration: 5 * 60, // s
    },
    race: {
        totalDistance: 0, // m, route distance provided manually or calculated from gpx
        maxSpeed: 0, // km/h

        // Current values
        speed: 0,
        distance: 0,
        grade: 0,

        videoFile: null,
        gpxData: [],
        currentGPXID: 0,
        opponents: [{
            name: "You",
            distance: 0,
        }],

        currentBPM: 0,
        energy: 0, // J

        startedAt: null,
        finishedAt: null,

        points: 0,

        csc: {
            speed: 0,
            distance: 0,
            points: 0, // Amount of recieved datapoints from csc sensor
        },
    },
};

export default state;
