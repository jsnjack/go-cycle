<template>
    <tr class="opponent">
        <td :style="nameStyle" class="name">{{ opponent.name }}</td>
        <td v-if="diffDistance" class="diff">{{ diffDistance }} <span class="units">m</span></td>
    </tr>
</template>
<script>
import vuex from "vuex";

let your_color = "#3e92cc";
let opponent_color = "#d36135";

export default {
    name: 'Opponent',
    props: {
        opponent: {
            type: Object,
            required: true
        }
    },
    computed: {
        ...vuex.mapState([
            "race",
        ]),
        nameStyle: function () {
            let style = {
                color: your_color
            };
            if (this.opponent.name !== "You") {
                style.color = opponent_color;
            }
            return style;
        },
        diffDistance: function() {
            let diff = this.race.opponents[0].distance - this.opponent.distance;
            return -Math.round(diff * 10) / 10;
        }
    }
};
</script>

<style scoped>
.opponent {
    font-size: 1.5rem;
}
.name {
    padding-right: 1rem;
}
.units {
    font-size: 50%;
}
</style>
