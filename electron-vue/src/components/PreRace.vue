<template>
    <div class="container">
        <h1>Settings</h1>
        <section>
            <h2>Body</h2>
            <div class="row">
                <span>Weight, kg</span>
                <input type="number" :value="user.weight"/>
            </div>
            <div class="row">
                <span>Gender</span>
                <input type="text" :value="user.gender"/>
            </div>
            <div class="row">
                <span>Age</span>
                <input type="number" :value="user.age"/>
            </div>
        </section>

        <section>
            <h2>Bicycle</h2>
                <span>Wheel Size, mm</span>
                <input type="number" :value="user.wheelSize"/>
        </section>

        <section>
            <h2>Video</h2>
                <span>Location</span>
                <input id="video" type="file" @change="saveFileReference"/>
        </section>

        <div class="controls-container">
            <button class="button-control" @click="onBack">Back</button>
            <button class="button-control" @click="onStart">Start</button>
        </div>
    </div>
</template>
<script>
    import vuex from "vuex";

    export default {
        name: 'PreRace',
        computed: {
                ...vuex.mapState([
                    "race",
                    "user",
                ]),
        },
        methods: {
            onBack() {
                this.$router.push("connect");
            },
            onStart() {
                this.$router.push("race");
                this.$store.commit("START_RACE");
            },
            saveFileReference(event) {
                let objectURL = window.URL.createObjectURL(event.target.files[0]);
                this.$store.commit("VIDEOFILE_URL", objectURL);
            },
        }
    };
</script>

<style scoped>
    @import url("../assets/style.css");

    .container {
        margin: 1em;
        font-size: 1.5rem;
    }

    input {
        font-size: 1.5rem;
        background: transparent;
        color: white;
        border: 2px solid white;
        width: 100px;
    }
    input[type='file'] {
        width: 100%;
    }
    input:focus{
        border: 2px solid white;
    }

</style>
