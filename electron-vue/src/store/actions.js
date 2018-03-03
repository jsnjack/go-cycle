const actions = {
    ws_startScanning({state}) {
        let data = {
            type: "app.bt:scan",
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
};

export default actions;
