<template>
    <div id="widget-progress">
            <div class="progress progress-slow" :style="opponentSlow"></div>
            <div class="progress progress-fast" :style="opponentFast"></div>
    </div>
</template>
<script>
import vuex from "vuex";

let yourColor = "#3e92cc";
let opponentColor = "#d36135";

export default {
    name: "WidgetProgress",
    computed: {
        ...vuex.mapGetters([
            "routeProgress",
        ]),
        ...vuex.mapState([
            "race",
        ]),
        opponentProgress: function() {
            if (this.race.opponents.length > 1) {
                return Math.round(this.race.opponents[1].distance / this.race.totalDistance * 100 * 100) / 100 || 0;
            }
            return 0;
        },
        progressDiff: function() {
            return this.routeProgress - this.opponentProgress;
        },
        opponentSlow: function() {
            let style;
            if (this.progressDiff >=0 ) {
                // User is faster
                style = {
                    width: this.opponentProgress + "%",
                    background: opponentColor,
                };
            } else {
                style = {
                    width: this.routeProgress + "%",
                    background: yourColor,
                };
            }
            return style;
        },
        opponentFast: function() {
            let style = {
                width: Math.abs(this.progressDiff) + "%",
            };
            if (this.progressDiff >= 0) {
                // User is faster
                style.background = yourColor;
            } else {
                style.background = opponentColor;
            }
            return style;
        },
    },
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
