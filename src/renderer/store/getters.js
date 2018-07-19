
const getters = {
    routeProgress(state) {
        return Math.round(state.race.distance/state.race.totalDistance * 100 * 100) / 100 || 0;
    },
    isRaceFinished(state) {
        return state.race.startedAt && state.race.finishedAt;
    },
    isRaceInProgress(state) {
        return state.race.startedAt && !state.race.finishedAt;
    },
};

export default getters;
