<template>
    <div id="dashboard">
        <table>
            <tr>
                <td>
                    <HRIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ hrBPM }} bpm
                </td>
            </tr>
            <tr>
                <td>
                    <SpeedIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ speed }}
                </td>
            </tr>
            <tr>
                <td>
                    <DistanceIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ distance }}
                </td>
            </tr>
        </table>
    </div>
</template>
<script>
import HRIcon from '../assets/heart-rate.svg';
import DistanceIcon from '../assets/distance.svg';
import SpeedIcon from '../assets/speed-meter.svg';

export default {
    name: 'Dashboard',

    components: {
        HRIcon, DistanceIcon, SpeedIcon,
    },

    created() {
        this.connect();
    },

    computed: {
        speed: function() {
            let value = this.cscRevPerSec * this.distancePerRev * 3.6 * 1000 * 1000;
            value = value.toFixed(1);
            return `${value} km/h`;
        },
        distance: function() {
            let value = this.distancePerRev * this.cscRevolutions;
            value = Math.round(value);
            return `${value} m`;
        },
    },

    methods: {
        connect() {
            console.log("this", this);
            console.log('Connecting...');
            let ws = new WebSocket('ws://localhost:8000/ws');

            // Listen for messages
            ws.addEventListener('message', event => {
                this.handleWSMessage(JSON.parse(event.data));
            });

            ws.addEventListener('close', event => {
                console.warn('WS connection closed', event);
                setTimeout(this.connect, 1000);
            });

            ws.addEventListener('error', event => {
                console.warn('WS connection error', event);
            });
        },
        handleWSMessage(msg) {
            if (msg.type === "hr_data") {
                this.hrBPM = msg.data.bpm;
            } else if (msg.type === "csc_data") {
                this.cscRevolutions += msg.data.revolutions;
                this.cscRevPerSec = msg.data.rev_per_sec;
            } else {
                console.info("Unhandled message", msg);
            }
        },
    },

    data() {
        return {
            distancePerRev: Math.PI * (622 + 28*2) / 1000, // meters
            hrBPM: 0,
            cscRevolutions: 0,
            cscRevPerSec: 0,
        };
    },
};
</script>

<style scoped>
table {
    opacity: 0.5;
    background: rgb(36, 36, 36);
    padding: 2rem;
}
.icon {
    height: 5rem;
    width: 5rem;
    margin: 0.25rem 2rem 0.25rem 0.25rem;
    fill: white;
}

.measurement {
    vertical-align: middle;
    font-size: 3rem;
    color: white;
}

</style>

