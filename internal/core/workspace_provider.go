package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type WorkspaceProvider struct {
	indexer *Indexer

	mu    sync.RWMutex
	cache map[string]map[string]any
}

func NewWorkspaceProvider(indexer *Indexer) *WorkspaceProvider {
	p := &WorkspaceProvider{
		indexer: indexer,
		cache:   make(map[string]map[string]any),
	}
	indexer.SetProvider(p)
	return p
}

func (p *WorkspaceProvider) Invalidate(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.cache, id)
}

func (p *WorkspaceProvider) getRoot(id string) (string, error) {
	ws, ok := p.indexer.Workspace(id)
	if !ok {
		return "", fmt.Errorf("workspace %q not found", id)
	}
	return ws.Root, nil
}

func (p *WorkspaceProvider) loadJSON(id, key, filename string, factory func() any) (any, error) {
	p.mu.RLock()
	if c, ok := p.cache[id]; ok {
		if v, ok := c[key]; ok {
			p.mu.RUnlock()
			return v, nil
		}
	}
	p.mu.RUnlock()

	root, err := p.getRoot(id)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(root, ".devserver", "context", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	val := factory()
	if err := json.Unmarshal(data, val); err != nil {
		return nil, err
	}

	p.mu.Lock()
	if p.cache[id] == nil {
		p.cache[id] = make(map[string]any)
	}
	p.cache[id][key] = val
	p.mu.Unlock()

	return val, nil
}

func (p *WorkspaceProvider) Load(id string) (*WorkspaceContext, error) {
	v, err := p.loadJSON(id, "workspace", "workspace.json", func() any { return &WorkspaceContext{} })
	if err != nil {
		return nil, err
	}
	return v.(*WorkspaceContext), nil
}

func (p *WorkspaceProvider) Architecture(id string) (*ArchitectureInfo, error) {
	v, err := p.loadJSON(id, "architecture", "architecture.json", func() any { return &ArchitectureInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*ArchitectureInfo), nil
}

func (p *WorkspaceProvider) Dependencies(id string) (*DependenciesInfo, error) {
	v, err := p.loadJSON(id, "dependencies", "dependencies.json", func() any { return &DependenciesInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*DependenciesInfo), nil
}

func (p *WorkspaceProvider) Routes(id string) (*RoutesInfo, error) {
	v, err := p.loadJSON(id, "routes", "routes.json", func() any { return &RoutesInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*RoutesInfo), nil
}

func (p *WorkspaceProvider) Database(id string) (*DatabaseInfo, error) {
	v, err := p.loadJSON(id, "database", "database.json", func() any { return &DatabaseInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*DatabaseInfo), nil
}

func (p *WorkspaceProvider) Environment(id string) (*EnvironmentInfo, error) {
	v, err := p.loadJSON(id, "environment", "environment.json", func() any { return &EnvironmentInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*EnvironmentInfo), nil
}

func (p *WorkspaceProvider) Git(id string) (*GitInfo, error) {
	v, err := p.loadJSON(id, "git", "git.json", func() any { return &GitInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*GitInfo), nil
}

func (p *WorkspaceProvider) Health(id string) (*HealthInfo, error) {
	v, err := p.loadJSON(id, "health", "health.json", func() any { return &HealthInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*HealthInfo), nil
}

func (p *WorkspaceProvider) Knowledge(id string) (*KnowledgeInfo, error) {
	v, err := p.loadJSON(id, "knowledge", "knowledge.json", func() any { return &KnowledgeInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*KnowledgeInfo), nil
}

func (p *WorkspaceProvider) Tasks(id string) (*TaskInfo, error) {
	v, err := p.loadJSON(id, "tasks", "tasks.json", func() any { return &TaskInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*TaskInfo), nil
}

func (p *WorkspaceProvider) Project(id string) (*ProjectSummary, error) {
	v, err := p.loadJSON(id, "project", "project.json", func() any { return &ProjectSummary{} })
	if err != nil {
		return nil, err
	}
	return v.(*ProjectSummary), nil
}

func (p *WorkspaceProvider) Plugins(id string) (*PluginInfo, error) {
	v, err := p.loadJSON(id, "plugins", "plugins.json", func() any { return &PluginInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*PluginInfo), nil
}

func (p *WorkspaceProvider) Infrastructure(id string) (*InfrastructureInfo, error) {
	v, err := p.loadJSON(id, "infrastructure", "infrastructure.json", func() any { return &InfrastructureInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*InfrastructureInfo), nil
}

func (p *WorkspaceProvider) Services(id string) (*[]ServiceInfo, error) {
	v, err := p.loadJSON(id, "services", "services.json", func() any { return &[]ServiceInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*[]ServiceInfo), nil
}

func (p *WorkspaceProvider) Deployments(id string) (*DeploymentInfo, error) {
	v, err := p.loadJSON(id, "deployments", "deployments.json", func() any { return &DeploymentInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*DeploymentInfo), nil
}

func (p *WorkspaceProvider) Domains(id string) (*[]DomainInfo, error) {
	v, err := p.loadJSON(id, "domains", "domains.json", func() any { return &[]DomainInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*[]DomainInfo), nil
}

func (p *WorkspaceProvider) Logs(id string) (*[]LogEntry, error) {
	v, err := p.loadJSON(id, "logs", "logs.json", func() any { return &[]LogEntry{} })
	if err != nil {
		return nil, err
	}
	return v.(*[]LogEntry), nil
}

func (p *WorkspaceProvider) AI(id string) (*AIContextInfo, error) {
	v, err := p.loadJSON(id, "ai", "ai.json", func() any { return &AIContextInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*AIContextInfo), nil
}

func (p *WorkspaceProvider) MCP(id string) (*MCPInfo, error) {
	v, err := p.loadJSON(id, "mcp", "mcp.json", func() any { return &MCPInfo{} })
	if err != nil {
		return nil, err
	}
	return v.(*MCPInfo), nil
}

func (p *WorkspaceProvider) Index(id string) (*IndexManifest, error) {
	v, err := p.loadJSON(id, "index", "index.json", func() any { return &IndexManifest{} })
	if err != nil {
		return nil, err
	}
	return v.(*IndexManifest), nil
}

// Special case for Cache since it's located in .devserver/cache/
func (p *WorkspaceProvider) CacheInfo(id string) (*CacheManifest, error) {
	p.mu.RLock()
	if c, ok := p.cache[id]; ok {
		if v, ok := c["cacheInfo"]; ok {
			p.mu.RUnlock()
			return v.(*CacheManifest), nil
		}
	}
	p.mu.RUnlock()

	root, err := p.getRoot(id)
	if err != nil {
		return nil, err
	}

	path := filepath.Join(root, ".devserver", "cache", "index.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	val := &CacheManifest{}
	if err := json.Unmarshal(data, val); err != nil {
		return nil, err
	}

	p.mu.Lock()
	if p.cache[id] == nil {
		p.cache[id] = make(map[string]any)
	}
	p.cache[id]["cacheInfo"] = val
	p.mu.Unlock()

	return val, nil
}

// Derived aggregations for specific endpoints

func (p *WorkspaceProvider) Overview(id string) (map[string]any, error) {
	// Need Workspace, Project, Health, Git, Plugins
	ws, err := p.Load(id)
	if err != nil { return nil, err }
	proj, err := p.Project(id)
	if err != nil { return nil, err }
	health, err := p.Health(id)
	if err != nil { return nil, err }
	gitInfo, err := p.Git(id)
	if err != nil { return nil, err }
	plugins, err := p.Plugins(id)
	if err != nil { return nil, err }

	return map[string]any{
		"workspace": ws,
		"project":   proj,
		"health":    health,
		"git":       gitInfo,
		"plugins":   plugins,
	}, nil
}

func (p *WorkspaceProvider) Doctor(id string) (map[string]any, error) {
	health, err := p.Health(id)
	if err != nil { return nil, err }
	proj, err := p.Project(id)
	if err != nil { return nil, err }
	gitInfo, err := p.Git(id)
	if err != nil { return nil, err }
	plugins, err := p.Plugins(id)
	if err != nil { return nil, err }
	deps, err := p.Dependencies(id)
	if err != nil { return nil, err }
	routes, err := p.Routes(id)
	if err != nil { return nil, err }
	infra, err := p.Infrastructure(id)
	if err != nil { return nil, err }

	score := 0
	if health != nil { score = health.Score }
	branch := ""
	if gitInfo != nil { branch = gitInfo.Branch }

	return map[string]any{
		"health":  health,
		"project": proj,
		"git":     gitInfo,
		"plugins": plugins,
		"checks": []map[string]string{
			{"label": "Dependencies", "status": statusFor(len(deps.Backend)+len(deps.Frontend)+len(deps.Database) > 0, "warning"), "detail": "Analyzed"},
			{"label": "Routes", "status": statusFor(len(routes.API) > 0, "passed"), "detail": fmt.Sprintf("%d routes indexed", len(routes.API))},
			{"label": "Git", "status": statusFor(branch != "", "warning"), "detail": branch},
			{"label": "Health", "status": statusFor(score >= 40, "warning"), "detail": fmt.Sprintf("Score %d", score)},
			{"label": "Infrastructure", "status": statusFor(len(infra.Tools) > 0, "warning"), "detail": fmt.Sprintf("%d tools checked", len(infra.Tools))},
		},
		"recommendations": []string{
			"Review dependency surface and plugins.",
			"Refresh route metadata after new saves.",
			"Keep the generated context in sync with the latest changes.",
		},
	}, nil
}

type WorkspaceFilePreview struct {
	Path       string    `json:"path"`
	Kind       string    `json:"kind"`
	Language   string    `json:"language"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modifiedAt"`
	Content    string    `json:"content"`
	Truncated  bool      `json:"truncated"`
}

func (p *WorkspaceProvider) FilePreview(id, relPath string) (*WorkspaceFilePreview, error) {
	root, err := p.getRoot(id)
	if err != nil {
		return nil, err
	}

	cleaned := filepath.Clean(strings.TrimSpace(relPath))
	if cleaned == "." || cleaned == string(filepath.Separator) || strings.HasPrefix(cleaned, "..") {
		return nil, fmt.Errorf("invalid file path %q", relPath)
	}

	absPath := filepath.Join(root, cleaned)
	rel, err := filepath.Rel(root, absPath)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(rel, "..") {
		return nil, fmt.Errorf("invalid file path %q", relPath)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}

	const maxPreviewBytes = 64 * 1024
	preview := data
	truncated := false
	if len(preview) > maxPreviewBytes {
		preview = preview[:maxPreviewBytes]
		truncated = true
	}

	return &WorkspaceFilePreview{
		Path:       filepath.ToSlash(cleaned),
		Kind:       fileKind(cleaned),
		Language:   languageForFile(cleaned),
		Size:       info.Size(),
		ModifiedAt: info.ModTime().UTC(),
		Content:    strings.ReplaceAll(string(preview), "\r\n", "\n"),
		Truncated:  truncated,
	}, nil
}

func (p *WorkspaceProvider) Files(id string, query string) (map[string]any, error) {
	cacheInfo, err := p.CacheInfo(id)
	if err != nil {
		return nil, err
	}
	
	files := cacheInfo.Files
	if query != "" {
		filtered := make([]WorkspaceFileInfo, 0)
		queryLower := strings.ToLower(query)
		for _, f := range files {
			if strings.Contains(strings.ToLower(f.Path), queryLower) {
				filtered = append(filtered, f)
			}
		}
		files = filtered
	}

	return map[string]any{"files": files, "total": len(files)}, nil
}

func fileKind(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return "Go source"
	case ".ts", ".tsx":
		return "TypeScript"
	case ".js", ".jsx":
		return "JavaScript"
	case ".json":
		return "JSON"
	case ".yaml", ".yml":
		return "YAML"
	case ".md":
		return "Markdown"
	case ".sql":
		return "SQL"
	case ".php":
		return "PHP"
	case ".py":
		return "Python"
	case "":
		return "Folder"
	default:
		return strings.TrimPrefix(strings.ToUpper(filepath.Ext(path)), ".")
	}
}

func languageForFile(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return "go"
	case ".ts", ".tsx":
		return "ts"
	case ".js", ".jsx":
		return "js"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".md":
		return "markdown"
	case ".sql":
		return "sql"
	case ".php":
		return "php"
	case ".py":
		return "python"
	default:
		return "text"
	}
}

func (p *WorkspaceProvider) Settings(id string) (map[string]any, error) {
	ws, err := p.Load(id)
	if err != nil { return nil, err }
	idx, err := p.Index(id)
	if err != nil { return nil, err }
	c, err := p.CacheInfo(id)
	if err != nil { return nil, err }

	return map[string]any{
		"workspace": ws,
		"index":     idx,
		"cache":     c,
	}, nil
}
