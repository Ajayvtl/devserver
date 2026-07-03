const fs = require("fs");
const path = require("path");
const config = require("./config");
const utils = require("./utils");

function save(name, data) {
    utils.json(
        path.join(config.output, name),
        data
    );
}

module.exports = {
    save
};