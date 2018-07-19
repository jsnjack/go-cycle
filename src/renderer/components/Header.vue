<template>
    <div id="header">
        <div>
            <span @click="finishRace"><FinishIcon v-show="isRaceInProgress" class="icon finish"/></span>
            <HRConnectedIcon class="icon" :class="{offline: isHRConnected}"/>
            <CSCConnectedIcon class="icon" :class="{offline: isCSCConnected}"/>
            <ConnectedIcon class="icon" :class="{offline: isOffline}"/>
        </div>
    </div>
</template>
<script>
import vuex from "vuex";
import ConnectedIcon from "../assets/connected.svg";
import HRConnectedIcon from "../assets/hr-connected.svg";
import CSCConnectedIcon from "../assets/speed-meter.svg";
import FinishIcon from "../assets/finish.svg";


export default {
    name: "Header",
    components: {ConnectedIcon, HRConnectedIcon, CSCConnectedIcon, FinishIcon},
    computed: {
        ...vuex.mapState([
            "ws",
            "devices",
            "race",
        ]),
        isOffline: function() {
            return !this.ws.connected;
        },
        isHRConnected: function() {
            return !this.devices.hr.connected;
        },
        isCSCConnected: function() {
            return !this.devices.csc.connected;
        },
        isRaceInProgress: function() {
            return this.race.startedAt && !this.race.finishedAt;
        },
    },
    methods: {
        finishRace() {
            this.$store.dispatch("finish_race");
        },
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
    background: rgb(36, 36, 36);
    opacity: 0.2;
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
.finish {
    margin-right: 1rem;
    fill: #fde74c;
}
.finish:hover {
    cursor: pointer;
}
</style>
