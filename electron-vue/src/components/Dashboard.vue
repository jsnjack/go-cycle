<template>
    <div>
    <h1>go-cycle</h1>
    <table>
        <tr>
            <td>
                <HRIcon class="icon"/>
            </td>
            <td>
                {{ hrBPM }} bpm
            </td>
        </tr>
        <tr>
            <td>
                <SpeedIcon class="icon"/>
            </td>
            <td>
                {{ speed }}
            </td>
        </tr>
        <tr>
            <td>
                <DistanceIcon class="icon"/>
            </td>
            <td>
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
    computed: {
        speed: function() {
            let value = this.wheel / this.cscRevPerSec;
            return `${value} km/h`;
        },
        distance: function() {
            let value = this.wheel * this.cscRevolutions;
            return `${value} m`;
        },
    },
    data() {
        return {
            wheel: 100,
            hrBPM: 0,
            cscRevolutions: 0,
            cscRevPerSec: 0,
        };
    },
};
</script>

<style scoped>
.icon {
    height: 5rem;
    width: 5rem;
}
</style>

