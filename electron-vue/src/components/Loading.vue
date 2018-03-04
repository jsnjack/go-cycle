<template>
    <div id="loading">
        Connecting...
    </div>
</template>
<script>
import vuex from "vuex";

export default {
    name: 'Loading',
};
</script>

<style scoped>
    #loading {
        display: flex;
        flex-direction: column;
        justify-content: center;
        align-items: center;
        text-align: center;
        min-height: 80vh;
        font-size: 200%;
        text-decoration: none;
        color: white;
    }
</style>
