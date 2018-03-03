<template>
    <div id="header">
        <div>
            <CSCConnectedIcon class="icon" :class="{offline: isCSCConnected}"/>
            <HRConnectedIcon class="icon" :class="{offline: isHRConnected}"/>
            <ConnectedIcon class="icon" :class="{offline: isOffline}"/>
        </div>
    </div>
</template>
<script>
import vuex from "vuex";
import ConnectedIcon from '../assets/connected.svg';
import HRConnectedIcon from '../assets/hr-connected.svg';
import CSCConnectedIcon from '../assets/speed-meter.svg';


export default {
    name: 'Header',
    components: {ConnectedIcon, HRConnectedIcon, CSCConnectedIcon},
    computed: {
        ...vuex.mapState([
            "ws",
            "devices"
        ]),
        isOffline: function () {
            return !this.ws.connected;
        },
        isHRConnected: function () {
            return !this.devices.hr.connected;
        },
        isCSCConnected: function () {
            return !this.devices.csc.connected;
        }
    },
};
</script>

<style scoped>
#header {
    z-index: 100;
    width: 100%;
    position: absolute;
    top: 0;
    left: 0;
    text-align: right;
    background: linear-gradient(rgba(0, 0, 0, 0.2));
}
.icon {
    height: 1.2rem;
    width: 1.2rem;
    margin: 0.15rem;
    fill: white;
}
.icon.offline {
    fill: rgb(77, 77, 77);
}
</style>
