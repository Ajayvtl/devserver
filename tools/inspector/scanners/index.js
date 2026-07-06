module.exports = [
    {
        name: "Tree",
        output: "tree.json",
        scanner: require("./tree")
    },
    {
        name: "Filesystem",
        output: "filesystem.json",
        scanner: require("./filesystem")
    },
    {
        name: "Manifest",
        output: "manifest.json",
        scanner: require("./manifest")
    },
    {
        name: "dependencies",
        output: "dependencies.json",
        scanner: require("./dependencies")
    }
];