<template>
    <div id="widget-elevation">
    </div>
</template>
<script>
import vuex from "vuex";
import * as d3 from "d3";

const yourColor = "#3e92cc";
const opponentColor = "#d36135";

let markerYou;
let markerOpponent;
let path;

export default {
    name: "WidgetElevation",

    computed: {
        ...vuex.mapGetters(["routeProgress"]),
        ...vuex.mapState(["race"]),
        distanceYou: function() {
            return this.race.opponents[0].distance;
        },
        distanceOpponent: function() {
            return this.race.opponents[1].distance;
        },
    },

    mounted() {
        const chartMargins = {
            top: 10,
            right: 10,
            bottom: 10,
            left: 10,
        };
        const chartWidth = window.innerWidth;
        const chartHeight = 100;
        const width = chartWidth - chartMargins.right - chartMargins.left;
        const height = chartHeight - chartMargins.top - chartMargins.bottom;

        const svg = d3
            .select("#widget-elevation")
            .append("div")
            .append("svg-container", true)
            .append("svg")
            .attr("preserveAspectRatio", "xMinYMin meet")
            .attr("viewBox", "0 0 " + width + " " + height)
            .classed("svg-content-responsive", true);

        const g = svg.append("g");

        const xScale = d3
            .scaleLinear()
            .range([0, width])
            .domain(d3.extent(this.race.gpxData, (d) => d.distance));

        const yScale = d3
            .scaleLinear()
            .range([height, 0])
            .domain(d3.extent(this.race.gpxData, (d) => d.elevation));

        const areaFn = d3
            .area()
            .x((d) => xScale(d.distance))
            .y0(yScale(d3.min(this.race.gpxData, (d) => d.elevation)))
            .y1((d) => yScale(d.elevation))
            .curve(d3.curveLinear);

        const line = d3
            .line()
            .x(function(d) {
                return xScale(d.distance);
            })
            .y(function(d) {
                return yScale(d.elevation);
            })
            .curve(d3.curveLinear);

        // Area chart
        g
            .append("path")
            .datum(this.race.gpxData)
            .attr("fill", "lightgrey")
            .attr("d", areaFn)
            .node();

        // Line chart for calculating progress
        path = g
            .append("path")
            .datum(this.race.gpxData)
            .attr("fill", "none")
            .attr("d", line)
            .node();

        markerYou = svg
            .append("circle")
            .attr("id", "marker-you")
            .attr("r", 5)
            .attr("fill", yourColor)
            .attr("transform", "translate(0,0)");

        markerOpponent = svg
            .append("circle")
            .attr("id", "marker-opponent")
            .attr("r", 5)
            .attr("fill", opponentColor)
            .attr("transform", "translate(0,0)");
    },
    watch: {
        distanceYou: function(val) {
            if (markerYou && path) {
                let co = path.getPointAtLength(
                    val * path.getTotalLength() / this.race.totalDistance
                );
                markerYou.attr(
                    "transform",
                    "translate(" + co.x + "," + co.y + ")"
                );
            }
        },
        distanceOpponent: function(val) {
            if (markerOpponent && path) {
                let co = path.getPointAtLength(
                    val * path.getTotalLength() / this.race.totalDistance
                );
                markerOpponent.attr(
                    "transform",
                    "translate(" + co.x + "," + co.y + ")"
                );
            }
        },
    },
};
</script>

<style scoped>
#widget-elevation {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    opacity: 0.5;
}
.svg-container {
    display: inline-block;
    position: relative;
    width: 100%;
    padding-bottom: 100%; /* aspect ratio */
    vertical-align: top;
    overflow: hidden;
}
.svg-content-responsive {
    display: inline-block;
    position: absolute;
    top: 10px;
    left: 10px;
    right: 10px;
}
</style>
