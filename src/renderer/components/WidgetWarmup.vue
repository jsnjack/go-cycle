<template>
    <div v-show="!race.startedAt && user.warmUpDuration" id="widget-warmup">
        <div class="warmup-title">WARMUP</div>
        <div class="warmup-left"> {{ warmupTime }}</div>
        <div class="warmup-button">
            <button class="button-control finish-button" @click="onWarmupFinish">Start</button>
        </div>
    </div>
</template>

<script>
import vuex from "vuex";
import {formatTime} from "../utils/time";

export default {
    name: "WidgetWarmup",
    computed: {
        ...vuex.mapState(["race", "user"]),
        warmupTimeLeft: function() {
            let warmupLeft =
                this.now -
                this.warmupStartedAt -
                this.user.warmUpDuration * 1000;
            if (warmupLeft > 0) {
                this.onWarmupFinish();
            } else {
                this.schedule();
            }
            return warmupLeft;
        },
        warmupTime: function() {
            return formatTime(Math.abs(this.warmupTimeLeft));
        },
    },
    methods: {
        onWarmupFinish() {
            if (!this.race.finishedAt && !this.warmupFinished) {
                this.warmupFinished = true;
                this.$store.commit("START_RACE");
                this.$store.dispatch("update_opponent", 0);
            }
        },
        schedule() {
            setTimeout(() => {
                this.now = new Date().getTime();
            }, 1000);
        },
    },
    data() {
        return {
            warmupStartedAt: new Date().getTime(),
            warmupFinished: false,
            now: new Date().getTime(),
        };
    },
};
</script>

<style scoped>
@import url("../assets/style.css");

#widget-warmup {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    font-size: 10rem;
    font-weight: bold;
}
.warmup-title {
    font-size: 25%;
    text-align: center;
}
.warmup-button {
    text-align: center;
    vertical-align: top;
}
.finish-button {
    background: transparent;
}
</style>
