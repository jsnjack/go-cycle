<template>
    <div id="widget-progress">
        <div class="progress" :style="progressStyle"></div>
    </div>
</template>
<script>
import vuex from "vuex";

export default {
    name: 'WidgetProgress',
    computed: {
        ...vuex.mapGetters([
            "routeProgress",
        ]),
        progressStyle: function() {
            return {
                width: this.routeProgress + "%"
            }
        }
    }
};
</script>

<style scoped>
    #widget-progress {
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
    }
    .progress {
        opacity: 0.5;
        background: #3e92cc;
        text-align: right;
        height: 0.5em;
    }
</style>
