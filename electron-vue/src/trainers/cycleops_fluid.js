// https://thebikegeek.blogspot.nl/2009/12/while-we-wait-for-better-and-better.html
function powerCurve(speed) {
    return 0.0115 * Math.pow(speed / 1.6, 3) - 0.0137* Math.pow(speed / 1.6, 2) + 8.9788 * speed / 1.6;
}

export default powerCurve;
