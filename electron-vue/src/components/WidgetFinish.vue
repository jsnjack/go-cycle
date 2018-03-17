<template>
    <div v-show="race.startedAt && distanceLeft < 700 || isRaceFinished" id="widget-finish">
        <div v-if="!isRaceFinished">
            {{ distanceLeft }}<span class="units">m</span>
        </div>
        <div v-else-if="coolDown">
            <div class="cool-down-title">COOL DOWN</div>
            <div class="coll-down-timer">{{ collingDownTime }}</div>
            <div class="cool-down-button">
                <button class="button-control finish-button" @click="onFinish">Finish</button>
            </div>
        </div>
        <div v-else>
            FINISH &#x1f44f;
        </div>
    </div>
</template>

<script>
import vuex from "vuex";
import {formatTime} from "../utils/time";

export default {
    name: 'WidgetFinish',
    computed: {
        ...vuex.mapState([
            "race",
        ]),
        ...vuex.mapGetters([
            "distanceLeft",
            "isRaceFinished"
        ]),
        collingDownTime: function () {
            return formatTime(this.now - this.race.finishedAt);
        }
    },
    methods: {
        onFinish() {
            this.$router.push("afterrace");
        }
    },
    watch: {
        isRaceFinished: function (val) {
            if (val) {
                setInterval(()=> {
                    this.now = new Date();
                }, 1000);
                setTimeout(()=>{
                    this.coolDown = true;
                }, 10000);
            }
        }
    },
    data() {
        return{
            coolDown: false,
            now: new Date(),
        }
    }
};
</script>

<style scoped>
    @import url("../assets/style.css");

    #widget-finish {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        font-size: 10rem;
        font-weight: bold;
    }
    .units {
        font-size: 50%;
    }
    .cool-down-title {
        font-size: 25%;
        text-align: center;
    }
    .cool-down-title {
        text-align: center;
    }
    .cool-down-button {
        text-align: center;
        vertical-align: top;
    }
    .finish-button {
        background: transparent;
    }
</style>
