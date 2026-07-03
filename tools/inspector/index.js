const utils = require("./utils");
const config = require("./config");
const tree = require("./scanners/tree");
const writer = require("./writers");

utils.ensureDirectory(config.output);

console.log("Aurora Inspector");
console.log("----------------------------");

const treeResult = tree(config.root);

writer.save(
    "tree.json",
    treeResult
);

console.log("✓ Tree generated");

console.log("Done");