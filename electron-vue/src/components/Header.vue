<template>
    <div id="header">
        <div>
            <ConnectedIcon class="icon" :class="{offline: isOffline}"/>
        </div>
    </div>
</template>
<script>
import vuex from "vuex";
import ConnectedIcon from '../assets/connected.svg';


export default {
    name: 'Header',
    components: {ConnectedIcon},
    computed: {
        ...vuex.mapState([
            "ws",
        ]),
        isOffline: function () {
            return !this.ws.connected;
        },
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
