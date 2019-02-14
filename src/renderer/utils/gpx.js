const R = 6371000;
// Check Exponential Moving Average https://docs.oracle.com/cd/E57185_01/IRWUG/ch12s07s05.html

/* eslint-disable max-len */
// A template for creating a gpx file
const gpxTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<gpx creator="Go-cycle" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd">
 <metadata>
  <time>2018-03-10T11:19:44Z</time>
 </metadata>
 <trk>
  <name>Go-cycle activity</name>
  <trkseg>
  </trkseg>
 </trk>
</gpx>
`;
/* eslint-enable max-len */

function toRad(n) {
    return n * Math.PI / 180;
}

// Calculate distance assuming the Earth is spherical
function getDistanceHaversine(lat1, lon1, lat2, lon2) {
    let x1 = lat2 - lat1;
    let dLat = toRad(x1);
    let x2 = lon2 - lon1;
    let dLon = toRad(x2);
    let a = Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(toRad(lat1)) * Math.cos(toRad(lat2)) *
    Math.sin(dLon / 2) * Math.sin(dLon / 2);
    let c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    let d = R * c;
    return d;
}

const getDistance = getDistanceHaversine;

// Read imported gpx file
function readBlob(blob) {
    return new Promise((resolve, reject) => {
        let reader = new FileReader();
        reader.readAsText(blob);
        reader.onload = function() {
            let doc = new DOMParser().parseFromString(
                reader.result,
                "text/xml"
            );
            resolve(doc);
        };
    });
}

// Create gpx file out of recorded data
function createGPX(points, startedAt, gpxData) {
    // More info https://developers.strava.com/docs/uploads/
    let started = performance.now();
    let offset = 1;
    let DOMParser = require("xmldom").DOMParser;
    let doc = new DOMParser().parseFromString(gpxTemplate, "text/xml");
    doc.getElementsByTagName("time")[0].textContent = startedAt;
    let trkseg = doc.getElementsByTagName("trkseg")[0];
    for (let i=0; i<points; i++) {
        let data = localStorage.getItem("trkpt_" + i);
        if (data) {
            let dataObj = JSON.parse(data);
            let trkpt = doc.createElement("trkpt");
            let time = doc.createElement("time");
            time.textContent = dataObj.time;
            let extensions = doc.createElement("extensions");
            let distance = doc.createElement("distance");
            distance.textContent = dataObj.distance;
            if (gpxData.length > 0) {
                let P = getCoordinatesFromDistance(dataObj.distance, gpxData, offset);
                if (P[0] && P[1]) {
                    trkpt.setAttribute("lat", P[0]);
                    trkpt.setAttribute("lon", P[1]);
                }
                offset = P[2];
            }
            if (dataObj.hr) {
                let hr = doc.createElement("heartrate");
                hr.textContent = dataObj.hr;
                extensions.appendChild(hr);
            }

            if (dataObj.cadence) {
                let cadence = doc.createElement("cadence");
                cadence.textContent = dataObj.cadence;
                extensions.appendChild(cadence);
            }

            let power = doc.createElement("power");
            power.textContent = dataObj.power;
            extensions.appendChild(power);

            extensions.appendChild(distance);

            trkpt.appendChild(time);
            trkpt.appendChild(extensions);
            trkseg.appendChild(trkpt);
        }
    }
    let XMLSerializer = require("xmldom").XMLSerializer;
    let serialized = new XMLSerializer().serializeToString(doc);
    let finished = performance.now();
    console.debug("Creating GPX took, ms:", finished - started);
    return serialized;
}


function getCoordinatesFromDistance(distance, gpxData, offset) {
    for (let i=offset; i<gpxData.length; i++) {
        if (gpxData[i].distance > distance) {
            let D = gpxData[i].distance - gpxData[i-1].distance;
            let d = distance - gpxData[i-1].distance;
            return [
                gpxData[i-1].lat + d / D * (gpxData[i].lat - gpxData[i-1].lat),
                gpxData[i-1].lon + d / D * (gpxData[i].lon - gpxData[i-1].lon),
                i-1,
            ];
        }
    }
    return [null, null, 0];
}

// Extracts data from gpx file, like distance
function extractDataFromGPX(doc) {
    let started = performance.now();
    let distance = 0;
    let points = doc.getElementsByTagName("trkpt");
    let container = new Array(points.length);
    let baseTime;
    if (points[0].getElementsByTagName("time").length) {
        baseTime = new Date(points[0].getElementsByTagName("time")[0].textContent).getTime();
    }
    container[0] = {
        distance: 0,
        grade: 0,
        elevation: parseFloat(points[0].getElementsByTagName("ele")[0].textContent),
        lat: parseFloat(points[0].getAttribute("lat")),
        lon: parseFloat(points[0].getAttribute("lon")),
    };
    if (baseTime) {
        container[0].time = 0;
    }
    for (let i=0; i<points.length - 1; i++) {
        let fragment = utils.getDistance(
            points[i].getAttribute("lat"),
            points[i].getAttribute("lon"),
            points[i+1].getAttribute("lat"),
            points[i+1].getAttribute("lon")
        );
        distance += fragment;
        let ele = parseFloat(points[i+1].getElementsByTagName("ele")[0].textContent);
        container[i+1] = {
            distance: distance,
            lat: parseFloat(points[i+1].getAttribute("lat")),
            lon: parseFloat(points[i+1].getAttribute("lon")),
            elevation: ele,
            grade: 0,
        };
        if (baseTime) {
            container[i+1].time = new Date(
                points[i+1].getElementsByTagName("time")[0].textContent
            ).getTime() - baseTime;
        }
    }
    container = calculateGrade(smoothMovingAverage(container));
    let finished = performance.now();
    console.debug("Extracting data from GPX took, ms:", finished - started);
    return container;
}

// Smootherns elevation data
function smoothMovingAverage(container, sampleSize=10) {
    for (let i=0; i<container.length - 1; i++) {
        let sum = 0;
        if (i < container.length - 1 - sampleSize) {
            for (let j=i; j < i + sampleSize; j++) {
                sum = sum + container[j].elevation;
            }
        } else {
            for (let j=container.length - 1 - sampleSize; j < container.length - 1; j++) {
                sum = sum + container[j].elevation;
            }
        }
        let ma = sum / sampleSize;
        container[i].elevation = ma;
    }
    return container;
}

// Calculate grade
function calculateGrade(container) {
    for (let i=1; i<container.length - 1; i++) {
        let fragment = container[i].distance - container[i-1].distance;
        let grade = (container[i].elevation - container[i-1].elevation) / fragment;
        if (grade === Infinity) {
            grade = 0;
        }
        container[i].grade = grade;
    }
    return container;
}

const utils = {
    readBlob, getDistance, createGPX, extractDataFromGPX,
};

export default utils;
