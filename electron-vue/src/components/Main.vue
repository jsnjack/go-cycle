<template>
    <div>
        <div class="cover">
            <video-bg
                :sources="['http://pespantelis.github.io/vue-videobg/demo/assets/timelapse.mp4']"
                img="http://pespantelis.github.io/vue-videobg/demo/assets/bg.jpg"
            >
                <Dashboard/>
            </video-bg>
        </div>
    </div>
</template>
<script>
import Dashboard from './Dashboard';
import VideoBg from 'vue-videobg'


export default {
    name: 'Main',

    components: {
        Dashboard, VideoBg
    },
    data() {
        return {};
    },
};
</script>

<style scoped>

</style>
