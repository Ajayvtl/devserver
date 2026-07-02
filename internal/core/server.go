package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Ajayvtl/devserver/internal/state"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/rs/zerolog"
)

type APIServer struct {
	log   zerolog.Logger
	store *state.StoreDB
	indexer *Indexer
	hub   *TaskHub
	addr  string
}

func NewAPIServer(log zerolog.Logger, store *state.StoreDB, indexer *Indexer, addr string) *APIServer {
	return &APIServer{
		log:   log,
		store: store,
		indexer: indexer,
		hub:   NewTaskHub(log),
		addr:  addr,
	}
}

func (s *APIServer) Run(ctx context.Context) error {
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
	mux.HandleFunc("/ws/tasks", s.hub.handleWS)

	if s.indexer != nil {
		go func() {
			_ = s.indexer.Run(ctx)
		}()
	}

	server := &http.Server{
		Addr:    s.addr,
		Handler: s.withCORS(mux),
	}

	errCh := make(chan error, 1)
	go func() {
		s.log.Info().Str("addr", s.addr).Msg("api server listening")
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	task, err := s.store.CreateTask(r.Context(), "setup", "Bootstrap DevServer", "Applying configuration and provisioning services")
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	go s.simulateTask(task.ID, "setup", "Bootstrap DevServer", "Applying configuration and provisioning services", func() error {
		return s.store.CompleteSetup(context.Background(), req)
	})

	writeJSON(w, map[string]any{"taskId": task.ID, "status": "started"})
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
		task, _ := s.store.CreateTask(r.Context(), "project", "Save project", "Persisting project metadata and settings")
		go s.simulateTask(task.ID, "project", "Save project", "Persisting project metadata and settings", nil)
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
	if s.indexer != nil {
		if r.URL.Path == "/api/devcenter" {
			if r.Method != http.MethodGet {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			ctxData, ok := s.indexer.Context("devserver")
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			writeJSON(w, devcenterSectionListFromContext(ctxData))
			return
		}
		key := strings.TrimPrefix(r.URL.Path, "/api/devcenter/")
		key = strings.TrimSuffix(key, "/")
		ctxData, ok := s.indexer.Context("devserver")
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writeJSON(w, devcenterSectionFromContext(key, ctxData))
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
	report := map[string]any{}
	if ctxData, ok := s.indexer.Context("devserver"); ok {
		report = map[string]any{
			"project": ctxData.Project,
			"checks": []map[string]string{
				{"label": "Dependencies", "status": statusFor(len(ctxData.Dependencies.Backend)+len(ctxData.Dependencies.Frontend)+len(ctxData.Dependencies.Database) > 0, "warning"), "detail": strings.Join(ctxData.Dependencies.Backend, ", ")},
				{"label": "Routes", "status": statusFor(len(ctxData.Routes.API) > 0, "passed"), "detail": fmt.Sprintf("%d routes indexed", len(ctxData.Routes.API))},
				{"label": "Git", "status": statusFor(ctxData.Git.Branch != "", "warning"), "detail": ctxData.Git.Branch},
				{"label": "Health", "status": statusFor(ctxData.Health.Score >= 40, "warning"), "detail": fmt.Sprintf("Score %d", ctxData.Health.Score)},
			},
			"recommendations": []string{
				"Review dependency surface and plugins.",
				"Refresh route metadata after new saves.",
				"Keep the generated context in sync with the latest changes.",
			},
		}
	} else {
		detail, err := s.store.ProjectDetail(r.Context(), slug)
		if err != nil {
			writeJSONError(w, http.StatusNotFound, err)
			return
		}
		report = map[string]any{
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
		}
	}
	writeJSON(w, report)
}

func (s *APIServer) handleWorkspaces(w http.ResponseWriter, r *http.Request) {
	if s.indexer == nil {
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
	ctxData, ok := s.indexer.Context(id)
	if !ok {
		writeJSONError(w, http.StatusNotFound, fmt.Errorf("workspace %q not indexed", id))
		return
	}
	switch section {
	case "context":
		writeJSON(w, ctxData)
	case "overview":
		writeJSON(w, map[string]any{
			"workspace": ctxData.Workspace,
			"project":   ctxData.Project,
			"health":    ctxData.Health,
			"git":       ctxData.Git,
			"plugins":   ctxData.Plugins,
		})
	case "files":
		query := r.URL.Query().Get("q")
		files := ctxData.Files
		if query != "" {
			filtered := make([]WorkspaceFileInfo, 0)
			for _, f := range files {
				if strings.Contains(strings.ToLower(f.Path), strings.ToLower(query)) {
					filtered = append(filtered, f)
				}
			}
			files = filtered
		}
		writeJSON(w, map[string]any{"files": files, "total": len(files)})
	case "repository", "git":
		writeJSON(w, ctxData.Git)
	case "environment":
		writeJSON(w, ctxData.Environment)
	case "infrastructure":
		writeJSON(w, ctxData.Infrastructure)
	case "services":
		writeJSON(w, ctxData.Services)
	case "deployments":
		writeJSON(w, ctxData.Deployments)
	case "database":
		writeJSON(w, ctxData.Database)
	case "domains":
		writeJSON(w, ctxData.Domains)
	case "logs":
		writeJSON(w, ctxData.Logs)
	case "ai":
		writeJSON(w, ctxData.AI)
	case "knowledge":
		writeJSON(w, ctxData.Knowledge)
	case "doctor":
		writeJSON(w, map[string]any{
			"health":  ctxData.Health,
			"project": ctxData.Project,
			"git":     ctxData.Git,
			"plugins": ctxData.Plugins,
			"checks": []map[string]string{
				{"label": "Dependencies", "status": statusFor(len(ctxData.Dependencies.Backend)+len(ctxData.Dependencies.Frontend)+len(ctxData.Dependencies.Database) > 0, "warning"), "detail": strings.Join(ctxData.Dependencies.Backend, ", ")},
				{"label": "Routes", "status": statusFor(len(ctxData.Routes.API) > 0, "passed"), "detail": fmt.Sprintf("%d routes indexed", len(ctxData.Routes.API))},
				{"label": "Git", "status": statusFor(ctxData.Git.Branch != "", "warning"), "detail": ctxData.Git.Branch},
				{"label": "Health", "status": statusFor(ctxData.Health.Score >= 40, "warning"), "detail": fmt.Sprintf("Score %d", ctxData.Health.Score)},
				{"label": "Infrastructure", "status": statusFor(len(ctxData.Infrastructure.Tools) > 0, "warning"), "detail": fmt.Sprintf("%d tools checked", len(ctxData.Infrastructure.Tools))},
			},
			"recommendations": []string{
				"Review dependency surface and plugins.",
				"Refresh route metadata after new saves.",
				"Keep the generated context in sync with the latest changes.",
			},
		})
	case "settings":
		writeJSON(w, map[string]any{
			"workspace": ctxData.Workspace,
			"index":     ctxData.Index,
			"cache":     ctxData.Cache,
		})
	case "mcp":
		writeJSON(w, ctxData.MCP)
	case "routes":
		writeJSON(w, ctxData.Routes)
	case "dependencies":
		writeJSON(w, ctxData.Dependencies)
	case "health":
		writeJSON(w, ctxData.Health)
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

func (s *APIServer) simulateTask(taskID, scope, name, detail string, done func() error) {
	phases := []struct {
		progress int
		state    string
		detail   string
	}{
		{20, "Running", detail},
		{45, "Running", "Installing dependencies and generating configs"},
		{70, "Running", "Applying validation and platform registration"},
		{100, "Done", "Task completed successfully"},
	}

	for _, phase := range phases {
		time.Sleep(450 * time.Millisecond)
		_ = s.store.UpdateTask(context.Background(), taskID, phase.progress, phase.state, phase.detail)
		s.hub.broadcast(state.TaskEvent{TaskID: taskID, Name: name, Progress: phase.progress, State: phase.state, Detail: phase.detail, Scope: scope})
	}

	if done != nil {
		_ = done()
	}
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

type TaskHub struct {
	log     zerolog.Logger
	mu      sync.Mutex
	clients map[*websocket.Conn]struct{}
}

func NewTaskHub(log zerolog.Logger) *TaskHub {
	return &TaskHub{
		log:     log,
		clients: map[*websocket.Conn]struct{}{},
	}
}

func (h *TaskHub) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		return
	}
	h.mu.Lock()
	h.clients[conn] = struct{}{}
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		_ = conn.Close(websocket.StatusNormalClosure, "")
	}()

	ctx := r.Context()
	for {
		var ignore map[string]any
		if err := wsjson.Read(ctx, conn, &ignore); err != nil {
			return
		}
	}
}

func (h *TaskHub) broadcast(event state.TaskEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.clients {
		_ = wsjson.Write(context.Background(), conn, event)
	}
}

func (h *TaskHub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.clients {
		_ = conn.Close(websocket.StatusNormalClosure, "")
		delete(h.clients, conn)
	}
}
