package core

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Ajayvtl/devserver/internal/events"
	"github.com/rs/zerolog"
)

type WorkspaceSpec struct {
	ID   string
	Root string
}

type FileStamp struct {
	ModTime time.Time
	Size    int64
	Hash    string
}

type WorkspaceSummary struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Root           string    `json:"root"`
	Kind           string    `json:"kind"`
	Framework      string    `json:"framework"`
	Languages      []string  `json:"languages"`
	Runtime        string    `json:"runtime"`
	PackageManager string    `json:"packageManager"`
	GeneratedAt    time.Time `json:"generatedAt"`
}

type WorkspaceFileInfo struct {
	Path       string    `json:"path"`
	Kind       string    `json:"kind"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type ProjectSummary struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Root           string   `json:"root"`
	Framework      string   `json:"framework"`
	Languages      []string `json:"languages"`
	Runtime        string   `json:"runtime"`
	PackageManager string   `json:"packageManager"`
	Repository     string   `json:"repository"`
}

type ArchitectureInfo struct {
	Kind           string   `json:"kind"`
	Framework      string   `json:"framework"`
	Languages      []string `json:"languages"`
	Runtime        string   `json:"runtime"`
	PackageManager string   `json:"packageManager"`
	EntryPoints    []string `json:"entryPoints"`
	Notes          []string `json:"notes"`
}

type DependenciesInfo struct {
	Frontend []string `json:"frontend"`
	Backend  []string `json:"backend"`
	Database []string `json:"database"`
	Tools    []string `json:"tools"`
}

type RouteNode struct {
	Path     string      `json:"path"`
	Children []RouteNode `json:"children,omitempty"`
}

type RoutesInfo struct {
	Next []string    `json:"next"`
	Go   []string    `json:"go"`
	API  []string    `json:"api"`
	Tree []RouteNode `json:"tree"`
}

type DatabaseInfo struct {
	Kind        string   `json:"kind"`
	Files       []string `json:"files"`
	Migrations  []string `json:"migrations"`
	Environment []string `json:"environment"`
}

type EnvironmentInfo struct {
	Files   []string          `json:"files"`
	Keys    []string          `json:"keys"`
	Secrets []string          `json:"secrets"`
	Preview map[string]string `json:"preview"`
}

type GitInfo struct {
	Branch        string   `json:"branch"`
	Commit        string   `json:"commit"`
	Clean         bool     `json:"clean"`
	Ahead         int      `json:"ahead"`
	Behind        int      `json:"behind"`
	RecentCommits []string `json:"recentCommits"`
	ChangedFiles  []string `json:"changedFiles"`
}

type HealthInfo struct {
	Score int    `json:"score"`
	Build string `json:"build"`
	Tests string `json:"tests"`
	Lint  string `json:"lint"`
}

type PluginInfo struct {
	Detected     []string `json:"detected"`
	Capabilities []string `json:"capabilities"`
}

type KnowledgeInfo struct {
	Files  []string `json:"files"`
	Topics []string `json:"topics"`
}

type TaskInfo struct {
	Suggested []string `json:"suggested"`
}

type InfraToolInfo struct {
	Name       string `json:"name"`
	Installed  bool   `json:"installed"`
	Version    string `json:"version"`
	Healthy    bool   `json:"healthy"`
	Configured bool   `json:"configured"`
}

type InfrastructureInfo struct {
	Tools      []InfraToolInfo `json:"tools"`
	OS         string          `json:"os"`
	Arch       string          `json:"arch"`
	Hostname   string          `json:"hostname"`
}

type ServiceInfo struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Port       string `json:"port"`
	PID        string `json:"pid"`
}

type DeploymentEntry struct {
	ID          string `json:"id"`
	Branch      string `json:"branch"`
	Commit      string `json:"commit"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
	Environment string `json:"environment"`
}

type DeploymentInfo struct {
	Entries []DeploymentEntry `json:"entries"`
	Current string            `json:"current"`
}

type DomainInfo struct {
	Host   string `json:"host"`
	SSL    string `json:"ssl"`
	Target string `json:"target"`
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Source    string `json:"source"`
}

type AIContextInfo struct {
	Workspace    WorkspaceSummary `json:"workspace"`
	Git          GitInfo          `json:"git"`
	Architecture ArchitectureInfo `json:"architecture"`
	Tasks        TaskInfo         `json:"tasks"`
	Dependencies DependenciesInfo `json:"dependencies"`
	Routes       RoutesInfo       `json:"routes"`
	Database     DatabaseInfo     `json:"database"`
	GeneratedAt  time.Time        `json:"generatedAt"`
}

type MCPProvider struct {
	Name       string `json:"name"`
	Enabled    bool   `json:"enabled"`
	Healthy    bool   `json:"healthy"`
	Endpoint   string `json:"endpoint"`
	Protocol   string `json:"protocol"`
}

type MCPInfo struct {
	Providers []MCPProvider `json:"providers"`
}

type IndexManifest struct {
	Version     int               `json:"version"`
	GeneratedAt time.Time         `json:"generatedAt"`
	Workspace   string            `json:"workspace"`
	Files       map[string]string `json:"files"`
}

type CacheManifest struct {
	UpdatedAt time.Time          `json:"updatedAt"`
	Fingerprint string           `json:"fingerprint"`
	Changed    []string          `json:"changed"`
	Files      []WorkspaceFileInfo `json:"files"`
}

type WorkspaceContext struct {
	Workspace      WorkspaceSummary   `json:"workspace"`
	Project        ProjectSummary     `json:"project"`
	Architecture   ArchitectureInfo   `json:"architecture"`
	Dependencies   DependenciesInfo   `json:"dependencies"`
	Routes         RoutesInfo         `json:"routes"`
	Database       DatabaseInfo       `json:"database"`
	Environment    EnvironmentInfo    `json:"environment"`
	Git            GitInfo            `json:"git"`
	Tasks          TaskInfo           `json:"tasks"`
	Health         HealthInfo         `json:"health"`
	Plugins        PluginInfo         `json:"plugins"`
	Knowledge      KnowledgeInfo      `json:"knowledge"`
	Infrastructure InfrastructureInfo `json:"infrastructure"`
	Services       []ServiceInfo      `json:"services"`
	Deployments    DeploymentInfo     `json:"deployments"`
	Domains        []DomainInfo       `json:"domains"`
	Logs           []LogEntry         `json:"logs"`
	AI             AIContextInfo      `json:"ai"`
	MCP            MCPInfo            `json:"mcp"`
	Index          IndexManifest      `json:"index"`
	Files          []WorkspaceFileInfo `json:"files"`
	Cache          CacheManifest      `json:"cache"`
}

