<template>
    <tr class="opponent">
        <td :style="nameStyle" class="name">{{ opponent.name }}</td>
        <td v-if="opponent.diffDistance != 0" class="diff">{{ opponent.diffDistance }} <span class="units">m</span></td>
    </tr>
</template>
<script>
import vuex from "vuex";

let yourColor = "#3e92cc";
let opponentColor = "#d36135";

export default {
    name: "Opponent",
    props: {
        opponent: {
            type: Object,
            required: true,
        },
    },
    computed: {
        ...vuex.mapState([
            "race",
        ]),
        nameStyle: function() {
            let style = {
                color: yourColor,
            };
            if (this.opponent.name !== "You") {
                style.color = opponentColor;
            }
            return style;
        },
    },
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
