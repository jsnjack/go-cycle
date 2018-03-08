<template>
    <div id="afterrace">
        <h1>Results</h1>
        <table>
            <tr>
                <td>Distance</td>
                <td>{{ totalDistance }} km</td>
            </tr>
            <tr>
                <td>Time</td>
                <td>{{ duration }}</td>
            </tr>
            <tr>
                <td>Average speed</td>
                <td>{{ averageSpeed }} km/h</td>
            </tr>
            <tr>
                <td>Max speed</td>
                <td>{{ maxSpeed }} km/h</td>
            </tr>
            <tr>
                <td>Calories</td>
                <td>{{ race.calories }} kcal</td>
            </tr>
        </table>

        <div class="controls-container">
            <button class="button-control" @click="onNew">New race</button>
        </div>
    </div>
</template>
<script>
import vuex from "vuex";
import {formatTime} from "../utils/time";

export default {
    name: 'AfterRace',
    computed: {
        ...vuex.mapGetters([
            "distance",
            "maxSpeed",
        ]),
        ...vuex.mapState([
            "race",
        ]),
        duration: function () {
            return formatTime(this.race.finishedAt - this.race.startedAt);
        },
        averageSpeed: function () {
            let avg = Math.round(this.race.distance / this.duration * 3.6 * 1000 * 1000 * 10) / 10 || 0;
            return avg;
        },
        totalDistance: function () {
            // Total distance in km
            return Math.round(this.race.distance / 1000 * 10) / 10 || 0;
        }
    },
    methods: {
        onNew: function () {
            this.$router.push("prerace");
        }
    }
};
</script>

<style scoped>
    @import url("../assets/style.css");
    #afterrace {
        margin-top: 3rem;
    }
    table {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        font-size: 1.5rem;
        opacity: 0.9;
    }
    td {
        padding: 1rem 2rem;
    }
    tr td:first-child {
        text-align: right;
    }
    tr td:last-child {
        text-align: left;
    }
</style>
