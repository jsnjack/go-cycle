<template>
    <div id="widget-opponents" v-show="race.opponents.length > 1">
        <table>
            <Opponent v-for="item in orderedOpponents"
                :key="item.name"
                :opponent="item">
            </Opponent>
        </table>
    </div>
</template>
<script>
import vuex from "vuex";
import Opponent from "./Opponent";
import _ from "lodash";


export default {
    name: "WidgetOpponents",
    components: {
        Opponent,
    },
    computed: {
        ...vuex.mapState([
            "race",
        ]),
        orderedOpponents: function() {
            return _.reverse(_.orderBy(this.race.opponents, ["distance"]));
        },
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
