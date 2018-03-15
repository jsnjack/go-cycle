const actions = {
    ws_startScanning({state}) {
        let data = {
            type: "app.bt:scan",
            data: {},
        };
        state.ws.obj.sendMessage(data);
    },
    ws_stopScanning({state}) {
        let data = {
            type: "app.bt:scan_stop",
            data: {},
        };
        state.ws.obj.sendMessage(data);
    },
    ws_connectDevice({state, commit}, id) {
        commit("DEVICE_CONNECTING", id);
        let data = {
            type: "app.device:connect",
            data: {
                id: id,
            },
        };
        state.ws.obj.sendMessage(data);
    },
    finish_race({state, commit}) {
        commit("FINISH_RACE");
    },
    update_opponent({state, commit, dispatch}, pointID) {
        if (pointID < state.race.gpxData.length && state.race.gpxData[0].time >= 0) {
            let timeout = state.race.gpxData[pointID + 1].time - state.race.gpxData[pointID].time;
            setTimeout(()=>{
                commit("SET_OPPONENT_DISTANCE", {id: 1, distance: state.race.gpxData[pointID + 1].distance});
                dispatch("update_opponent", pointID + 1);
            }, timeout);
        }
    },
};

export default actions;
