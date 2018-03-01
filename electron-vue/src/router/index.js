import Vue from "vue";
import Router from "vue-router";
import Connect from "@/components/Connect";
import Loading from "@/components/Loading";

Vue.use(Router);

export default new Router({
    routes: [
        {
            path: "/",
            name: "loading",
            component: Loading,
        },
        {
            path: "/connect",
            name: "connect",
            component: Connect,
        },
    ],
});
