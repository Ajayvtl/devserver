const path = require("path");

module.exports = {
    root: process.cwd(),

    output: path.join(process.cwd(), "inspection"),

    ignore: [
        ".git",
        "node_modules",
        ".next",
        "dist",
        "build",
        "coverage",
        "inspection",
        "vendor",
        "tmp"
    ],

    extensions: [
        ".js",
        ".jsx",
        ".ts",
        ".tsx",
        ".go",
        ".sql",
        ".json",
        ".md",
        ".yaml",
        ".yml"
    ]
};