// https://www.gribble.org/cycling/power_v_speed.html

const g = 9.8067; // gravitational constant, m/s^2
const a = 0.509; // Frontal area, m^2
const cd = 0.63; // Drag coeff
const crr = 0.005; // Coeff rolling resistence
const rho = 1.226; // Air resistence

// speed from the sensor m/s
function speedFromSensor(revolutions, time, circumference) {
    let speed = revolutions * circumference / time;
    return speed;
}

// Speed from ms to kmh
function toKMH(speed) {
    let kmh = speed * 3600 / 1000;
    return kmh;
}

// Speed from kmh to ms
function toMS(speed) {
    let ms = speed * 1000 / 3600;
    return ms;
}

// Get real speed
function getRealSpeed(speedFromSensor, grade, mass) {
    let realSpeed = 0;
    let speedMS = toMS(speedFromSensor);
    let fakePower = fRolling(0, mass) * speedMS + dragCoeff() * speedMS * speedMS * speedMS;

    let aa = dragCoeff();
    let bb = 0;
    let cc = fRolling(grade, mass) + fGravity(grade, mass);
    let dd = -fakePower;
    let roots = [0];
    try {
        roots = solveCubic(aa, bb, cc, dd);
    } catch (e) {
    };
    realSpeed = toKMH(getMinPositive(roots));
    return realSpeed;
}

// Rolling resistence force
function fRolling(grade, mass) {
    let f = g * Math.cos(Math.atan(grade)) * mass * crr;
    return f;
}

// Gravity force
function fGravity(grade, mass) {
    let f = g * Math.sin(Math.atan(grade)) * mass;
    return f;
}

// Draf coefficient, Fdrag = dargCoef*speed^2
function dragCoeff() {
    return 0.5 * cd * a * rho;
}

function cuberoot(x) {
    let y = Math.pow(Math.abs(x), 1/3);
    return x < 0 ? -y : y;
}

function solveCubic(a, b, c, d) {
    if (Math.abs(a) < 1e-8) { // Quadratic case, ax^2+bx+c=0
        a = b; b = c; c = d;
        if (Math.abs(a) < 1e-8) { // Linear case, ax+b=0
            a = b; b = c;
            if (Math.abs(a) < 1e-8) { // Degenerate case
                return [];
            }
            return [-b/a];
        }

        let D = b*b - 4*a*c;
        if (Math.abs(D) < 1e-8) {
            return [-b/(2*a)];
        } else if (D > 0) {
            return [(-b+Math.sqrt(D))/(2*a), (-b-Math.sqrt(D))/(2*a)];
        }
        return [];
    }

    // Convert to depressed cubic t^3+pt+q = 0 (subst x = t - b/3a)
    let p = (3*a*c - b*b)/(3*a*a);
    let q = (2*b*b*b - 9*a*b*c + 27*a*a*d)/(27*a*a*a);
    let roots;

    if (Math.abs(p) < 1e-8) { // p = 0 -> t^3 = -q -> t = -q^1/3
        roots = [cuberoot(-q)];
    } else if (Math.abs(q) < 1e-8) { // q = 0 -> t^3 + pt = 0 -> t(t^2+p)=0
        roots = [0].concat(p < 0 ? [Math.sqrt(-p), -Math.sqrt(-p)] : []);
    } else {
        let D = q*q/4 + p*p*p/27;
        if (Math.abs(D) < 1e-8) { // D = 0 -> two roots
            roots = [-1.5*q/p, 3*q/p];
        } else if (D > 0) { // Only one real root
            let u = cuberoot(-q/2 - Math.sqrt(D));
            roots = [u - p/(3*u)];
        } else { // D < 0, three roots, but needs to use complex numbers/trigonometric solution
            let u = 2*Math.sqrt(-p/3);
            let t = Math.acos(3*q/p/u)/3; // D < 0 implies p < 0 and acos argument in [-1..1]
            let k = 2*Math.PI/3;
            roots = [u*Math.cos(t), u*Math.cos(t-k), u*Math.cos(t-2*k)];
        }
    }

    // Convert back from depressed cubic
    for (let i = 0; i < roots.length; i++) {
        roots[i] -= b/(3*a);
    }
    return roots;
}

// Return minimum positive value from list
function getMinPositive(a) {
    let speed = 0;
    for (let i=0; i<a.length; i++) {
        if (a[i] > 0 && a[i] > speed) {
            speed = a[i];
        }
    }
    return speed;
}

const real = {
    speedFromSensor, toKMH, getRealSpeed, toMS,
};

export default real;
