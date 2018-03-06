import "xmldom";

const R = 6371000;

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

const utils = {
    readBlob, getDistance,
};

export default utils;
