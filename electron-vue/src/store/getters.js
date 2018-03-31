
const getters = {
    currentSpeed(state) {
        return state.race.csc.speed || 0;
    },
    distance(state, getters) {
        // Current distance in meters
        let value = 0 * state.race.totalRevolutions || 0;
        value = Math.round(value);
        return value;
    },
    distanceLeft(state, getters) {
        // Distance left in meters
        return Math.round(state.race.totalDistance - getters.distance) || 0;
    },
    routeProgress(state, getters) {
        return Math.round(getters.distance/state.race.totalDistance * 100 * 100) / 100 || 0;
    },
    isRaceFinished(state) {
        return state.race.startedAt && state.race.finishedAt;
    },
    isRaceInProgress(state) {
        return state.race.startedAt && !state.race.finishedAt;
    },
};

export default getters;
