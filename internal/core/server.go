package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Ajayvtl/devserver/internal/capabilities"
	"github.com/Ajayvtl/devserver/internal/commands"
	rt "github.com/Ajayvtl/devserver/internal/runtime"
	"github.com/Ajayvtl/devserver/internal/state"
	"github.com/Ajayvtl/devserver/internal/tasks"
	"github.com/rs/zerolog"
)

type APIServer struct {
	log      zerolog.Logger
	store    *state.StoreDB
	indexer  *Indexer
	provider *WorkspaceProvider
	engine   *tasks.Engine
	stream   *tasks.EventStream
	commands *commands.Engine
	caps     *capabilities.Registry
	addr     string
	ctx      context.Context
	cancel   context.CancelFunc
	server   *http.Server
	status   rt.Status
}

func NewAPIServer(log zerolog.Logger, store *state.StoreDB, indexer *Indexer, provider *WorkspaceProvider, engine *tasks.Engine, stream *tasks.EventStream, cmdBus *commands.Engine, caps *capabilities.Registry, addr string) *APIServer {
	return &APIServer{
		log:      log,
		store:    store,
		indexer:  indexer,
		provider: provider,
		engine:   engine,
		stream:   stream,
		commands: cmdBus,
		caps:     caps,
		addr:     addr,
		status:   rt.StatusStopped,
	}
}

func (s *APIServer) Name() string { return "core.APIServer" }

func (s *APIServer) Initialize(ctx context.Context) error {
	s.status = rt.StatusStarting
	return nil
}

func (s *APIServer) Start(ctx context.Context) error {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.status = rt.StatusRunning
	go s.runInternal(s.ctx)
	return nil
}

func (s *APIServer) Stop(ctx context.Context) error {
	if s.cancel != nil {
		s.cancel()
	}
	if s.server != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.server.Shutdown(shutdownCtx)
	}
	s.status = rt.StatusStopped
	return nil
}

func (s *APIServer) Status() rt.Status { return s.status }

func (s *APIServer) Health() rt.Health { return rt.HealthHealthy }

func (s *APIServer) runInternal(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/bootstrap/decision", s.handleBootstrapDecision)
	mux.HandleFunc("/api/setup", s.handleSetup)
	mux.HandleFunc("/api/setup/complete", s.handleSetupComplete)
	mux.HandleFunc("/api/auth/login", s.handleLogin)
	mux.HandleFunc("/api/knowledge", s.handleKnowledge)
	mux.HandleFunc("/api/dashboard", s.handleDashboard)
	mux.HandleFunc("/api/projects", s.handleProjects)
	mux.HandleFunc("/api/projects/", s.handleProjectBySlug)
	mux.HandleFunc("/api/settings", s.handleSettings)
	mux.HandleFunc("/api/devcenter", s.handleDevCenter)
	mux.HandleFunc("/api/devcenter/", s.handleDevCenter)
	mux.HandleFunc("/api/doctor/", s.handleDoctor)
	mux.HandleFunc("/api/workspaces/", s.handleWorkspaces)
	mux.HandleFunc("/api/tasks", s.handleTasks)
	mux.HandleFunc("/api/tasks/", s.handleTaskByID)
	mux.HandleFunc("/api/commands", s.handleCommands)
	mux.HandleFunc("/api/commands/", s.handleCommandByID)
	mux.HandleFunc("/api/capabilities", s.handleCapabilities)
	mux.HandleFunc("/ws/events", s.stream.HandleWS)
	mux.HandleFunc("/ws/tasks", s.stream.HandleWS)

	s.server = &http.Server{
		Addr:    s.addr,
		Handler: s.withCORS(mux),
	}

	errCh := make(chan error, 1)
	go func() {
		s.log.Info().Str("addr", s.addr).Msg("api server listening")
		errCh <- s.server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.status = rt.StatusError
			return err
		}
		return nil
	}
}

func (s *APIServer) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *APIServer) handleBootstrapDecision(w http.ResponseWriter, r *http.Request) {
	data, err := s.store.BootstrapDecision(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, data)
}

