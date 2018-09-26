const installer = require("electron-installer-redhat");
const rpmConfig = require("./rpm");


function createRPM() {
    console.log("Creating rpm package...", rpmConfig.version);

    installer(rpmConfig, function(err) {
        if (err) {
            console.error(err, err.stack);
            process.exit(1);
        }

        console.log("Successfully created package at " + rpmConfig.dest);
    });
}

createRPM();