type workspaceState struct {
	spec     WorkspaceSpec
	snapshot map[string]FileStamp
	context  WorkspaceContext
}

type Indexer struct {
	log        zerolog.Logger
	mu         sync.RWMutex
	workspaces map[string]*workspaceState
	bus        events.Bus
	provider   *WorkspaceProvider
}

func (i *Indexer) SetProvider(p *WorkspaceProvider) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.provider = p
}

func (i *Indexer) Workspace(id string) (*WorkspaceSpec, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	ws, ok := i.workspaces[id]
	if !ok {
		return nil, false
	}
	return &ws.spec, true
}

func NewIndexer(log zerolog.Logger, specs []WorkspaceSpec, bus events.Bus) *Indexer {
	workspaces := make(map[string]*workspaceState, len(specs))
	for _, spec := range specs {
		if spec.ID == "" || spec.Root == "" {
			continue
		}
		workspaces[spec.ID] = &workspaceState{spec: spec, snapshot: map[string]FileStamp{}}
	}
	return &Indexer{log: log, workspaces: workspaces, bus: bus}
}

// AddWorkspace registers a workspace for indexing at runtime.
func (i *Indexer) AddWorkspace(id, root string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.workspaces[id] = &workspaceState{
		spec:     WorkspaceSpec{ID: id, Root: root},
		snapshot: map[string]FileStamp{},
	}
}

// RemoveWorkspace unregisters a workspace from indexing.
func (i *Indexer) RemoveWorkspace(id string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.workspaces, id)
}

// WorkspaceIDs returns all currently registered workspace IDs.
func (i *Indexer) WorkspaceIDs() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	ids := make([]string, 0, len(i.workspaces))
	for id := range i.workspaces {
		ids = append(ids, id)
	}
	return ids
}

func (i *Indexer) Run(ctx context.Context) error {
	if err := i.RefreshAll(ctx); err != nil {
		return err
	}
	
	ch := i.bus.Subscribe(events.WorkspaceChanged)
	for {
		select {
		case <-ctx.Done():
			return nil
		case evt := <-ch:
			payload, ok := evt.Payload.(events.WorkspaceChangedEvent)
			if ok {
				i.log.Debug().Str("workspace", payload.WorkspaceID).Msg("received change event, refreshing")
				_ = i.Refresh(ctx, payload.WorkspaceID)
			}
		}
	}
}

