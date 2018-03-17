const getters = {
    currentSpeed(state, getters) {
        let value = Math.round(state.race.currentRevPerSec * getters.distancePerRev * 3.6 * 1000 * 1000 * 10) / 10 || 0;
        return value;
    },
    maxSpeed(state, getters) {
        let value = Math.round(state.race.maxRevPerSec * getters.distancePerRev * 3.6 * 1000 * 1000 * 10) / 10 || 0;
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
        return state.race.simpleRouteDistance;
    },
    distanceLeft(state, getters) {
        // Distance left in meters
        return Math.round(getters.routeDistance - getters.distance) || 0;
    },
    routeProgress(state, getters) {
        return Math.round(getters.distance/getters.routeDistance * 100 * 100) / 100 || 0;
    },
    isRaceFinished(state) {
        return state.race.startedAt && state.race.finishedAt;
    },
};

export default getters;
