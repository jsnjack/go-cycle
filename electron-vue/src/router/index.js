import Vue from "vue";
import Router from "vue-router";
import Connect from "@/components/Connect";
import Loading from "@/components/Loading";
import PreRace from "@/components/PreRace";
import Race from "@/components/Race";


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
        {
            path: "/prerace",
            name: "prerace",
            component: PreRace,
        },
        {
            path: "/race",
            name: "race",
            component: Race,
        },
    ],
});
