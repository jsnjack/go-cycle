<template>
    <div>
        <div class="cover">
            <div class="video-container">
                <video
                    :src="getVideoFile"
                    autoplay="autoplay"
                    muted="true"
                    loop="true"
                    width="100%"
                    height="100%">
                </video>
                <WidgetDevices/>
                <WidgetProgress v-if="race.gpxData.length === 0"/>
                <WidgetFinish/>
                <WidgetElevation v-if="race.gpxData.length > 0"/>
                <WidgetOpponents/>
                <WidgetWarmup/>
            </div>
        </div>
    </div>
</template>
<script>
import vuex from "vuex";
import WidgetDevices from "./WidgetDevices";
import WidgetProgress from "./WidgetProgress";
import WidgetFinish from "./WidgetFinish";
import WidgetElevation from "./WidgetElevation";
import WidgetOpponents from "./WidgetOpponents";
import WidgetWarmup from "./WidgetWarmup";

export default {
    name: "Race",

    components: {
        WidgetDevices, WidgetProgress, WidgetFinish, WidgetElevation, WidgetWarmup,
        WidgetOpponents,
    },

    computed: {
        ...vuex.mapState([
            "race",
        ]),
        ...vuex.mapGetters([
            "isRaceInProgress",
        ]),
        getVideoFile: function() {
            if (this.race.startedAt) {
                return [this.race.videoFile];
            }
        },
    },
};
</script>

<style scoped>
    .video-container {
        position: fixed;
        right: 0;
        bottom: 0;
        min-width: 100%;
        min-height: 100%;
        background: url("../assets/background.jpg") center;
        background-size: cover;
    }
    .video-container video {
        position: fixed;
        min-width: 100%;
        min-height: 100%;
        object-fit: cover;
    }
</style>
