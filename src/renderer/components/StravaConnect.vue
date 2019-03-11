<template>
  <webview
    :src="getSRC"
    @did-navigate="hnadleNavigation"
    ></webview>
</template>

<script>

const registrationSuccessURL = "https://gocycle.space/register/success?authid=";

export default {
    name: "StravaConnect",
    computed: {
        getSRC: function() {
            let url = "https://www.strava.com/oauth/authorize?client_id=24045" +
                      "&redirect_uri=https://gocycle.space/register&" +
                      "response_type=code&approval_scope=auto&scope=activity:write";
            return url;
        },
    },
    methods: {
        hnadleNavigation(ev) {
            if (ev.httpResponseCode === 200 && ev.url.startsWith(registrationSuccessURL)) {
                console.log("go-cycle-auth: success");
                this.$store.commit("UPDATE_USER", {
                    "stravaAccessToken": ev.url.substring(registrationSuccessURL.length),
                });
                this.$router.push("prerace");
            }
        },
    },
};
</script>

<style scoped>
webview {
  display: flex;
  height: 100vh;
  flex-grow: 1;
}
</style>