func (s *APIServer) handleSetup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := s.store.SetupWizardData(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *APIServer) handleSetupComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req state.SetupRequestData
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	taskID := s.submitTask("setup", "Bootstrap DevServer", "Applying configuration and provisioning services", func() error {
		if err := s.store.CompleteSetup(context.Background(), req); err != nil {
			return err
		}
		return nil
	})
	if taskID == "" {
		writeJSONError(w, http.StatusInternalServerError, fmt.Errorf("failed to queue setup task"))
		return
	}

	writeJSON(w, map[string]any{"taskId": taskID, "status": "started"})
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		data, err := s.store.LoginData(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case http.MethodPost:
		var req state.LoginRequestData
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}

		res, err := s.store.Auth(r.Context(), req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}

		writeJSON(w, res)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *APIServer) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data, err := s.store.DashboardData(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, data)
}

func (s *APIServer) handleProjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := s.store.ListProjects(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, items)
	case http.MethodPost:
		var form state.ProjectFormDataStore
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}
		saved, err := s.store.SaveProject(r.Context(), form)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		go s.submitTask("project", "Save project", "Persisting project metadata and settings", func() error {
			return nil
		})
		writeJSON(w, saved)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *APIServer) handleProjectBySlug(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	if path == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if path == "form" {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		form, err := s.store.ProjectFormData(r.Context(), "")
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, form)
		return
	}
	parts := strings.Split(path, "/")
	slug := parts[0]

	if len(parts) == 2 && parts[1] == "form" {
		form, err := s.store.ProjectFormData(r.Context(), slug)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, form)
		return
	}

	switch r.Method {
	case http.MethodGet:
		data, err := s.store.ProjectDetail(r.Context(), slug)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, err)
			return
		}
		writeJSON(w, data)
	case http.MethodPatch, http.MethodPut:
		var form state.ProjectFormDataStore
		if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}
		form.Slug = slug
		saved, err := s.store.SaveProject(r.Context(), form)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, saved)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *APIServer) handleSettings(w http.ResponseWriter, r *http.Request) {
	data, err := s.store.Settings(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, data)
}

func (s *APIServer) handleDevCenter(w http.ResponseWriter, r *http.Request) {
	if s.provider != nil {
		if r.URL.Path == "/api/devcenter" {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			ctxData, err := s.provider.Load("devserver")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, devcenterSectionListFromContext(*ctxData))
			return
		}
		key := strings.TrimPrefix(r.URL.Path, "/api/devcenter/")
		key = strings.TrimSuffix(key, "/")
		ctxData, err := s.provider.Load("devserver")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writeJSON(w, devcenterSectionFromContext(key, *ctxData))
		return
	}

	if r.URL.Path == "/api/devcenter" {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		sections, err := s.store.DevCenterSections(r.Context())
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, sections)
		return
	}
	key := strings.TrimPrefix(r.URL.Path, "/api/devcenter/")
	key = strings.TrimSuffix(key, "/")
	data, err := s.store.DevCenterSection(r.Context(), key)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, data)
}

func (s *APIServer) handleDoctor(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/api/doctor/")
	if slug == "" {
		slug = "atlas-commerce"
	}

	if s.provider != nil {
		report, err := s.provider.Doctor("devserver")
		if err == nil {
			writeJSON(w, report)
			return
		}
	}

	detail, err := s.store.ProjectDetail(r.Context(), slug)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, map[string]any{
		"project": detail.Project,
		"checks": []map[string]string{
			{"label": "Repository", "status": "passed", "detail": "Connected and fetchable."},
			{"label": "Environment", "status": "warning", "detail": "One secret is missing from the snapshot."},
			{"label": "Deployment target", "status": "passed", "detail": "Healthy and ready to accept rollout."},
		},
		"recommendations": []string{
			"Review environment variables.",
			"Validate SSL renewal coverage.",
			"Inspect queued deployment tasks.",
		},
	})
}

