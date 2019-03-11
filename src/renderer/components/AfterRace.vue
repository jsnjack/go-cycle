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
                <td>{{ calories }} kcal</td>
            </tr>
        </table>

        <div class="controls-container">
            <button v-show="!!user.stravaAccessToken && !uploaded"
                    class="button-control" @click="onUpload">Upload activity</button>
            <button class="button-control" @click="onNew">New race</button>
        </div>
    </div>
</template>
<script>
import vuex from "vuex";
import {formatTime} from "../utils/time";
import utils from "../utils/gpx";

let fs = require("fs");
const activityUploadURL = "https://gocycle.space/upload";

export default {
    name: "AfterRace",
    computed: {
        ...vuex.mapState([
            "race",
            "user",
        ]),
        duration: function() {
            return formatTime(this.race.finishedAt - this.race.startedAt);
        },
        averageSpeed: function() {
            let avg = Math.round(
                this.race.distance / (this.race.finishedAt - this.race.startedAt) * 3.6 * 1000 * 10
            ) / 10 || 0;
            return avg;
        },
        maxSpeed: function() {
            return Math.round(this.race.maxSpeed * 10) / 10 || 0;
        },
        totalDistance: function() {
            // Total distance in km
            return Math.round(this.race.distance / 1000 * 10) / 10 || 0;
        },
        calories: function() {
            return Math.round(this.race.calories);
        },
    },
    methods: {
        onNew: function() {
            this.$router.push("prerace");
        },
        onUpload: function() {
            let gpxData = utils.createGPX(this.race.points, this.race.startedAt.toISOString(), this.race.gpxData);
            let formData = new FormData();
            let filePath = new Date().getTime() + ".gpx";

            fs.writeFile(filePath, gpxData, "utf8", (err) => {
                if (err) {
                    console.log(err);
                    return;
                }
                console.log(`File ${filePath} created`);
            });

            formData.append("authid", this.user.stravaAccessToken);
            formData.append("file", new File([gpxData], "activity.gpx", {type: "text/xml"}));
            fetch(activityUploadURL, {
                method: "POST",
                body: formData,
            }).then((response) => response.json())
                .catch((error) => console.error("Error:", error))
                .then((response) => {
                    console.log("Success:", response);
                    this.uploaded = true;
                });
        },
    },
    data() {
        return {
            uploaded: false,
        };
    },
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
