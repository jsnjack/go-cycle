<template>
    <div class="device" @click="connect" :class="{connecting: deviceStatus}">
        <div class="thumbnail">
            <DeviceIcon class="icon"/>
        </div>
        <div class="name">
            {{ name }}
        </div>
    </div>
</template>
<script>
import DeviceIcon from '../assets/device.svg';

export default {
    name: 'Device',
    components: {
        DeviceIcon
    },
    props: {
        device: {
            type: Object,
            required: true
        }
    },
    methods: {
        connect() {
            if (!this.device.connecting) {
                this.$store.dispatch("ws_connectDevice", this.device.id);
            }
        }
    },
    computed: {
        name: function () {
            if (this.device.name.length > 20) {
                return this.device.name.substring(0, 17) + "...";
            }
            return this.device.name;
        },
        deviceStatus: function () {
            return this.device.connecting ? "connecting" : ""
        }
    }
};
</script>

<style scoped>
.device {
    cursor: pointer;
    display: inline-block;
    margin: 3em;
    padding: 1em;
    background: linear-gradient(rgba(0, 0, 0, 0.1));
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
}
</style>
