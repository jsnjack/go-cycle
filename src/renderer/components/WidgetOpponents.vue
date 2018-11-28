<template>
    <div id="widget-opponents" v-show="participants.length > 1">
        <table>
            <Opponent v-for="item in participants"
                :key="item.name"
                :opponent="item">
            </Opponent>
        </table>
    </div>
</template>
<script>
import vuex from "vuex";
import Opponent from "./Opponent";
const updateTablePeriod = 2000;


export default {
    name: "WidgetOpponents",
    components: {
        Opponent,
    },
    created() {
        for (let i=0; i<this.race.opponents.length; i++) {
            this.participants.push({
                name: this.race.opponents[i].name,
                diffDistance: 0,
            });
        }
        this.interval = setInterval(function() {
            this.updateDistance();
        }.bind(this), updateTablePeriod);
    },
    beforeDestroy: function() {
        clearInterval(this.interval);
    },
    methods: {
        updateDistance: function() {
            let diff = Math.round(this.race.opponents[1].distance - this.race.opponents[0].distance);
            if (diff > 0) {
                this.participants[0].diffDistance = diff;
                this.participants[0].name = this.race.opponents[1].name;
                this.participants[1].diffDistance = 0;
                this.participants[1].name = this.race.opponents[0].name;
            } else {
                this.participants[0].diffDistance = 0;
                this.participants[0].name = this.race.opponents[0].name;
                this.participants[1].diffDistance = diff;
                this.participants[1].name = this.race.opponents[1].name;
            }
        },
    },
    computed: {
        ...vuex.mapState([
            "race",
        ]),
    },
    data() {
        return {
            participants: [],
            interval: null,
        };
    },
};
</script>

<style scoped>
    #widget-opponents {
        position: fixed;
        top: 29px;
        right: 0;
    }
    table {
        opacity: 0.5;
        background: rgb(36, 36, 36);
        padding: 1rem;
    }
</style>
