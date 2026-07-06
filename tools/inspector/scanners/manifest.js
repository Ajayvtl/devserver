const fs = require("fs");
const path = require("path");
const config = require("../config");

const MANIFESTS = new Set([
    "go.mod",
    "go.sum",
    "package.json",
    "package-lock.json",
    "pnpm-lock.yaml",
    "yarn.lock",
    "Dockerfile",
    "docker-compose.yml",
    "docker-compose.yaml",
    "compose.yml",
    "compose.yaml",
    ".env",
    ".env.local",
    ".env.development",
    ".env.production",
    "tsconfig.json",
    "jsconfig.json",
    "next.config.js",
    "next.config.mjs",
    "next.config.ts",
    "vite.config.js",
    "vite.config.ts",
    "tailwind.config.js",
    "tailwind.config.ts",
    "eslint.config.js",
    ".eslintrc",
    ".eslintrc.json",
    ".prettierrc",
    ".prettierrc.json"
]);

function scan(root) {
    const repositoryRoot = path.resolve(root, "..");
 
    const manifests = [];

    walk(repositoryRoot);

    manifests.sort((a, b) => a.path.localeCompare(b.path));

    return manifests;

    function walk(dir) {

        for (const entry of fs.readdirSync(dir)) {

            if (config.ignore.includes(entry))
                continue;

            const full = path.join(dir, entry);
            const stat = fs.statSync(full);

            if (stat.isDirectory()) {
                walk(full);
                continue;
            }

            if (!MANIFESTS.has(entry))
                continue;

            manifests.push({
                name: entry,
                path: path.relative(config.root, full).replace(/\\/g, "/"),
                size: stat.size,
                modified: stat.mtime.toISOString()
            });
        }
    }
}

module.exports = scan;