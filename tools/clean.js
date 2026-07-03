const fs = require("fs");

const dirs = [
  "apps/.next",
  ".tmp",
  ".devserver/cache",
  ".devserver/snapshot"
];

for (const dir of dirs) {
  fs.rmSync(dir, {
    recursive: true,
    force: true
  });
  console.log("Removed", dir);
}