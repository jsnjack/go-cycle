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
            <div class="controls-container">
                <button class="button-control" @click="onNext">Next</button>
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
        methods: {
            onNext() {
                this.$router.push("prerace");
            },
        },
        data() {
            return{
                status: "Scanning..."
            }
        }
    };
</script>

<style scoped>
    @import url("../assets/style.css");

    #devices {
        margin-top: 3rem;
    }

    .container {
        display: flex;
        flex-flow: row wrap;
        width: 100%;
    }
</style>
