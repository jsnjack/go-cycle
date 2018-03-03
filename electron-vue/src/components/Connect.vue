<template>
    <div id="devices">
        <h1>Connect devices</h1>
        <div>
            <div v-show="!devices.availableDevices.length">
                {{ status }}
            </div>
            <div class="container">
                <Device v-for="item in devices.availableDevices"
                    :key="item.id"
                    :device="item">
                </Device>
            </div>
        </div>
    </div>
</template>
<script>
    import vuex from "vuex";
    import Device from "./Device";

    export default {
        name: 'Connect',
        components : {
            Device
        },
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
    .container {
        display: inline-block;
        width: 100%;
    }
</style>
