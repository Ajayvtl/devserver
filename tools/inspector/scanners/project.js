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

    const result = {
        name: path.basename(repositoryRoot),
        root: repositoryRoot,
        type: "unknown",
        languages: [],
        frameworks: [],
        packageManagers: [],
        module: null,
        package: null,
        docker: false
    };

    const goMod = path.join(repositoryRoot, "go.mod");
    if (exists(goMod)) {
        result.type = "go";
        result.languages.push("go");

        const content = readText(goMod);

        const match = content.match(/^module\s+(.+)$/m);

        if (match)
            result.module = match[1].trim();
    }

    const packageJson = path.join(repositoryRoot, "apps", "package.json");

    if (exists(packageJson)) {

        const pkg = readJson(packageJson);

        if (pkg) {

            result.package = pkg.name || null;

            if (result.type === "go")
                result.type = "mixed";
            else
                result.type = "node";

            if (!result.languages.includes("javascript"))
                result.languages.push("javascript");

            result.packageManagers.push("npm");

            const deps = {
                ...(pkg.dependencies || {}),
                ...(pkg.devDependencies || {})
            };

            if (deps.next)
                result.frameworks.push("Next.js");

            if (deps.react)
                result.frameworks.push("React");

            if (deps.express)
                result.frameworks.push("Express");

            if (deps.typescript && !result.languages.includes("typescript"))
                result.languages.push("typescript");
        }
    }

    if (
        exists(path.join(repositoryRoot, "Dockerfile")) ||
        exists(path.join(repositoryRoot, "docker-compose.yml")) ||
        exists(path.join(repositoryRoot, "compose.yml"))
    ) {
        result.docker = true;
    }

    result.languages.sort();
    result.frameworks.sort();
    result.packageManagers.sort();

    return result;
}

module.exports = scan;