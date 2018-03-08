const getters = {
    currentSpeed(state, getters) {
        let value = Math.round(state.race.currentRevPerSec * getters.distancePerRev * 3.6 * 1000 * 1000 * 10) / 10 || 0;
        value = value.toFixed(1);
        return value;
    },
    distancePerRev(state) {
        return parseInt(state.user.wheelSize, 10) / 1000 || 0; // meters
    },
    distance(state, getters) {
        // Total distance in meters
        let value = getters.distancePerRev * state.race.totalRevolutions || 0;
        value = Math.round(value);
        return value;
    },
    routeDistance(state) {
        // Route distance in meters
        let total = 0;
        if (state.race.gpxDistToElev.length) {
            return state.race.gpxDistToElev[state.race.gpxDistToElev.length - 1].distance;
        }
        return total;
    },
    distanceLeft(state) {
        // Distance left in meters
        return state.routeDistance - state.distance || 0;
    },
    routeProgress(state) {
        return Math.round(state.distance/state.routeDistance * 10) / 10 || 0;
    },
};

export default getters;
