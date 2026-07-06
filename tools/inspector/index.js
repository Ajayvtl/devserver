const utils = require("./utils");
const config = require("./config");

const tree = require("./scanners/tree");
const filesystem = require("./scanners/filesystem");
const manifest = require("./scanners/manifest");
const dependencies = require("./scanners/dependencies");
const writer = require("./writers");

const project = require("./scanners/project");

utils.ensureDirectory(config.output);

console.log("DevServer Inspector");
console.log("----------------------------");

const treeResult = tree(config.root);

writer.save("tree.json", treeResult);

console.log("✓ Tree generated");

const filesystemResult = filesystem(config.root);

writer.save("filesystem.json", filesystemResult);

console.log("✓ Filesystem generated");

const manifestResult = manifest(config.root);

writer.save("manifest.json", manifestResult);

console.log("✓ Manifest generated");

const projectResult = project(config.root);

writer.save(
    "project.json",
    projectResult
);

console.log("✓ Project generated");
const depenDenciesResult = dependencies(config.root);

writer.save(
    "dependencies.json",
    depenDenciesResult
);

console.log("✓ Dependencies generated");
console.log("Done");