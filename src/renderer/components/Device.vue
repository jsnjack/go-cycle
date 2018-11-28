<template>
    <div class="device" @click="connect" :class="{connecting: deviceStatus, connected: isConnected}">
        <div class="thumbnail">
            <DeviceIcon class="icon"/>
        </div>
        <div class="name">
            {{ device.name }}
        </div>
    </div>
</template>
<script>
import vuex from "vuex";
import DeviceIcon from "../assets/device.svg";

export default {
    name: "Device",
    components: {
        DeviceIcon,
    },
    props: {
        device: {
            type: Object,
            required: true,
        },
    },
    methods: {
        connect() {
            if (!this.device.connecting && !this.isConnected) {
                this.$store.dispatch("ws_connectDevice", this.device.id);
            }
        },
    },
    computed: {
        ...vuex.mapState([
            "devices",
        ]),
        deviceStatus: function() {
            return this.device.connecting ? "connecting" : "";
        },
        isConnected: function() {
            if ((this.device.id === this.devices.hr.id && this.devices.hr.connected) ||
                (this.device.id === this.devices.csc_speed.id && this.devices.csc_speed.connected) ||
                (this.device.id === this.devices.csc_cadence.id && this.devices.csc_cadence.connected)) {
                return true;
            }
            return false;
        },
    },
};
</script>

<style scoped>
.device {
    position: relative;
    cursor: pointer;
    width: 10rem;
    height: 10rem;
    margin: 3em;
    padding: 1em;
    background: linear-gradient(rgba(0, 0, 0, 0.1), rgba(0, 0, 0, 0.1));
}
.device.connected::after{
    position: absolute;
    bottom: -0.4rem;
    right: 0;
    content: "🗸";
    font-size: 2rem;
}
.device.connecting {
    filter: blur(2px);
}
.thumbnail {
    text-align: center;
}
.icon {
    height: 5rem;
    width: 5rem;
    fill: white;
    margin: 1em;
}
.name {
    vertical-align: middle;
    font-size: 2rem;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
