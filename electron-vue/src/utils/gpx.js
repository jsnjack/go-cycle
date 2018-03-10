const R = 6371000;

/* eslint-disable max-len */
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

function createGPX(points, distPerRev) {
    let DOMParser = require("xmldom").DOMParser;
    let doc = new DOMParser().parseFromString(gpxTemplate, "text/xml");
    doc.getElementsByTagName("time")[0].textContent = new Date().toISOString();
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
            distance.textContent = dataObj.rev * distPerRev;
            if (dataObj.hr) {
                let hr = doc.createElement("heartrate");
                hr.textContent = dataObj.hr;
                extensions.appendChild(hr);
            }
            extensions.appendChild(distance);
            trkpt.appendChild(time);
            trkpt.appendChild(extensions);
            trkseg.appendChild(trkpt);
        }
    }
    let XMLSerializer = require("xmldom").XMLSerializer;
    return new XMLSerializer().serializeToString(doc);
}

const utils = {
    readBlob, getDistance, createGPX,
};

export default utils;
