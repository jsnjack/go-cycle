// https://thebikegeek.blogspot.nl/2009/12/while-we-wait-for-better-and-better.html
import CurveInterpolator from "curve-interpolator";

function _powerCurve(speed) {
    return 0.0115 * Math.pow(speed / 1.6, 3) - 0.0137* Math.pow(speed / 1.6, 2) + 8.9788 * speed / 1.6;
}

let data = [];
data.length = 70;

for (let i=0; i<data.length; i++) {
    data[i] = {x: i, y: _powerCurve(i)};
}

const trainer = new CurveInterpolator(data);

export default trainer;
