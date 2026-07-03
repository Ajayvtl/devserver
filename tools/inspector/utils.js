const fs = require("fs");
const path = require("path");

function ensureDirectory(dir) {
    if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, { recursive: true });
    }
}

function read(file) {
    return fs.readFileSync(file, "utf8");
}

function write(file, data) {
    ensureDirectory(path.dirname(file));
    fs.writeFileSync(file, data);
}

function json(file, obj) {
    write(file, JSON.stringify(obj, null, 2));
}

module.exports = {
    ensureDirectory,
    read,
    write,
    json
};