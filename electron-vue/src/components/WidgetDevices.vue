<template>
    <div id="widget-devices">
        <table>
            <tr>
                <td>
                    <HRIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ race.currentBPM }} <span class="unit">bpm</span>
                </td>
            </tr>
            <tr>
                <td>
                    <SpeedIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ getCurrentSpeed }} <span class="unit">km/h</span>
                </td>
            </tr>
            <tr>
                <td>
                    <DistanceIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ getDistance }} <span class="unit">{{ getDistanceUnit }}</span>
                </td>
            </tr>
            <tr>
                <td>
                    <PowerIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ getPower }} <span class="unit">W</span>
                </td>
            </tr>
            <tr>
                <td>
                    <MountainIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ getGrade }} <span class="unit">%</span>
                </td>
            </tr>
            <tr>
                <td>
                    <CaloriesIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ getCalories }} <span class="unit">kCal</span>
                </td>
            </tr>
            <tr>
                <td>
                    <TimeIcon class="icon"/>
                </td>
                <td class="measurement">
                    {{ getRaceDuration }}
                </td>
            </tr>
        </table>
    </div>
</template>
<script>
import vuex from "vuex";
import HRIcon from '../assets/hr-connected.svg';
import MountainIcon from '../assets/mountain.svg';
import DistanceIcon from '../assets/distance.svg';
import SpeedIcon from '../assets/speed-meter.svg';
import CaloriesIcon from '../assets/fire.svg';
import TimeIcon from '../assets/time.svg';
import PowerIcon from '../assets/power.svg';
import {formatTime} from "../utils/time";


export default {
    name: 'WidgetDevices',

    components: {
        HRIcon, DistanceIcon, SpeedIcon, CaloriesIcon, TimeIcon, PowerIcon, MountainIcon
    },

    mounted () {
        setInterval(()=> {
            this.now = new Date();
        }, 1000);
    },

    computed: {
        ...vuex.mapState([
            "race"
        ]),
        ...vuex.mapGetters([
            "isRaceInProgress"
        ]),
        getDistance: function () {
            if (this.race.distance < 1000) {
                return Math.round(this.race.distance);
            }
            return Math.round(this.race.distance / 1000 * 10) / 10;
        },
        getDistanceUnit: function () {
            if (this.race.distance < 1000) {
                return "m";
            }
            return "km"
        },
        getCalories: function () {
            return Math.round(this.race.calories);
        },
        getRaceDuration: function () {
            if (this.race.startedAt) {
                return formatTime(this.now.getTime() - this.race.startedAt.getTime());
            }
            return 0;
        },
        getCurrentSpeed: function () {
            if(this.isRaceInProgress) {
                return this.race.speed.toFixed(1);
            }
            return this.race.csc.speed.toFixed(1);
        },
        getPower: function () {
            return Math.round(this.race.currentPower);
        },
        getGrade: function () {
            return Math.round(this.race.grade * 100 * 10) / 10 || 0.0;
        },
    },
    data () {
        return {
            now: new Date(),
        };
    }
};
</script>

<style scoped>
    #widget-devices {
        position: fixed;
        top: 29px;
        left: 0;
    }
    table {
        opacity: 0.5;
        background: rgb(36, 36, 36);
        padding: 1rem;
    }
    .icon {
        height: 2rem;
        width: 2rem;
        margin: 0.25rem 1rem 0.25rem 0.25rem;
        fill: white;
    }
    .measurement {
        vertical-align: middle;
        font-size: 2rem;
        color: white;
    }
    .measurement .unit {
        vertical-align: bottom;
        font-size: 50%;
        color: white;
    }
</style>

