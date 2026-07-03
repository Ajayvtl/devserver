const fs = require("fs");
const path = require("path");
const config = require("../config");

function walk(dir) {

    let items = [];

    for (const file of fs.readdirSync(dir)) {

        if (config.ignore.includes(file))
            continue;

        const full = path.join(dir, file);

        const stat = fs.statSync(full);

        if (stat.isDirectory()) {

            items.push({
                type: "directory",
                name: file,
                path: path.relative(config.root, full),
                children: walk(full)
            });

        } else {

            items.push({
                type: "file",
                name: file,
                path: path.relative(config.root, full),
                size: stat.size
            });

        }

    }

    return items;
}

module.exports = walk;