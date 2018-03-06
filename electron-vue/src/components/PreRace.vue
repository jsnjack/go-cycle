<template>
    <div id="prerace">
        <h1>Settings</h1>
        <div class="container">
            <div class="section">
                <h2>Personal info</h2>
                <div class="row">
                    <div>Weight, kg</div>
                    <input type="number" :value="user.weight" @input="updateWeight"/>
                </div>
                <div class="row">
                    <div>Gender</div>
                    <select v-model="user.gender" @input="updateGender">
                        <option value="f">Female</option>
                        <option value="m">Male</option>
                    </select>
                </div>
                <div class="row">
                    <div>Age</div>
                    <input type="number" :value="user.age" @input="updateAge"/>
                </div>
                <div class="row">
                    <div>Tyre Size</div>
                    <!-- https://www.cateye.com/data/resources/Tire_size_chart_ENG.pdf -->
                    <select v-model="user.wheelSize" @input="updateWheelSize">
                        <option value="2070">700x18C</option>
                        <option value="2080">700x19C</option>
                        <option value="2086">700x20C</option>
                        <option value="2096">700x23C</option>
                        <option value="2105">700x25C</option>
                        <option value="2136">700x28C</option>
                        <option value="2146">700x30C</option>
                        <option value="2155">700x32C</option>
                    </select>
                </div>
            </div>

            <div class="section">
                <h2>Route</h2>
                    <div class="row">
                        <div>Video</div>
                        <input id="video" type="file" @change="saveFileReference"/>
                    </div>

                    <div class="row">
                        <div>GPX track</div>
                        <input id="gpx_track" type="file" @change="gpxTrack"/>
                        {{ routeDistance }}
                    </div>
            </div>
        </div>

        <div class="controls-container">
            <button class="button-control" @click="onBack">Back</button>
            <button class="button-control" @click="onStart">Start</button>
        </div>
    </div>
</template>
<script>
    import vuex from "vuex";
    import utils from "../utils/gpx";

    export default {
        name: 'PreRace',
        computed: {
                ...vuex.mapState([
                    "race",
                    "user",
                ]),
                ...vuex.mapGetters([
                    "routeDistance",
                ]),
        },
        methods: {
            onBack() {
                this.$router.push("connect");
            },
            onStart() {
                this.$router.push("race");
                this.$store.commit("START_RACE");
                this.$store.dispatch("ws_stopScanning");
            },
            saveFileReference(event) {
                let objectURL = window.URL.createObjectURL(event.target.files[0]);
                this.$store.commit("VIDEOFILE_URL", objectURL);
            },
            gpxTrack(event) {
                let objectURL = window.URL.createObjectURL(event.target.files[0]);
                utils.readBlob(event.target.files[0]).then(doc => {
                    this.$store.commit("SET_GPX_DOC", doc);
                });
            },
            updateWeight(event) {
                this.$store.commit("UPDATE_USER_WEIGHT", event.target.value);
            },
            updateGender(event) {
                this.$store.commit("UPDATE_USER_GENDER", event.target.value);
            },
            updateAge(event) {
                this.$store.commit("UPDATE_USER_AGE", event.target.value);
            },
            updateWheelSize(event) {
                this.$store.commit("UPDATE_USER_WHEEL_SIZE", event.target.value);
            }
        }
    };
</script>

<style scoped>
    @import url("../assets/style.css");

    #prerace {
        margin-top: 3rem;
        overflow: hidden;
    }

    .container {
        margin: 1em;
        font-size: 1.5rem;
        display: flex;
    }

    .section {
        flex: 1;
    }

    input, select {
        font-size: 1.5rem;
        border: none;
        width: 100px;
        margin-left: 2rem;
    }
    input[type='file'] {
        width: 100%;
    }
    input:focus{
        border: none;
    }
    .row {
        margin-bottom: 1rem;
    }

</style>