func (i *Indexer) RefreshAll(ctx context.Context) error {
	i.mu.RLock()
	ids := make([]string, 0, len(i.workspaces))
	for id := range i.workspaces {
		ids = append(ids, id)
	}
	i.mu.RUnlock()
	sort.Strings(ids)
	for _, id := range ids {
		if err := i.Refresh(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (i *Indexer) Refresh(ctx context.Context, id string) error {
	i.mu.RLock()
	ws, ok := i.workspaces[id]
	i.mu.RUnlock()
	if !ok {
		return fmt.Errorf("workspace %q not found", id)
	}

	current, err := collectRelevantFiles(ws.spec.Root)
	if err != nil {
		return err
	}

	dirty := diffFiles(ws.snapshot, current)
	if len(dirty) == 0 && len(ws.snapshot) > 0 {
		return nil
	}

	next := composeContext(ws, current, dirtySections(dirty))
	if err := writeContext(ws.spec.Root, next); err != nil {
		return err
	}

	if err := writeCache(ws.spec.Root, next); err != nil {
		return err
	}
	i.mu.Lock()
	ws.snapshot = current
	ws.context = next
	p := i.provider
	i.mu.Unlock()
	
	if p != nil {
		p.Invalidate(id)
	}

	i.bus.Publish(events.WorkspaceIndexed, events.WorkspaceIndexedEvent{
		WorkspaceID:  id,
		ChangedFiles: len(dirty),
	})

	for _, file := range dirty {
    i.log.Debug().
        Str("workspace", id).
        Str("file", file).
        Msg("changed file")
	}
	i.log.Info().Str("workspace", id).Int("changed_files", len(dirty)).Msg("workspace indexed")
	return nil
}

func (i *Indexer) Context(id string) (WorkspaceContext, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	ws, ok := i.workspaces[id]
	if !ok {
		return WorkspaceContext{}, false
	}
	return ws.context, true
}

func collectRelevantFiles(root string) (map[string]FileStamp, error) {
	out := map[string]FileStamp{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			if shouldSkipDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if !isRelevantFile(rel) {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		sum, err := fileHash(path)
		if err != nil {
			return err
		}
		out[rel] = FileStamp{ModTime: info.ModTime().UTC(), Size: info.Size(), Hash: sum}
		return nil
	})
	return out, err
}

func fileHash(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha1.Sum(data)
	return fmt.Sprintf("%x", sum[:]), nil
}

func diffFiles(prev, current map[string]FileStamp) []string {
	changed := make([]string, 0)
	for path, now := range current {
		before, ok := prev[path]
		if !ok || before.ModTime != now.ModTime || before.Size != now.Size || before.Hash != now.Hash {
			changed = append(changed, path)
		}
	}
	for path := range prev {
		if _, ok := current[path]; !ok {
			changed = append(changed, path)
		}
	}
	return changed
}

func isRelevantFile(rel string) bool {
    if strings.Contains(rel, ".devserver/context/") ||
        strings.Contains(rel, ".devserver/cache/") ||
        strings.Contains(rel, ".devserver/snapshot/") {
        return false
    }

    base := filepath.Base(rel)

    if strings.HasPrefix(base, ".env") ||
        base == "go.mod" ||
        base == "go.sum" ||
        base == "package.json" ||
        base == "composer.json" ||
        base == "pyproject.toml" ||
        strings.HasPrefix(base, "requirements") ||
        base == "pnpm-lock.yaml" ||
        base == "yarn.lock" ||
        base == "package-lock.json" ||
        base == "docker-compose.yml" ||
        base == "docker-compose.yaml" ||
        base == "Dockerfile" ||
        strings.HasSuffix(base, ".sql") {
        return true
    }

    if strings.HasSuffix(base, ".go") ||
        strings.HasSuffix(base, ".js") ||
        strings.HasSuffix(base, ".jsx") ||
        strings.HasSuffix(base, ".ts") ||
        strings.HasSuffix(base, ".tsx") ||
        strings.HasSuffix(base, ".py") ||
        strings.HasSuffix(base, ".php") ||
        strings.HasSuffix(base, ".md") ||
        strings.HasSuffix(base, ".yaml") ||
        strings.HasSuffix(base, ".yml") ||
        strings.HasSuffix(base, ".json") {
        return true
    }

    if strings.Contains(rel, ".devserver/knowledge/") {
        return true
    }

    return false
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git",
     "node_modules",
     ".next",
     "build",
     "dist",
     "vendor",
     ".tmp",
	 ".turbo",
	 ".vercel",
     "coverage",
     ".devserver":
    return true
	default:
		return false
	}
}

func composeContext(ws *workspaceState, files map[string]FileStamp, dirty []string) WorkspaceContext {
	base := ws.context
	workspace := detectWorkspace(ws.spec, files)
	project := detectProject(workspace)
	inventory := summarizeFiles(files)

	if len(dirty) == 0 {
		dirty = []string{"workspace", "project", "architecture", "dependencies", "routes", "database", "environment", "git", "tasks", "health", "plugins", "knowledge", "infrastructure", "index"}
	}

	if needs(dirty, "workspace") {
		base.Workspace = workspace
	}
	if needs(dirty, "project") {
		base.Project = project
	}
	if needs(dirty, "architecture") {
		base.Architecture = detectArchitecture(workspace, files)
	}
	if needs(dirty, "dependencies") {
		base.Dependencies = detectDependencies(ws.spec.Root, files)
	}
	if needs(dirty, "routes") {
		base.Routes = detectRoutes(ws.spec.Root, files)
	}
	if needs(dirty, "database") {
		base.Database = detectDatabase(ws.spec.Root, files)
	}
	if needs(dirty, "environment") {
		base.Environment = detectEnvironment(ws.spec.Root, files)
	}
	if needs(dirty, "git") {
		base.Git = detectGit(ws.spec.Root)
	}
	if needs(dirty, "health") {
		base.Health = detectHealth(files, base.Git)
	}
	if needs(dirty, "plugins") {
		base.Plugins = detectPlugins(base.Workspace, base.Dependencies, base.Database, base.Git)
	}
	if needs(dirty, "knowledge") {
		base.Knowledge = detectKnowledge(ws.spec.Root, files)
	}
	if needs(dirty, "tasks") {
		base.Tasks = detectTasks(base.Workspace, base.Dependencies, base.Routes, base.Database)
	}
	if needs(dirty, "infrastructure") {
		base.Infrastructure = detectInfrastructure()
	}
	base.Services = detectServices(base.Infrastructure)
	base.Deployments = DeploymentInfo{Current: base.Git.Commit, Entries: buildDeploymentHistory(base.Git)}
	base.Domains = detectDomains(ws.spec.Root)
	base.Logs = readRecentLogs(ws.spec.Root)
	base.AI = composeAIContext(base)
	base.MCP = detectMCP(ws.spec.Root)
	base.Files = inventory
	base.Cache = detectCache(files, dirty, base.Routes, base.Dependencies, base.Git)

	base.Index = IndexManifest{
		Version:     1,
		GeneratedAt: time.Now().UTC(),
		Workspace:   ws.spec.ID,
		Files: map[string]string{
			"workspace":     "workspace.json",
			"project":       "project.json",
			"architecture":  "architecture.json",
			"dependencies":  "dependencies.json",
			"routes":        "routes.json",
			"api":           "api.json",
			"database":      "database.json",
			"environment":   "environment.json",
			"git":           "git.json",
			"tasks":         "tasks.json",
			"health":        "health.json",
			"plugins":       "plugins.json",
			"knowledge":     "knowledge.json",
			"index":         "index.json",
			"cache":         "../cache/index.json",
			"cache-scan":    "../cache/scan.json",
			"cache-git":     "../cache/git.json",
			"cache-routes":   "../cache/routes.json",
			"cache-deps":     "../cache/dependencies.json",
		},
	}
	return base
}

func needs(dirty []string, key string) bool {
	for _, item := range dirty {
		if item == key {
			return true
		}
	}
	return false
}

func dirtySections(files []string) []string {
	sections := map[string]struct{}{}
	add := func(name string) {
		sections[name] = struct{}{}
	}
	for _, rel := range files {
		base := filepath.Base(rel)
		switch {
		case strings.Contains(rel, ".devserver/knowledge/") || strings.HasSuffix(base, "README.md") || strings.HasSuffix(base, "architecture.md"):
			add("knowledge")
		case strings.HasPrefix(base, ".env"):
			add("environment")
			add("project")
			add("index")
		case base == "go.mod" || base == "go.sum" || base == "package.json" || base == "composer.json" || base == "pyproject.toml" || strings.HasPrefix(base, "requirements") || base == "pnpm-lock.yaml" || base == "yarn.lock" || base == "package-lock.json":
			add("workspace")
			add("project")
			add("architecture")
			add("dependencies")
			add("plugins")
			add("health")
			add("index")
		case strings.HasSuffix(base, ".go") || strings.HasSuffix(base, ".ts") || strings.HasSuffix(base, ".tsx") || strings.HasSuffix(base, ".js") || strings.HasSuffix(base, ".jsx") || strings.HasSuffix(base, ".py") || strings.HasSuffix(base, ".php"):
			add("workspace")
			add("project")
			add("architecture")
			add("routes")
			add("tasks")
			add("knowledge")
			add("index")
		case strings.HasSuffix(base, ".sql") || strings.Contains(rel, "migrations/"):
			add("database")
			add("dependencies")
			add("plugins")
			add("project")
			add("index")
		case base == "Dockerfile" || base == "docker-compose.yml" || base == "docker-compose.yaml":
			add("workspace")
			add("dependencies")
			add("plugins")
			add("health")
			add("index")
		default:
			add("project")
			add("index")
		}
	}
	out := make([]string, 0, len(sections))
	for key := range sections {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func detectWorkspace(spec WorkspaceSpec, files map[string]FileStamp) WorkspaceSummary {
	root := strings.TrimSpace(spec.Root)
	name := filepath.Base(root)
	framework := detectFramework(root, files)
	languages := detectLanguages(files)
	return WorkspaceSummary{
		ID:             spec.ID,
		Name:           prettifyName(name),
		Root:           root,
		Kind:           workspaceKind(framework, languages),
		Framework:      framework,
		Languages:      languages,
		Runtime:        detectRuntime(framework, languages),
		PackageManager: detectPackageManager(files),
		GeneratedAt:    time.Now().UTC(),
	}
}

func detectProject(ws WorkspaceSummary) ProjectSummary {
	return ProjectSummary{
		ID:             ws.ID,
		Name:           ws.Name,
		Root:           ws.Root,
		Framework:      ws.Framework,
		Languages:      ws.Languages,
		Runtime:        ws.Runtime,
		PackageManager: ws.PackageManager,
		Repository:     detectRepository(ws.Root),
	}
}

func detectArchitecture(ws WorkspaceSummary, files map[string]FileStamp) ArchitectureInfo {
	return ArchitectureInfo{
		Kind:           ws.Kind,
		Framework:      ws.Framework,
		Languages:      ws.Languages,
		Runtime:        ws.Runtime,
		PackageManager: ws.PackageManager,
		EntryPoints:    entryPoints(detectRoutes(ws.Root, files)),
		Notes:          []string{"Read-only metadata generated by the workspace indexer."},
	}
}

func detectDependencies(root string, files map[string]FileStamp) DependenciesInfo {
	frontend := []string{}
	backend := []string{}
	database := []string{}
	tools := []string{}

	for _, rel := range pathsByBase(files, "package.json") {
		if pkg, ok := readPackageJSON(filepath.Join(root, rel)); ok {
			for name := range pkg.Dependencies {
				classifyDependency(name, &frontend, &backend, &database, &tools)
			}
			for name := range pkg.DevDependencies {
				classifyDependency(name, &frontend, &backend, &database, &tools)
			}
		}
	}
	for _, dep := range readGoModDeps(filepath.Join(root, "go.mod")) {
		classifyDependency(dep, &frontend, &backend, &database, &tools)
	}
	if hasFile(files, "composer.json") {
		backend = appendUnique(backend, "php")
	}
	return DependenciesInfo{
		Frontend: normalizeList(frontend),
		Backend:  normalizeList(backend),
		Database: normalizeList(database),
		Tools:    normalizeList(tools),
	}
}

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func readPackageJSON(path string) (packageJSON, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return packageJSON{}, false
	}
	var out packageJSON
	if err := json.Unmarshal(data, &out); err != nil {
		return packageJSON{}, false
	}
	return out, true
}

func readGoModDeps(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	out := []string{"go"}
	for _, line := range strings.Split(string(data), "\n") {
		lower := strings.ToLower(line)
		switch {
		case strings.Contains(lower, "postgres"):
			out = appendUnique(out, "postgres")
		case strings.Contains(lower, "mysql"):
			out = appendUnique(out, "mysql")
		case strings.Contains(lower, "redis"):
			out = appendUnique(out, "redis")
		}
	}
	return out
}

func classifyDependency(name string, frontend, backend, database, tools *[]string) {
	n := strings.ToLower(name)
	switch {
	case n == "next" || n == "react" || n == "react-dom" || n == "tailwindcss" || strings.Contains(n, "vite"):
		*frontend = appendUnique(*frontend, name)
	case strings.Contains(n, "express") || strings.Contains(n, "fastify") || strings.Contains(n, "koa") || strings.Contains(n, "gin") || strings.Contains(n, "fiber") || strings.Contains(n, "echo") || strings.Contains(n, "flask") || strings.Contains(n, "django") || strings.Contains(n, "laravel"):
		*backend = appendUnique(*backend, name)
	case strings.Contains(n, "postgres") || strings.Contains(n, "mysql") || strings.Contains(n, "sqlite") || strings.Contains(n, "prisma") || strings.Contains(n, "gorm") || strings.Contains(n, "typeorm"):
		*database = appendUnique(*database, name)
	default:
		*tools = appendUnique(*tools, name)
	}
}

func detectRoutes(root string, files map[string]FileStamp) RoutesInfo {
	nextRoutes := []string{}
	goRoutes := []string{}
	for rel := range files {
		if strings.HasSuffix(rel, ".go") {
			goRoutes = append(goRoutes, extractGoRoutes(filepath.Join(root, rel))...)
		}
		if route := extractNextRoute(rel); route != "" {
			nextRoutes = append(nextRoutes, route)
		}
	}
	nextRoutes = normalizeList(nextRoutes)
	goRoutes = normalizeList(goRoutes)
	apiRoutes := normalizeList(append(append([]string{}, nextRoutes...), goRoutes...))
	return RoutesInfo{
		Next: nextRoutes,
		Go:   goRoutes,
		API:  apiRoutes,
		Tree: buildRouteTree(apiRoutes),
	}
}

func extractGoRoutes(path string) []string {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	re := regexp.MustCompile(`HandleFunc\(\s*"([^"]+)"|Handle\(\s*"([^"]+)"|HandlePrefix\(\s*"([^"]+)"`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	out := make([]string, 0, len(matches))
	for _, match := range matches {
		for i := 1; i < len(match); i++ {
			if strings.TrimSpace(match[i]) != "" {
				out = append(out, match[i])
				break
			}
		}
	}
	return out
}

func extractNextRoute(rel string) string {
	rel = filepath.ToSlash(rel)
	parts := strings.Split(rel, "/")
	for i := 0; i < len(parts); i++ {
		if parts[i] == "app" {
			tail := parts[i+1:]
			if len(tail) == 0 {
				return ""
			}
			last := tail[len(tail)-1]
			if isPageFile(last) || isRouteFile(last) {
				return normalizeRouteSegments(tail[:len(tail)-1])
			}
		}
		if parts[i] == "pages" {
			tail := parts[i+1:]
			if len(tail) == 0 {
				return "/"
			}
			last := tail[len(tail)-1]
			if isPageFile(last) {
				return normalizeRouteSegments(tail[:len(tail)-1])
			}
		}
	}
	return ""
}

func normalizeRouteSegments(segments []string) string {
	if len(segments) == 0 {
		return "/"
	}
	cleaned := make([]string, 0, len(segments))
	for _, segment := range segments {
		switch segment {
		case "page.tsx", "page.ts", "page.jsx", "page.js", "route.ts", "route.js", "index.tsx", "index.ts", "index.jsx", "index.js":
			continue
		}
		if strings.HasPrefix(segment, "[") && strings.HasSuffix(segment, "]") {
			cleaned = append(cleaned, ":"+strings.Trim(segment, "[]"))
			continue
		}
		cleaned = append(cleaned, segment)
	}
	if len(cleaned) == 0 {
		return "/"
	}
	return "/" + strings.Join(cleaned, "/")
}

func isPageFile(name string) bool {
	switch name {
	case "page.tsx", "page.ts", "page.jsx", "page.js", "index.tsx", "index.ts", "index.jsx", "index.js":
		return true
	default:
		return false
	}
}

func isRouteFile(name string) bool {
	return name == "route.ts" || name == "route.js"
}

func buildRouteTree(routes []string) []RouteNode {
	root := map[string]*RouteNode{}
	for _, route := range routes {
		route = strings.Trim(route, "/")
		if route == "" {
			continue
		}
		parts := strings.Split(route, "/")
		insertRoute(root, parts)
	}
	return flattenRouteTree(root)
}

func insertRoute(tree map[string]*RouteNode, parts []string) {
	if len(parts) == 0 {
		return
	}
	key := parts[0]
	node, ok := tree[key]
	if !ok {
		node = &RouteNode{Path: "/" + key}
		tree[key] = node
	}
	if len(parts) == 1 {
		return
	}
	children := map[string]*RouteNode{}
	for _, child := range node.Children {
		childKey := strings.TrimPrefix(child.Path, "/")
		childCopy := child
		children[childKey] = &childCopy
	}
	insertRoute(children, parts[1:])
	node.Children = flattenRouteTree(children)
}

func flattenRouteTree(tree map[string]*RouteNode) []RouteNode {
	keys := make([]string, 0, len(tree))
	for key := range tree {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	out := make([]RouteNode, 0, len(keys))
	for _, key := range keys {
		if node := tree[key]; node != nil {
			out = append(out, *node)
		}
	}
	return out
}

func detectDatabase(root string, files map[string]FileStamp) DatabaseInfo {
	kind := ""
	switch {
	case contains(files, "postgres") || hasDependency(files, root, "postgres"):
		kind = "postgres"
	case contains(files, "mysql") || hasDependency(files, root, "mysql"):
		kind = "mysql"
	case contains(files, "sqlite") || hasDependency(files, root, "sqlite"):
		kind = "sqlite"
	}
	return DatabaseInfo{
		Kind:        kind,
		Files:       filesWithSuffix(files, ".sql"),
		Migrations:  filesMatching(files, "migrations/"),
		Environment: environmentFiles(files),
	}
}

func detectEnvironment(root string, files map[string]FileStamp) EnvironmentInfo {
	envFiles := []string{}
	keys := []string{}
	secrets := []string{}
	preview := map[string]string{}
	for rel := range files {
		if !strings.HasPrefix(filepath.Base(rel), ".env") {
			continue
		}
		envFiles = append(envFiles, rel)
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			keys = appendUnique(keys, key)
			preview[key] = val
			if strings.Contains(strings.ToUpper(key), "SECRET") || strings.Contains(strings.ToUpper(key), "TOKEN") || strings.Contains(strings.ToUpper(key), "PASSWORD") {
				secrets = appendUnique(secrets, key)
			}
		}
	}
	return EnvironmentInfo{Files: normalizeList(envFiles), Keys: normalizeList(keys), Secrets: normalizeList(secrets), Preview: preview}
}

func detectGit(root string) GitInfo {
	base := GitInfo{Clean: true}
	if _, err := os.Stat(filepath.Join(root, ".git")); err != nil {
		return base
	}

	branch := trimGitOutput(runGit(root, "rev-parse", "--abbrev-ref", "HEAD"))
	commit := trimGitOutput(runGit(root, "rev-parse", "HEAD"))
	status := trimGitOutput(runGit(root, "status", "--porcelain"))
	ahead, behind := gitAheadBehind(root)
	recent := splitGitLines(runGit(root, "log", "--pretty=format:%h %s", "-n", "5"))
	changed := splitGitLines(runGit(root, "status", "--porcelain"))
	if branch == "" {
		branch = detectBranchFromHead(root)
	}
	if commit == "" {
		commit = detectCommitFromHead(root)
	}
	if status == "" {
		status = ""
	}
	return GitInfo{
		Branch:        branch,
		Commit:        commit,
		Clean:         status == "",
		Ahead:         ahead,
		Behind:        behind,
		RecentCommits: recent,
		ChangedFiles:  changed,
	}
}

func readGitRef(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func runGit(root string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func splitGitLines(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	lines := strings.Split(value, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return normalizeList(out)
}

func trimGitOutput(value string) string {
	return strings.TrimSpace(value)
}

func gitAheadBehind(root string) (int, int) {
	out := runGit(root, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	if out == "" {
		return 0, 0
	}
	parts := strings.Fields(out)
	if len(parts) != 2 {
		return 0, 0
	}
	ahead := atoiSafe(parts[0])
	behind := atoiSafe(parts[1])
	return ahead, behind
}

func detectBranchFromHead(root string) string {
	headPath := filepath.Join(root, ".git", "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}
	head := strings.TrimSpace(string(data))
	if strings.HasPrefix(head, "ref: ") {
		return filepath.Base(strings.TrimSpace(strings.TrimPrefix(head, "ref: ")))
	}
	return ""
}

func detectCommitFromHead(root string) string {
	headPath := filepath.Join(root, ".git", "HEAD")
	data, err := os.ReadFile(headPath)
	if err != nil {
		return ""
	}
	head := strings.TrimSpace(string(data))
	if strings.HasPrefix(head, "ref: ") {
		ref := strings.TrimSpace(strings.TrimPrefix(head, "ref: "))
		if refCommit := readGitRef(filepath.Join(root, ".git", ref)); refCommit != "" {
			return refCommit
		}
	}
	return head
}

func detectPlugins(ws WorkspaceSummary, deps DependenciesInfo, db DatabaseInfo, git GitInfo) PluginInfo {
	detected := []string{}
	add := func(v string) {
		if v != "" {
			detected = appendUnique(detected, v)
		}
	}
	add(ws.Framework)
	for _, dep := range deps.Frontend {
		switch strings.ToLower(dep) {
		case "react":
			add("React")
		case "next":
			add("Next.js")
		case "tailwindcss":
			add("Tailwind")
		}
	}
	for _, lang := range ws.Languages {
		add(lang)
	}
	if db.Kind != "" {
		add(strings.Title(db.Kind))
	}
	if git.Branch != "" {
		add("Git")
	}
	if ws.PackageManager != "" {
		add(strings.Title(ws.PackageManager))
	}
	if _, err := os.Stat(filepath.Join(ws.Root, "Dockerfile")); err == nil {
		add("Docker")
	}

	capabilities := []string{}
	if anyMatch(detected, "Next.js", "React", "Laravel", "Django", "FastAPI") {
		capabilities = append(capabilities, "Web")
	}
	if anyMatch(detected, "Go", "Node", "Python", "PHP") {
		capabilities = append(capabilities, "API")
	}
	if db.Kind != "" {
		capabilities = append(capabilities, "Database")
	}
	if anyMatch(detected, "Redis") {
		capabilities = append(capabilities, "Redis")
	}
	if anyMatch(detected, "Docker") {
		capabilities = append(capabilities, "Docker")
	}
	capabilities = append(capabilities, "CLI", "WebSocket")
	return PluginInfo{Detected: normalizeList(detected), Capabilities: normalizeList(capabilities)}
}

func detectHealth(files map[string]FileStamp, git GitInfo) HealthInfo {
	score := 0
	if len(files) > 0 {
		score += 40
	}
	if git.Branch != "" {
		score += 20
	}
	return HealthInfo{Score: score, Build: "unknown", Tests: "unknown", Lint: "unknown"}
}

func detectKnowledge(root string, files map[string]FileStamp) KnowledgeInfo {
	out := []string{}
	topics := []string{}
	for rel := range files {
		if strings.Contains(rel, ".devserver/knowledge/") || strings.HasSuffix(rel, "README.md") || strings.HasSuffix(rel, "architecture.md") {
			out = append(out, rel)
			topics = append(topics, filepath.Base(rel))
		}
	}
	return KnowledgeInfo{Files: normalizeList(out), Topics: normalizeList(topics)}
}

func summarizeFiles(files map[string]FileStamp) []WorkspaceFileInfo {
	out := make([]WorkspaceFileInfo, 0, len(files))
	for rel, stamp := range files {
		out = append(out, WorkspaceFileInfo{
			Path:       rel,
			Kind:       fileKind(rel),
			Size:       stamp.Size,
			ModifiedAt: stamp.ModTime,
		})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Path == out[j].Path {
			return out[i].ModifiedAt.Before(out[j].ModifiedAt)
		}
		return out[i].Path < out[j].Path
	})
	return out
}

func fileKind(rel string) string {
	switch strings.ToLower(filepath.Ext(rel)) {
	case ".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".php", ".json", ".yaml", ".yml", ".toml", ".md":
		return "file"
	}
	if strings.Contains(filepath.Base(rel), ".") {
		return "file"
	}
	return "folder"
}

func detectCache(files map[string]FileStamp, dirty []string, routes RoutesInfo, deps DependenciesInfo, git GitInfo) CacheManifest {
	changed := append([]string{}, dirty...)
	sort.Strings(changed)
	fingerprintParts := make([]string, 0, len(files))
	inventory := summarizeFiles(files)
	for _, file := range inventory {
		fingerprintParts = append(fingerprintParts, fmt.Sprintf("%s:%d:%s", file.Path, file.Size, file.ModifiedAt.UTC().Format(time.RFC3339Nano)))
	}
	fingerprintParts = append(fingerprintParts, strings.Join(routes.API, ","), strings.Join(deps.Frontend, ","), strings.Join(deps.Backend, ","), strings.Join(deps.Database, ","), git.Commit, git.Branch)
	sum := sha1.Sum([]byte(strings.Join(fingerprintParts, "|")))
	return CacheManifest{
		UpdatedAt:   time.Now().UTC(),
		Fingerprint: fmt.Sprintf("%x", sum[:]),
		Changed:     changed,
		Files:       inventory,
	}
}

func detectTasks(ws WorkspaceSummary, deps DependenciesInfo, routes RoutesInfo, db DatabaseInfo) TaskInfo {
	tasks := []string{}
	if ws.Framework != "" || len(routes.API) > 0 {
		tasks = append(tasks, "sync routes")
	}
	if len(deps.Frontend)+len(deps.Backend)+len(deps.Database) > 0 {
		tasks = append(tasks, "refresh dependencies")
	}
	if db.Kind != "" {
		tasks = append(tasks, "validate database")
	}
	return TaskInfo{Suggested: normalizeList(tasks)}
}

func workspaceKind(framework string, langs []string) string {
	if framework != "" {
		return framework
	}
	if anyMatch(langs, "go") {
		return "Go"
	}
	if anyMatch(langs, "python") {
		return "Python"
	}
	return "workspace"
}

func detectFramework(root string, files map[string]FileStamp) string {
	for _, rel := range pathsByBase(files, "package.json") {
		if pkg, ok := readPackageJSON(filepath.Join(root, rel)); ok {
			if hasKey(pkg.Dependencies, "next") || hasKey(pkg.DevDependencies, "next") {
				return "Next.js"
			}
			if hasKey(pkg.Dependencies, "react") || hasKey(pkg.DevDependencies, "react") {
				return "React"
			}
			if hasKey(pkg.Dependencies, "express") || hasKey(pkg.DevDependencies, "express") {
				return "Express"
			}
		}
	}
	if hasFile(files, "go.mod") {
		return "Go"
	}
	if hasFile(files, "composer.json") {
		return "Laravel"
	}
	if hasFile(files, "requirements.txt") || hasFile(files, "pyproject.toml") {
		if hasKeyword(files, root, "fastapi") {
			return "FastAPI"
		}
		if hasKeyword(files, root, "django") {
			return "Django"
		}
		return "Python"
	}
	return ""
}

func detectLanguages(files map[string]FileStamp) []string {
	out := []string{}
	if hasFile(files, "package.json") || hasExt(files, ".ts") || hasExt(files, ".tsx") || hasExt(files, ".js") || hasExt(files, ".jsx") {
		out = append(out, "TypeScript", "JavaScript")
	}
	if hasFile(files, "go.mod") || hasExt(files, ".go") {
		out = append(out, "Go")
	}
	if hasExt(files, ".py") || hasFile(files, "requirements.txt") || hasFile(files, "pyproject.toml") {
		out = append(out, "Python")
	}
	if hasExt(files, ".php") || hasFile(files, "composer.json") {
		out = append(out, "PHP")
	}
	return normalizeList(out)
}

func detectPackageManager(files map[string]FileStamp) string {
	switch {
	case hasAnyBase(files, "pnpm-lock.yaml"):
		return "pnpm"
	case hasAnyBase(files, "yarn.lock"):
		return "yarn"
	case hasAnyBase(files, "package-lock.json"):
		return "npm"
	default:
		return ""
	}
}

func detectRuntime(framework string, langs []string) string {
	switch {
	case strings.Contains(strings.ToLower(framework), "next"), strings.Contains(strings.ToLower(framework), "react"), strings.Contains(strings.ToLower(framework), "express"):
		return "node"
	case anyMatch(langs, "go"):
		return "go"
	case anyMatch(langs, "python"):
		return "python"
	case anyMatch(langs, "php"):
		return "php"
	default:
		return ""
	}
}

func detectRepository(root string) string {
	if _, err := os.Stat(filepath.Join(root, ".git")); err == nil {
		return "git"
	}
	return ""
}

func hasFile(files map[string]FileStamp, name string) bool {
	for rel := range files {
		if filepath.Base(rel) == name {
			return true
		}
	}
	return false
}

func hasAnyBase(files map[string]FileStamp, name string) bool {
	for rel := range files {
		if filepath.Base(rel) == name {
			return true
		}
	}
	return false
}

func hasKey(values map[string]string, key string) bool {
	for k := range values {
		if strings.EqualFold(k, key) {
			return true
		}
	}
	return false
}

func hasDependency(files map[string]FileStamp, root, name string) bool {
	for _, rel := range pathsByBase(files, "package.json") {
		pkg, ok := readPackageJSON(filepath.Join(root, rel))
		if !ok {
			continue
		}
		if hasKey(pkg.Dependencies, name) || hasKey(pkg.DevDependencies, name) {
			return true
		}
	}
	return false
}

func hasExt(files map[string]FileStamp, ext string) bool {
	for rel := range files {
		if strings.HasSuffix(strings.ToLower(rel), ext) {
			return true
		}
	}
	return false
}

func hasKeyword(files map[string]FileStamp, root, keyword string) bool {
	for rel := range files {
		data, err := os.ReadFile(filepath.Join(root, rel))
		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(string(data)), strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}

func contains(files map[string]FileStamp, keyword string) bool {
	for rel := range files {
		if strings.Contains(strings.ToLower(rel), keyword) {
			return true
		}
	}
	return false
}

func filesWithSuffix(files map[string]FileStamp, suffix string) []string {
	out := []string{}
	for rel := range files {
		if strings.HasSuffix(rel, suffix) {
			out = append(out, rel)
		}
	}
	return normalizeList(out)
}

func filesMatching(files map[string]FileStamp, token string) []string {
	out := []string{}
	for rel := range files {
		if strings.Contains(rel, token) {
			out = append(out, rel)
		}
	}
	return normalizeList(out)
}

func environmentFiles(files map[string]FileStamp) []string {
	out := []string{}
	for rel := range files {
		if strings.HasPrefix(filepath.Base(rel), ".env") {
			out = append(out, rel)
		}
	}
	return normalizeList(out)
}

func pathsByBase(files map[string]FileStamp, name string) []string {
	out := []string{}
	for rel := range files {
		if filepath.Base(rel) == name {
			out = append(out, rel)
		}
	}
	return normalizeList(out)
}

func entryPoints(routes RoutesInfo) []string {
	out := append([]string{}, routes.Next...)
	out = append(out, routes.Go...)
	return normalizeList(out)
}

func appendUnique(values []string, next string) []string {
	for _, value := range values {
		if strings.EqualFold(value, next) {
			return values
		}
	}
	return append(values, next)
}

func normalizeList(values []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		key := strings.ToLower(value)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func anyMatch(values []string, needles ...string) bool {
	for _, value := range values {
		for _, needle := range needles {
			if strings.EqualFold(value, needle) {
				return true
			}
		}
	}
	return false
}

func prettifyName(name string) string {
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	if name == "" {
		return "workspace"
	}
	return strings.Title(name)
}

func writeContext(root string, ctx WorkspaceContext) error {
	base := filepath.Join(root, ".devserver", "context")
	if err := os.MkdirAll(base, 0o755); err != nil {
		return err
	}
	files := map[string]any{
		"workspace.json":      ctx.Workspace,
		"project.json":        ctx.Project,
		"architecture.json":   ctx.Architecture,
		"dependencies.json":   ctx.Dependencies,
		"routes.json":         ctx.Routes,
		"api.json":            map[string]any{"routes": ctx.Routes.API, "generatedAt": ctx.Index.GeneratedAt},
		"database.json":       ctx.Database,
		"environment.json":    ctx.Environment,
		"git.json":            ctx.Git,
		"tasks.json":          ctx.Tasks,
		"health.json":         ctx.Health,
		"plugins.json":        ctx.Plugins,
		"knowledge.json":      ctx.Knowledge,
		"infrastructure.json": ctx.Infrastructure,
		"services.json":       ctx.Services,
		"deployments.json":    ctx.Deployments,
		"domains.json":        ctx.Domains,
		"ai.json":             ctx.AI,
		"mcp.json":            ctx.MCP,
		"index.json":          ctx.Index,
	}
	for name, payload := range files {
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(base, name), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func writeCache(root string, ctx WorkspaceContext) error {
	base := filepath.Join(root, ".devserver", "cache")
	if err := os.MkdirAll(base, 0o755); err != nil {
		return err
	}
	files := map[string]any{
		"index.json":          ctx.Cache,
		"scan.json":           ctx.Cache,
		"git.json":            ctx.Git,
		"routes.json":         ctx.Routes,
		"dependencies.json":   ctx.Dependencies,
		"infrastructure.json": ctx.Infrastructure,
		"services.json":       ctx.Services,
		"deployments.json":    ctx.Deployments,
		"domains.json":        ctx.Domains,
		"ai.json":             ctx.AI,
		"mcp.json":            ctx.MCP,
	}
	for name, payload := range files {
		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(base, name), data, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func atoiSafe(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// detectInfrastructure probes the system for installed tools.
func detectInfrastructure() InfrastructureInfo {
	hostname, _ := os.Hostname()
	tools := []InfraToolInfo{
		probeToolVersion("Node.js", "node", "--version"),
		probeToolVersion("Go", "go", "version"),
		probeToolVersion("Python", "python3", "--version"),
		probeToolVersion("PHP", "php", "--version"),
		probeToolVersion("Redis", "redis-cli", "--version"),
		probeToolVersion("PostgreSQL", "psql", "--version"),
		probeToolVersion("MySQL", "mysql", "--version"),
		probeToolVersion("Docker", "docker", "--version"),
		probeToolVersion("Nginx", "nginx", "-v"),
		probeToolVersion("PM2", "pm2", "--version"),
		probeToolVersion("Git", "git", "--version"),
	}
	return InfrastructureInfo{
		Tools:    tools,
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Hostname: hostname,
	}
}

func probeToolVersion(name, binary, flag string) InfraToolInfo {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary, flag)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return InfraToolInfo{Name: name, Installed: false, Version: "", Healthy: false, Configured: false}
	}
	version := strings.TrimSpace(string(out))
	// Extract just the first meaningful line
	if idx := strings.Index(version, "\n"); idx > 0 {
		version = version[:idx]
	}
	return InfraToolInfo{
		Name:       name,
		Installed:  true,
		Version:    version,
		Healthy:    true,
		Configured: true,
	}
}

func detectServices(infra InfrastructureInfo) []ServiceInfo {
	services := []ServiceInfo{}
	for _, tool := range infra.Tools {
		if !tool.Installed {
			continue
		}
		switch tool.Name {
		case "Redis":
			services = append(services, ServiceInfo{Name: "Redis", Status: probeServiceStatus("redis-cli", "ping"), Port: "6379"})
		case "PostgreSQL":
			services = append(services, ServiceInfo{Name: "PostgreSQL", Status: probeServiceStatus("pg_isready", ""), Port: "5432"})
		case "MySQL":
			services = append(services, ServiceInfo{Name: "MySQL", Status: "installed", Port: "3306"})
		case "Nginx":
			services = append(services, ServiceInfo{Name: "Nginx", Status: probeServiceStatus("nginx", "-t"), Port: "80"})
		case "Docker":
			services = append(services, ServiceInfo{Name: "Docker", Status: probeServiceStatus("docker", "info"), Port: ""})
		}
	}
	return services
}

func probeServiceStatus(binary, arg string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if arg == "" {
		cmd = exec.CommandContext(ctx, binary)
	} else {
		cmd = exec.CommandContext(ctx, binary, arg)
	}
	_, err := cmd.CombinedOutput()
	if err != nil {
		return "stopped"
	}
	return "running"
}

func buildDeploymentHistory(git GitInfo) []DeploymentEntry {
	entries := []DeploymentEntry{}
	for i, commit := range git.RecentCommits {
		parts := strings.SplitN(commit, " ", 2)
		hash := commit
		message := ""
		if len(parts) == 2 {
			hash = parts[0]
			message = parts[1]
		}
		entries = append(entries, DeploymentEntry{
			ID:          fmt.Sprintf("deploy-%d", i+1),
			Branch:      git.Branch,
			Commit:      hash,
			Status:      "succeeded",
			Timestamp:   message,
			Environment: "production",
		})
	}
	return entries
}

func detectDomains(root string) []DomainInfo {
	domains := []DomainInfo{}
	// Check nginx config fragments
	nginxConf := filepath.Join(root, "nginx.conf")
	if data, err := os.ReadFile(nginxConf); err == nil {
		re := regexp.MustCompile(`server_name\s+([^;]+);`)
		matches := re.FindAllStringSubmatch(string(data), -1)
		for _, match := range matches {
			if len(match) > 1 {
				for _, host := range strings.Fields(match[1]) {
					if host != "_" && host != "" {
						domains = append(domains, DomainInfo{Host: host, SSL: "Pending", Target: "localhost"})
					}
				}
			}
		}
	}
	return domains
}

func readRecentLogs(root string) []LogEntry {
	entries := []LogEntry{}
	logDir := filepath.Join(root, ".devserver", "logs")
	if _, err := os.Stat(logDir); err != nil {
		return entries
	}
	_ = filepath.WalkDir(logDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !strings.HasSuffix(d.Name(), ".log") && !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		lines := strings.Split(string(data), "\n")
		start := 0
		if len(lines) > 50 {
			start = len(lines) - 50
		}
		for _, line := range lines[start:] {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			entries = append(entries, LogEntry{
				Timestamp: time.Now().UTC().Format(time.RFC3339),
				Level:     "info",
				Message:   line,
				Source:    d.Name(),
			})
		}
		return nil
	})
	if len(entries) > 100 {
		entries = entries[len(entries)-100:]
	}
	return entries
}

func composeAIContext(ctx WorkspaceContext) AIContextInfo {
	return AIContextInfo{
		Workspace:    ctx.Workspace,
		Git:          ctx.Git,
		Architecture: ctx.Architecture,
		Tasks:        ctx.Tasks,
		Dependencies: ctx.Dependencies,
		Routes:       ctx.Routes,
		Database:     ctx.Database,
		GeneratedAt:  time.Now().UTC(),
	}
}

func detectMCP(root string) MCPInfo {
	providers := []MCPProvider{
		{Name: "DevServer Context", Enabled: true, Healthy: true, Endpoint: "/api/workspaces/{id}/ai", Protocol: "http"},
	}
	// Check for .devserver/mcp.json configuration
	mcpPath := filepath.Join(root, ".devserver", "mcp.json")
	if data, err := os.ReadFile(mcpPath); err == nil {
		var cfg struct {
			Providers []MCPProvider `json:"providers"`
		}
		if json.Unmarshal(data, &cfg) == nil && len(cfg.Providers) > 0 {
			providers = append(providers, cfg.Providers...)
		}
	}
	return MCPInfo{Providers: providers}
}
