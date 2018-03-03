<template>
    <div id="app">
        <Header/>
        <router-view/>
    </div>
</template>

<script>
import vuex from "vuex";
import Header from "./components/Header";
import Vue from "vue";
import wsMessageHandler from "./store/communication";

export default {
    name: "App",
    components: {Header},
    mounted() {
        let ws = new WebSocket(this.$store.state.ws.url);
        ws.sendMessage = function(obj) {
            ws.send(JSON.stringify(obj));
        };
        // Listen for messages
        ws.addEventListener('message', event => {
            wsMessageHandler(this, event.data);
        });

        ws.addEventListener('close', event => {
            this.$store.commit("WS_DISCONNECTED");
            setTimeout(this.connect, this.$store.state.ws.reconnectTimeout);
        });

        ws.addEventListener('error', event => {
            console.warn('WS connection error', event);
        });

        ws.addEventListener('open', event => {
            this.$store.commit("WS_CONNECTED", ws);
            if (this.$store.state.devices.availableDevices.length === 0) {
                // we need to connect new devices
                this.$router.push("connect");
            }
        });
    },
    computed: {
        ...vuex.mapState([
        ]),
    },
};
</script>
<style>
html,
body {
    margin: 0;
    padding: 0;
    background:#489671;
    color: white;
    font-size: 16px;
}
</style>
