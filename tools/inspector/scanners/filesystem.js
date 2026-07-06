const fs = require("fs");
const path = require("path");
const config = require("../config");

const LANGUAGE_MAP = {
    ".go": "go",
    ".js": "javascript",
    ".jsx": "javascript",
    ".ts": "typescript",
    ".tsx": "typescript",
    ".json": "json",
    ".sql": "sql",
    ".md": "markdown",
    ".yaml": "yaml",
    ".yml": "yaml"
};

function language(ext) {
    return LANGUAGE_MAP[ext.toLowerCase()] || "unknown";
}

function scan(root) {
    const repositoryRoot = path.resolve(root, "..");
    const results = [];

    walk(repositoryRoot);

    return results;

    function walk(dir) {

        const entries = fs.readdirSync(dir);

        for (const entry of entries) {

            if (config.ignore.includes(entry))
                continue;

            const full = path.join(dir, entry);

            const stat = fs.statSync(full);

            const relative = path.relative(config.root, full);

            if (stat.isDirectory()) {

                results.push({
                    type: "directory",
                    name: entry,
                    path: relative.replaceAll("\\", "/"),
                    extension: "",
                    language: "",
                    size: 0,
                    modified: stat.mtime.toISOString()
                });

                walk(full);

                continue;
            }

            const ext = path.extname(entry);

            results.push({
                type: "file",
                name: entry,
                path: relative.replaceAll("\\", "/"),
                extension: ext,
                language: language(ext),
                size: stat.size,
                modified: stat.mtime.toISOString()
            });

        }

    }

}

module.exports = scan;