func (s *APIServer) handleWorkspaces(w http.ResponseWriter, r *http.Request) {
	if s.indexer == nil || s.provider == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/api/workspaces/")
	if path == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	parts := strings.Split(path, "/")
	id := parts[0]
	if len(parts) < 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	section := parts[1]
	if section == "index" && r.Method == http.MethodPost {
		if err := s.indexer.Refresh(r.Context(), id); err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, map[string]string{"status": "queued"})
		return
	}

	switch section {
	case "context":
		data, err := s.provider.Load(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "overview":
		data, err := s.provider.Overview(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "filecontent":
		filePath := r.URL.Query().Get("path")
		if filePath == "" {
			writeJSONError(w, http.StatusBadRequest, fmt.Errorf("path query parameter is required"))
			return
		}
		content, err := s.provider.ReadFile(id, filePath)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, map[string]string{"content": content})
	case "files":
		query := r.URL.Query().Get("q")
		data, err := s.provider.Files(id, query)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "repository", "git":
		data, err := s.provider.Git(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "environment":
		data, err := s.provider.Environment(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "infrastructure":
		data, err := s.provider.Infrastructure(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "services":
		data, err := s.provider.Services(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "deployments":
		data, err := s.provider.Deployments(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "database":
		data, err := s.provider.Database(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "domains":
		data, err := s.provider.Domains(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "logs":
		data, err := s.provider.Logs(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "ai":
		data, err := s.provider.AI(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "knowledge":
		data, err := s.provider.Knowledge(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "doctor":
		data, err := s.provider.Doctor(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "settings":
		data, err := s.provider.Settings(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "mcp":
		data, err := s.provider.MCP(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "routes":
		data, err := s.provider.Routes(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "dependencies":
		data, err := s.provider.Dependencies(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	case "health":
		data, err := s.provider.GetHealth(id)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, data)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func devcenterSectionListFromContext(ctxData WorkspaceContext) []map[string]string {
	return []map[string]string{
		{"key": "architecture", "title": "Architecture", "subtitle": ctxData.Architecture.Framework + " architecture and boundaries."},
		{"key": "dependencies", "title": "Dependencies", "subtitle": "Frontend, backend, database, and tool dependencies."},
		{"key": "api", "title": "API", "subtitle": "Generated endpoint and route surface."},
		{"key": "database", "title": "Database", "subtitle": "Generated schema and migration context."},
		{"key": "tasks", "title": "Tasks", "subtitle": "Suggested maintenance and generation tasks."},
		{"key": "knowledge", "title": "Knowledge", "subtitle": "Generated learning and documentation context."},
		{"key": "project-doctor", "title": "Project Doctor", "subtitle": "Health and risk checks from the indexer."},
	}
}

func devcenterSectionFromContext(key string, ctxData WorkspaceContext) state.DevCenterSectionDataStore {
	section := state.DevCenterSectionDataStore{
		Key:      key,
		Subtitle: "Generated from workspace metadata.",
		Metadata: map[string]string{
			"workspace": ctxData.Workspace.Name,
			"framework": ctxData.Workspace.Framework,
		},
	}

	switch key {
	case "architecture":
		section.Title = "Architecture"
		section.Description = "Workspace architecture and boundaries generated from the current codebase."
		section.Bullets = []string{
			"Framework: " + blankAs(ctxData.Architecture.Framework, "unknown"),
			"Runtime: " + blankAs(ctxData.Architecture.Runtime, "unknown"),
			"Languages: " + strings.Join(ctxData.Architecture.Languages, ", "),
			"Entry points: " + strings.Join(ctxData.Architecture.EntryPoints, ", "),
		}
		section.Examples = ctxData.Architecture.Notes
		section.Knowledge = []state.KnowledgeArticleData{}
		section.CodeSample = "Generated architecture metadata is stored in .devserver/context/architecture.json"
	case "dependencies":
		section.Title = "Dependencies"
		section.Description = "Detected dependencies grouped by surface."
		section.Bullets = []string{
			"Frontend: " + strings.Join(ctxData.Dependencies.Frontend, ", "),
			"Backend: " + strings.Join(ctxData.Dependencies.Backend, ", "),
			"Database: " + strings.Join(ctxData.Dependencies.Database, ", "),
			"Tools: " + strings.Join(ctxData.Dependencies.Tools, ", "),
		}
		section.Examples = []string{"plugins.json", "dependencies.json"}
		section.CodeSample = "Dependencies are derived from package manifests and workspace files."
	case "api":
		section.Title = "API"
		section.Description = "Generated API and route surface."
		section.Bullets = []string{
			"Endpoints: " + strings.Join(ctxData.Routes.API, ", "),
			"Next routes: " + strings.Join(ctxData.Routes.Next, ", "),
			"Go routes: " + strings.Join(ctxData.Routes.Go, ", "),
		}
		section.Examples = []string{"routes.json", "api.json"}
		section.CodeSample = "API metadata is read-only and regenerated by the workspace indexer."
	case "database":
		section.Title = "Database"
		section.Description = "Detected database footprint and migrations."
		section.Bullets = []string{
			"Kind: " + blankAs(ctxData.Database.Kind, "unknown"),
			"Files: " + strings.Join(ctxData.Database.Files, ", "),
			"Migrations: " + strings.Join(ctxData.Database.Migrations, ", "),
		}
		section.Examples = []string{"database.json", "migrations/"}
		section.CodeSample = "Database metadata is inferred from schema and migration files."
	case "tasks":
		section.Title = "Tasks"
		section.Description = "Suggested follow-up tasks based on current workspace metadata."
		section.Bullets = ctxData.Tasks.Suggested
		section.Examples = []string{"tasks.json"}
		section.CodeSample = "Task suggestions are generated from dependency, route, and database signals."
	case "knowledge":
		section.Title = "Knowledge"
		section.Description = "Generated learning context and documentation sources."
		section.Bullets = []string{
			"Docs: " + strings.Join(ctxData.Knowledge.Files, ", "),
			"Topics: " + strings.Join(ctxData.Knowledge.Topics, ", "),
		}
		section.Examples = []string{"README.md", "architecture.md", ".devserver/knowledge/"}
		section.CodeSample = "Knowledge is generated from docs and workspace metadata."
	case "project-doctor":
		section.Title = "Project Doctor"
		section.Description = "Health and risk checks built on indexer output."
		section.Bullets = []string{
			"Health score: " + fmt.Sprintf("%d", ctxData.Health.Score),
			"Git branch: " + blankAs(ctxData.Git.Branch, "unknown"),
			"Detected plugins: " + strings.Join(ctxData.Plugins.Detected, ", "),
		}
		section.Examples = []string{"health.json", "git.json", "plugins.json"}
		section.CodeSample = "Doctor consumes indexer metadata only."
	default:
		section.Title = "Architecture"
		section.Description = "Workspace architecture and boundaries generated from the current codebase."
	}

	return section
}

func blankAs(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func (s *APIServer) handleKnowledge(w http.ResponseWriter, r *http.Request) {
	scope := r.URL.Query().Get("scope")
	items, err := s.store.KnowledgeArticles(r.Context(), scope)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, items)
}

func (s *APIServer) handleTasks(w http.ResponseWriter, r *http.Request) {
	items, err := s.store.Tasks(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, items)
}

func (s *APIServer) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	task, err := s.store.GetTask(r.Context(), id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, task)
}

func (s *APIServer) handleCommands(w http.ResponseWriter, r *http.Request) {
	if s.commands == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, s.commands.List())
	case http.MethodPost:
		var req commands.Command
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}
		record, err := s.commands.Submit(r.Context(), &req)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, record)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *APIServer) handleCommandByID(w http.ResponseWriter, r *http.Request) {
	if s.commands == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/commands/")
	if id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	record, err := s.commands.Get(id)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, record)
}

func (s *APIServer) handleCapabilities(w http.ResponseWriter, r *http.Request) {
	if s.caps == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if capName := strings.TrimSpace(r.URL.Query().Get("capability")); capName != "" {
		bindings, err := capabilities.NewResolver(s.caps).Resolve(capabilities.Capability(capName))
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, bindings)
		return
	}

	out := make(map[string][]capabilities.Metadata)
	for _, provider := range s.caps.Providers() {
		for _, cap := range provider.Capabilities() {
			meta := provider.Metadata(cap)
			meta.Provider = provider.Name()
			meta.Capability = cap
			out[string(cap)] = append(out[string(cap)], meta)
		}
	}
	writeJSON(w, out)
}

func (s *APIServer) submitTask(scope, name, detail string, work func() error) string {
	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())
	_, _ = s.store.CreateTaskWithID(context.Background(), taskID, scope, name, detail)

	task := &tasks.Task{
		ID:          taskID,
		Name:        name,
		WorkspaceID: scope,
		Type:        scope,
		Priority:    tasks.TaskPriorityNormal,
		Metadata:    map[string]string{"scope": scope, "detail": detail},
	}

	record, err := s.engine.Submit(task)
	if err != nil {
		_ = s.store.UpdateTask(context.Background(), taskID, 0, "Failed", err.Error())
		s.log.Error().Err(err).Str("task", name).Msg("submit task failed")
		return ""
	}

	// Legacy closure execution inline to avoid breaking legacy endpoints before Runners are made
	if work != nil {
		go work()
	}

	return record.ID
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func writeJSONError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	writeJSON(w, map[string]string{"error": err.Error()})
}

func statusFor(condition bool, fallback string) string {
	if condition {
		return "passed"
	}
	return fallback
}
