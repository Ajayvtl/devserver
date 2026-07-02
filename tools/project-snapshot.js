#!/usr/bin/env node

const fs = require("fs");
const path = require("path");

const ROOT = process.cwd();
const OUTPUT = path.join(ROOT, ".devserver", "snapshot");

const IGNORE = new Set([
  ".git",
  ".next",
  "node_modules",
  "vendor",
  "dist",
  "build",
  ".tmp",
  ".turbo",
  ".cache",
  ".idea",
  ".vscode",
  ".DS_Store",
  "coverage"
]);

const IMPORTANT_FILES = [
  "go.mod",
  "go.sum",
  "package.json",
  "tsconfig.json",
  "next.config.js",
  "next.config.ts",
  "tailwind.config.ts",
  "tailwind.config.js",
  ".gitignore",
  "README.md"
];

fs.mkdirSync(OUTPUT, { recursive: true });

function shouldIgnore(name) {
  return IGNORE.has(name);
}

function walk(dir, depth = 0) {
  let result = [];

  const entries = fs
    .readdirSync(dir, { withFileTypes: true })
    .sort((a, b) => a.name.localeCompare(b.name));

  for (const entry of entries) {
    if (shouldIgnore(entry.name)) continue;

    const full = path.join(dir, entry.name);
    const rel = path.relative(ROOT, full);

    if (entry.isDirectory()) {
      result.push(`${" ".repeat(depth * 2)}📁 ${entry.name}`);
      result.push(...walk(full, depth + 1));
    } else {
      result.push(`${" ".repeat(depth * 2)}📄 ${entry.name}`);
    }
  }

  return result;
}

fs.writeFileSync(
  path.join(OUTPUT, "tree.txt"),
  walk(ROOT).join("\n"),
  "utf8"
);

let inventory = [];

function collect(dir) {
  const entries = fs.readdirSync(dir, { withFileTypes: true });

  for (const e of entries) {
    if (shouldIgnore(e.name)) continue;

    const full = path.join(dir, e.name);

    if (e.isDirectory()) {
      collect(full);
    } else {
      const stat = fs.statSync(full);

      inventory.push({
        file: path.relative(ROOT, full),
        size: stat.size,
        modified: stat.mtime.toISOString()
      });
    }
  }
}

collect(ROOT);

fs.writeFileSync(
  path.join(OUTPUT, "files.json"),
  JSON.stringify(inventory, null, 2)
);

const configs = {};

for (const f of IMPORTANT_FILES) {
  const p = path.join(ROOT, f);

  if (fs.existsSync(p)) {
    configs[f] = fs.readFileSync(p, "utf8");
  }
}

const webPackage = path.join(ROOT, "apps", "web", "package.json");

if (fs.existsSync(webPackage)) {
  configs["apps/web/package.json"] = fs.readFileSync(webPackage, "utf8");
}

fs.writeFileSync(
  path.join(OUTPUT, "configs.json"),
  JSON.stringify(configs, null, 2)
);

const extensions = {};

for (const file of inventory) {
  const ext = path.extname(file.file).toLowerCase() || "none";
  extensions[ext] = (extensions[ext] || 0) + 1;
}

fs.writeFileSync(
  path.join(OUTPUT, "extensions.json"),
  JSON.stringify(extensions, null, 2)
);

const summary = {
  generated: new Date().toISOString(),
  root: ROOT,
  totalFiles: inventory.length,
  extensions
};

fs.writeFileSync(
  path.join(OUTPUT, "summary.json"),
  JSON.stringify(summary, null, 2)
);

console.log("");
console.log("✅ Project snapshot generated");
console.log(OUTPUT);
console.log("");