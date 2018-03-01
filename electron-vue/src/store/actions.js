const actions = {
    ws_startScanning({state}) {
        let data = {
            type: "command",
            data: "scan",
        };
        state.ws.obj.send(data);
    },
};

export default actions;
