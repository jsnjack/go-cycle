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
    text-align: right;
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
