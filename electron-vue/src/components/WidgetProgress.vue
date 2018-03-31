<template>
    <div id="widget-progress">
            <div class="progress progress-slow" :style="opponentSlow"></div>
            <div class="progress progress-fast" :style="opponentFast"></div>
    </div>
</template>
<script>
import vuex from "vuex";

let your_color = "#3e92cc";
let opponent_color = "#d36135";

export default {
    name: 'WidgetProgress',
    computed: {
        ...vuex.mapGetters([
            "routeProgress",
        ]),
        ...vuex.mapState([
            "race",
        ]),
        opponentProgress: function () {
            if (this.race.opponents.length > 1) {
                return Math.round(this.race.opponents[1].distance / this.race.totalDistance * 100 * 100) / 100 || 0;
            }
            return 0;
        },
        progressDiff: function () {
            return this.routeProgress - this.opponentProgress;
        },
        opponentSlow: function() {
            let style;
            if (this.progressDiff >=0 ) {
                // User is faster
                style = {
                    width: this.opponentProgress + "%",
                    background: opponent_color,
                }
            } else {
                style = {
                    width: this.routeProgress + "%",
                    background: your_color,
                }
            }
            return style;
        },
        opponentFast: function () {
            let style = {
                width: Math.abs(this.progressDiff) + "%"
            };
            if (this.progressDiff >= 0) {
                // User is faster
                style.background = your_color;
            } else {
                style.background = opponent_color;
            }
            return style;
        },
    }
};
</script>

<style scoped>
    #widget-progress {
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 1em;
    }
    .progress {
        height: 1em;
        opacity: 0.7;
        float: left;
    }
</style>
