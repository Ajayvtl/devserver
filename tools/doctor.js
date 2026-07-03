const fs = require("fs");

const checks = [
  ".git",
  "go.mod",
  "apps/package.json",
  "cmd/devserver/main.go"
];

for (const file of checks) {
  console.log(
    fs.existsSync(file)
      ? `✅ ${file}`
      : `❌ ${file}`
  );
}