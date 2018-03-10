<template>
    <div v-show="distanceLeft < 1000" id="widget-finish">
        <div v-if="!isRaceFinished">
            {{ distanceLeft }}<span class="units">m</span>
        </div>
        <div v-else>
            FINISH &#x1f44f;
        </div>
    </div>
</template>

<script>
import vuex from "vuex";

export default {
    name: 'WidgetFinish',
    computed: {
        ...vuex.mapGetters([
            "distanceLeft",
            "isRaceFinished"
        ]),
    }
};
</script>

<style scoped>
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
</style>
