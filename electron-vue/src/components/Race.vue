<template>
    <div>
        <div class="cover">
            <div class="video-container">
                <video
                    :src="getVideoFile"
                    autoplay="autoplay"
                    muted="true"
                    loop="true">
                </video>
                <WidgetDevices/>
                <WidgetProgress/>
            </div>
        </div>
    </div>
</template>
<script>
import vuex from "vuex";
import WidgetDevices from './WidgetDevices';
import WidgetProgress from './WidgetProgress';


export default {
    name: 'Race',

    components: {
        WidgetDevices, WidgetProgress
    },

    computed: {
        ...vuex.mapState([
            "race"
        ]),
        getVideoFile: function () {
            return [this.race.videoFile];
        }
    }
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
</style>
