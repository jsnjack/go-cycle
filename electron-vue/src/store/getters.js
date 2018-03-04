const getters = {
    currentSpeed(state, getters) {
        let value = Math.round(state.race.currentRevPerSec * getters.distancePerRev * 3.6 * 1000 * 1000 * 10) / 10 || 0;
        value = value.toFixed(1);
        return value;
    },
    distancePerRev(state) {
        return Math.Pi * state.user.wheelSize / 1000 || 0; // meters
    },
    distance(state, getters) {
        // Total distance in meters
        let value = getters.distancePerRev * state.race.totalRevolutions || 0;
        value = Math.round(value);
        return value;
    },
};

export default getters;
