
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
    calories(state, getters) {
        if (!state.race.startedAt) {
            return 0;
        }
        let timestamp = new Date();
        if (getters.isRaceFinished) {
            timestamp = state.race.finishedAt;
        }
        let e = state.race.avgPower * (timestamp.getTime() - state.race.startedAt.getTime()) / 1000;
        let calories = e / 4.184 / 1000;
        return calories;
    },
};

export default getters;
