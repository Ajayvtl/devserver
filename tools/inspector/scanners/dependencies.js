const fs = require("fs");
const path = require("path");

function exists(file) {
    return fs.existsSync(file);
}

function readJson(file) {
    try {
        return JSON.parse(fs.readFileSync(file, "utf8"));
    } catch {
        return null;
    }
}

function readText(file) {
    try {
        return fs.readFileSync(file, "utf8");
    } catch {
        return "";
    }
}

function scan(root) {

    const repositoryRoot = path.resolve(root, "..");

    const dependencies = [];

    const goMod = path.join(repositoryRoot, "go.mod");

    if (exists(goMod)) {

        const content = readText(goMod);

        const matches = [
            ...content.matchAll(/^\s*([A-Za-z0-9./_-]+)\s+v([^\s]+)$/gm)
        ];

        for (const match of matches) {

            dependencies.push({
                ecosystem: "go",
                name: match[1],
                version: match[2]
            });

        }

    }

    const packageJson = path.join(repositoryRoot, "apps", "package.json");

    if (exists(packageJson)) {

        const pkg = readJson(packageJson);

        if (pkg) {

            const all = {
                ...(pkg.dependencies || {}),
                ...(pkg.devDependencies || {})
            };

            for (const [name, version] of Object.entries(all)) {

                dependencies.push({
                    ecosystem: "npm",
                    name,
                    version
                });

            }

        }

    }

    dependencies.sort((a, b) => {

        if (a.ecosystem !== b.ecosystem)
            return a.ecosystem.localeCompare(b.ecosystem);

        return a.name.localeCompare(b.name);

    });

    return dependencies;

}

module.exports = scan;