<template>
    <div id="header">
        <div>
            <span>Status: {{ getWSStatus }}</span>
        </div>
    </div>
</template>
<script>
import vuex from "vuex";

export default {
    name: 'Header',
    computed: {
        ...vuex.mapState([
            "ws",
        ]),
        getWSStatus: function () {
            return this.ws.connected ? "Connected": "Offline";
        },
    },
};
</script>

<style scoped>
</style>
