const wsMessageHandler = function(app, data) {
    let msg = JSON.parse(data);
    switch (msg.type) {
        case "ws.device:discovered":
            app.$store.commit("DEVICE_DISCOVERED", msg.data);
            break;
        case "ws.device:status":
        console.log(msg.data);
            if (msg.data.status === "connected") {
                app.$store.commit("DEVICE_CONNECTED", msg.data);
            } else if (msg.data.status === "disconnected") {
                app.$store.commit("DEVICE_DISCONNECTED", msg.data);
            } else {
                console.warn("Unexpected device status", msg);
            }
            break;
        default:
            console.warn("Unhandled message", msg);
    }
};

export default wsMessageHandler;
