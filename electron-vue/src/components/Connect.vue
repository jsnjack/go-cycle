<template>
    <div id="devices">
        <h1>Connect devices</h1>
        <div>
            <div v-show="!devices.availableDevices.length">
                {{ status }}
            </div>
        </div>
    </div>
</template>
<script>
    import vuex from "vuex";

    export default {
        name: 'Connect',
        mounted () {
            setTimeout(() => {
                this.$store.dispatch("ws_startScanning");
            }, 100);
        },
        computed: {
            ...vuex.mapState([
                "devices",
            ]),
        },
        data() {
            return{
                status: "Scanning..."
            }
        }
    };
</script>

<style scoped>
    h1 {
        text-align: center;
    }
</style>
