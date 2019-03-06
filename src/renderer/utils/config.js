const electron = window.require("electron");
const path = window.require("path");
const fs = window.require("fs");
const userDataPath = (electron.app || electron.remote.app).getPath("userData");
const filePath = path.join(userDataPath, "config.json");


// Load configuration file
function loadConfig() {
    try {
        return JSON.parse(fs.readFileSync(filePath));
    } catch (error) {
        return {};
    }
}

// Save configuration file
function saveConfig(data={}) {
    fs.writeFile(filePath, JSON.stringify(data), "utf8", (err) => {
        if (err) {
            console.log(err);
            return;
        }
    });
}

export {loadConfig, saveConfig};
