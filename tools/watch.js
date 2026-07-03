#!/usr/bin/env node

const chokidar = require("chokidar");
const path = require("path");

const watcher = chokidar.watch(".", {
  ignored: [
    "**/.git/**",
    "**/node_modules/**",
    "**/.next/**",
    "**/.tmp/**",
    "**/.devserver/cache/**",
    "**/.devserver/context/**",
    "**/.devserver/snapshot/**"
  ],
  ignoreInitial: true,
  persistent: true
});

watcher.on("all", (event, file) => {
  console.log(
    `[${new Date().toLocaleTimeString()}] ${event} ${path.normalize(file)}`
  );
});

console.log("Watching...");