package state

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type BootstrapDecisionData struct {
	Destination          string `json:"destination"`
	SetupRequired        bool   `json:"setupRequired"`
	Authenticated        bool   `json:"authenticated"`
	InstallationProgress int    `json:"installationProgress"`
	Message              string `json:"message"`
}

type KnowledgeArticleData struct {
	Title    string `json:"title"`
	Question string `json:"question"`
	Summary  string `json:"summary"`
	Link     string `json:"link"`
}

type SetupStepData struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type SetupChecklistData struct {
	Label  string `json:"label"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

type SetupInstallationTypeData struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Detail string `json:"detail"`
}

type SetupWizardDataStore struct {
	Steps             []SetupStepData            `json:"steps"`
	Checks            []SetupChecklistData       `json:"checks"`
	Providers         []string                   `json:"providers"`
	InstallationTypes []SetupInstallationTypeData `json:"installationTypes"`
	KnowledgeArticles []KnowledgeArticleData     `json:"knowledgeArticles"`
	TaskID            string                     `json:"taskId"`
}

type LoginDataStore struct {
	Branding         string               `json:"branding"`
	Subtitle         string               `json:"subtitle"`
	SupportEmail     string               `json:"supportEmail"`
	KnowledgeArticles []KnowledgeArticleData `json:"knowledgeArticles"`
}

type MetricData struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Detail string `json:"detail"`
	Trend  string `json:"trend"`
	Tone   string `json:"tone"`
}

type SectionItemData struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Tone  string `json:"tone,omitempty"`
}

type SectionCardData struct {
	Title  string            `json:"title"`
	Detail string            `json:"detail"`
	Items  []SectionItemData `json:"items"`
}

type ActivityItemData struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	When   string `json:"when"`
	Tone   string `json:"tone"`
}

type TaskItemData struct {
	Title    string `json:"title"`
	Progress int    `json:"progress"`
	State    string `json:"state"`
	Detail   string `json:"detail"`
}

type RecommendationItemData struct {
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Reason string `json:"reason"`
}

type DashboardDataStore struct {
	Server          string                 `json:"server"`
	Notifications   int                    `json:"notifications"`
	Metrics         []MetricData           `json:"metrics"`
	Sections        []SectionCardData      `json:"sections"`
	Activity        []ActivityItemData     `json:"activity"`
	Tasks           []TaskItemData         `json:"tasks"`
	Recommendations []RecommendationItemData `json:"recommendations"`
	Knowledge       []KnowledgeArticleData `json:"knowledge"`
}

type ProjectData struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	Repository  string `json:"repository"`
	Environment string `json:"environment"`
	Status      string `json:"status"`
	UpdatedAt   string `json:"updatedAt"`
	Progress    int    `json:"progress"`
}

type ProjectDeploymentData struct {
	ID          string `json:"id"`
	Version     string `json:"version"`
	Status      string `json:"status"`
	Environment string `json:"environment"`
	DeployedAt  string `json:"deployedAt"`
	Note        string `json:"note"`
}

type ProjectDomainData struct {
	Host   string `json:"host"`
	SSL    string `json:"ssl"`
	Target string `json:"target"`
}

type ProjectLogData struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Time   string `json:"time"`
	Tone   string `json:"tone"`
}

type ProjectEnvironmentData struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	Visibility string `json:"visibility"`
}

type ProjectDatabaseData struct {
	Name   string `json:"name"`
	Engine string `json:"engine"`
	Status string `json:"status"`
}

type ProjectDetailDataStore struct {
	Project         ProjectData               `json:"project"`
	Repository      map[string]string         `json:"repository"`
	Environment     []ProjectEnvironmentData   `json:"environment"`
	Deployments     []ProjectDeploymentData    `json:"deployments"`
	Domains         []ProjectDomainData        `json:"domains"`
	Logs            []ProjectLogData           `json:"logs"`
	Database        ProjectDatabaseData        `json:"database"`
	Overview        []SectionCardData          `json:"overview"`
	Activity        []ActivityItemData         `json:"activity"`
	Tasks           []TaskItemData             `json:"tasks"`
	Recommendations []RecommendationItemData   `json:"recommendations"`
	Knowledge       []KnowledgeArticleData     `json:"knowledge"`
}

type ProjectFormDataStore struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	Repository   string `json:"repository"`
	Branch       string `json:"branch"`
	Environment  string `json:"environment"`
	Owner        string `json:"owner"`
	DeployTarget string `json:"deployTarget"`
	HealthCheck  string `json:"healthCheck"`
	AutoDeploy   bool   `json:"autoDeploy"`
	Domains      string `json:"domains"`
}

type ProjectFormOptionsData struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	SubmitLabel string `json:"submitLabel"`
}

type SettingsItemData struct {
	Label  string `json:"label"`
	Value  string `json:"value"`
	Detail string `json:"detail"`
}

type SettingsSectionData struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Items       []SettingsItemData `json:"items"`
}

type SettingsDataStore struct {
	Sections       []SettingsSectionData `json:"sections"`
	Preferences    []SettingsItemData    `json:"preferences"`
	Security       []SettingsItemData    `json:"security"`
	AIProviders    []SettingsItemData    `json:"aiProviders"`
	MCPIntegrations []SettingsItemData   `json:"mcpIntegrations"`
}

type ChecklistData struct {
	Label  string `json:"label"`
	Detail string `json:"detail"`
	Status string `json:"status"`
}

type DevCenterSectionDataStore struct {
	Key         string            `json:"key"`
	Title       string            `json:"title"`
	Subtitle    string            `json:"subtitle"`
	Description string            `json:"description"`
	Bullets     []string          `json:"bullets"`
	Examples    []string          `json:"examples"`
	Checklist   []ChecklistData   `json:"checklist"`
	Knowledge   []KnowledgeArticleData `json:"knowledge"`
	CodeSample  string            `json:"codeSample"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type DevCenterSectionSummaryData struct {
	Key      string `json:"key"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

type AuthResponseData struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type LoginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

type UserRecordData struct {
	Email    string
	Password string
	Role     string
	Remember bool
}

type SetupRequestData struct {
	LicenseAccepted  bool   `json:"licenseAccepted"`
	AdminName        string `json:"adminName"`
	AdminEmail       string `json:"adminEmail"`
	AdminPassword    string `json:"adminPassword"`
	Provider         string `json:"provider"`
	InstallationType string `json:"installationType"`
}

type TaskRecordData struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Scope    string `json:"scope"`
	Progress int    `json:"progress"`
	State    string `json:"state"`
	Detail   string `json:"detail"`
}

type StoreDB struct {
	db *sql.DB
}

func NewDB(path string) (*StoreDB, error) {
	if path == "" {
		path = "configs/devserver.db"
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if store, err := openStore("file:///"+filepath.ToSlash(absPath)+"?mode=rwc&cache=shared"); err == nil {
		return store, nil
	}
	return openStore(":memory:")
}

func openStore(dsn string) (*StoreDB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	store := &StoreDB{db: db}
	if err := store.init(context.Background()); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *StoreDB) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *StoreDB) init(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS bootstrap_state (id INTEGER PRIMARY KEY CHECK (id = 1), setup_required INTEGER NOT NULL, authenticated INTEGER NOT NULL, installation_progress INTEGER NOT NULL, message TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS projects (slug TEXT PRIMARY KEY, name TEXT NOT NULL, description TEXT NOT NULL, owner TEXT NOT NULL, repository TEXT NOT NULL, environment TEXT NOT NULL, status TEXT NOT NULL, updated_at TEXT NOT NULL, progress INTEGER NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS project_env (project_slug TEXT NOT NULL, "key" TEXT NOT NULL, value TEXT NOT NULL, visibility TEXT NOT NULL, PRIMARY KEY (project_slug, "key"));`,
		`CREATE TABLE IF NOT EXISTS deployments (id TEXT PRIMARY KEY, project_slug TEXT NOT NULL, version TEXT NOT NULL, status TEXT NOT NULL, environment TEXT NOT NULL, deployed_at TEXT NOT NULL, note TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS project_logs (id TEXT PRIMARY KEY, project_slug TEXT NOT NULL, title TEXT NOT NULL, detail TEXT NOT NULL, "time" TEXT NOT NULL, tone TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS users (email TEXT PRIMARY KEY, password TEXT NOT NULL, role TEXT NOT NULL, remember INTEGER NOT NULL, created_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS tasks (id TEXT PRIMARY KEY, scope TEXT NOT NULL, name TEXT NOT NULL, progress INTEGER NOT NULL, state TEXT NOT NULL, detail TEXT NOT NULL, updated_at TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS settings (section TEXT NOT NULL, "key" TEXT NOT NULL, value TEXT NOT NULL, detail TEXT NOT NULL, PRIMARY KEY (section, "key"));`,
		`CREATE TABLE IF NOT EXISTS ai_providers (name TEXT PRIMARY KEY, status TEXT NOT NULL, enabled INTEGER NOT NULL, detail TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS mcp_integrations (name TEXT PRIMARY KEY, status TEXT NOT NULL, detail TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS devcenter_sections ("key" TEXT PRIMARY KEY, title TEXT NOT NULL, subtitle TEXT NOT NULL, description TEXT NOT NULL, bullets_json TEXT NOT NULL, examples_json TEXT NOT NULL, checklist_json TEXT NOT NULL, knowledge_json TEXT NOT NULL, code_sample TEXT NOT NULL, metadata_json TEXT NOT NULL);`,
		`CREATE TABLE IF NOT EXISTS knowledge_articles (scope TEXT NOT NULL, title TEXT NOT NULL, question TEXT NOT NULL, summary TEXT NOT NULL, link TEXT NOT NULL);`,
	}
	for _, stmt := range stmts {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}
	return s.seed(ctx)
}

func (s *StoreDB) seed(ctx context.Context) error {
	if err := s.seedBootstrap(ctx); err != nil {
		return err
	}
	if err := s.seedProjects(ctx); err != nil {
		return err
	}
	if err := s.seedSettings(ctx); err != nil {
		return err
	}
	if err := s.seedTasks(ctx); err != nil {
		return err
	}
	if err := s.seedIntegrations(ctx); err != nil {
		return err
	}
	if err := s.seedKnowledge(ctx); err != nil {
		return err
	}
	if err := s.seedDevCenter(ctx); err != nil {
		return err
	}
	return nil
}

func (s *StoreDB) seedBootstrap(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM bootstrap_state`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err := s.db.ExecContext(ctx, `INSERT INTO bootstrap_state (id, setup_required, authenticated, installation_progress, message) VALUES (1, 1, 0, 32, ?)`, "Checking installation prerequisites and state store...")
	return err
}

func (s *StoreDB) seedProjects(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM projects`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	projects := []ProjectData{
		{Slug: "atlas-commerce", Name: "Atlas Commerce", Description: "Customer storefront, checkout, and internal ops dashboard.", Owner: "Ava Patel", Repository: "git@github.com:devserver/atlas-commerce.git", Environment: "production", Status: "healthy", UpdatedAt: "5 minutes ago", Progress: 88},
		{Slug: "docs-portal", Name: "Docs Portal", Description: "Product docs and knowledge hub for the platform.", Owner: "Milo Chen", Repository: "git@github.com:devserver/docs-portal.git", Environment: "staging", Status: "warning", UpdatedAt: "21 minutes ago", Progress: 63},
		{Slug: "internal-api", Name: "Internal API", Description: "Internal APIs powering auth, audit logs, and billing hooks.", Owner: "Priya Singh", Repository: "git@github.com:devserver/internal-api.git", Environment: "production", Status: "maintenance", UpdatedAt: "1 hour ago", Progress: 74},
		{Slug: "client-preview", Name: "Client Preview", Description: "Ephemeral preview environments for client-facing builds.", Owner: "Jordan Lee", Repository: "git@github.com:devserver/client-preview.git", Environment: "preview", Status: "blocked", UpdatedAt: "Yesterday", Progress: 29},
	}
	for _, project := range projects {
		if _, err := s.db.ExecContext(ctx, `INSERT INTO projects (slug, name, description, owner, repository, environment, status, updated_at, progress) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			project.Slug, project.Name, project.Description, project.Owner, project.Repository, project.Environment, project.Status, project.UpdatedAt, project.Progress); err != nil {
			return err
		}
	}
	for _, project := range projects {
		envs := []ProjectEnvironmentData{
			{Key: "APP_ENV", Value: project.Environment, Visibility: "public"},
			{Key: "LOG_LEVEL", Value: "info", Visibility: "public"},
			{Key: "DATABASE_URL", Value: "postgres://platform", Visibility: "secret"},
			{Key: "REDIS_URL", Value: "redis://cache", Visibility: "secret"},
		}
		for _, env := range envs {
			_, _ = s.db.ExecContext(ctx, `INSERT INTO project_env (project_slug, "key", value, visibility) VALUES (?, ?, ?, ?)`, project.Slug, env.Key, env.Value, env.Visibility)
		}
		_, _ = s.db.ExecContext(ctx, `INSERT INTO deployments (id, project_slug, version, status, environment, deployed_at, note) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			"dep-"+project.Slug+"-1", project.Slug, "v2.3.1", "succeeded", "production", "2 hours ago", "Smoke tests passed")
		_, _ = s.db.ExecContext(ctx, `INSERT INTO project_logs (id, project_slug, title, detail, "time", tone) VALUES (?, ?, ?, ?, ?, ?)`,
			"log-"+project.Slug+"-1", project.Slug, "Deployment validated", "The latest rollout completed smoke checks.", "4 minutes ago", "success")
	}
	return nil
}

func (s *StoreDB) seedSettings(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM settings`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	rows := []struct{ section, key, value, detail string }{
		{"platform", "default_environment", "production", "New projects inherit this environment unless overridden."},
		{"platform", "deployment_mode", "Rolling", "Updates roll through the fleet with health checks."},
		{"platform", "time_zone", "Asia/Kolkata", "Used for audit logs, tasks, and scheduled jobs."},
		{"notifications", "email_alerts", "Enabled", "Warnings and failures go to the admin team."},
		{"notifications", "task_summaries", "Every 6 hours", "Digest mode keeps noise manageable."},
		{"notifications", "slack_bridge", "Connected", "Workspace bridge currently points to #devserver-ops."},
		{"preferences", "theme", "Dark / Graphite", "The product ships in a calm high-contrast palette."},
		{"preferences", "default_shell", "/bin/bash", "Used when launching terminal sessions from projects."},
		{"preferences", "editor_format", "Prettier + ESLint", "Formatting presets are surfaced in project configs."},
		{"security", "session_ttl", "12 hours", "Admins can reduce this later without changing the UI."},
		{"security", "mfa_enforcement", "Recommended", "Enable for admin accounts before production use."},
		{"security", "audit_retention", "90 days", "Audit logs are retained on the platform database."},
	}
	for _, row := range rows {
		if _, err := s.db.ExecContext(ctx, `INSERT INTO settings (section, "key", value, detail) VALUES (?, ?, ?, ?)`, row.section, row.key, row.value, row.detail); err != nil {
			return err
		}
	}
	return nil
}

func (s *StoreDB) seedTasks(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM tasks`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	tasks := []TaskRecordData{
		{ID: "task-boot-1", Name: "Install nginx", Scope: "bootstrap", Progress: 100, State: "Done", Detail: "Completed and verified"},
		{ID: "task-boot-2", Name: "Configure postgres", Scope: "bootstrap", Progress: 72, State: "Running", Detail: "Creating database roles"},
		{ID: "task-boot-3", Name: "Prepare redis", Scope: "bootstrap", Progress: 38, State: "Queued", Detail: "Waiting for package source"},
	}
	for _, task := range tasks {
		if _, err := s.db.ExecContext(ctx, `INSERT INTO tasks (id, scope, name, progress, state, detail, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, task.ID, task.Scope, task.Name, task.Progress, task.State, task.Detail, time.Now().UTC().Format(time.RFC3339)); err != nil {
			return err
		}
	}
	return nil
}

func (s *StoreDB) seedIntegrations(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ai_providers`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := s.db.ExecContext(ctx, `INSERT INTO ai_providers (name, status, enabled, detail) VALUES
			('OpenAI', 'connected', 1, 'Configured and ready'),
			('Anthropic', 'available', 0, 'Can be enabled from settings'),
			('Local Model', 'available', 0, 'Local provider stub for offline testing'),
			('None', 'disabled', 0, 'Skip AI for now')`); err != nil {
			return err
		}
	}
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM mcp_integrations`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		if _, err := s.db.ExecContext(ctx, `INSERT INTO mcp_integrations (name, status, detail) VALUES
			('GitHub', 'connected', 'Repository access and pull request metadata'),
			('Slack', 'connected', 'Notifications and incident routing'),
			('Linear', 'available', 'Issue tracking integration stub'),
			('Grafana', 'available', 'Observability integration stub')`); err != nil {
			return err
		}
	}
	return nil
}

func (s *StoreDB) seedKnowledge(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM knowledge_articles`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	articles := map[string][]KnowledgeArticleData{
		"setup": {
			{Title: "What is bootstrap?", Question: "Why does DevServer need a setup wizard?", Summary: "Bootstrap is the controlled path from a clean machine to a managed platform.", Link: "#learn-bootstrap"},
			{Title: "System scan basics", Question: "What are we checking before install?", Summary: "We validate OS, memory, storage, networking, and service prerequisites.", Link: "#learn-scan"},
		},
		"login": {
			{Title: "JWT sessions", Question: "How does authentication stay stateless?", Summary: "The UI will post credentials, receive a token, and refresh transparently.", Link: "#learn-jwt"},
			{Title: "RBAC", Question: "How are permissions grouped?", Summary: "Roles like Admin, Developer, and Viewer are derived from the backend policy.", Link: "#learn-rbac"},
		},
		"dashboard": {
			{Title: "Projects", Question: "How do projects differ from servers?", Summary: "Projects bundle repository, environment, deployment, and monitoring state.", Link: "#learn-projects"},
			{Title: "Terminal", Question: "Will commands run locally or remotely?", Summary: "The executor layer will eventually support SSH and dry-run execution.", Link: "#learn-terminal"},
		},
	}
	for scope, list := range articles {
		for _, article := range list {
			if _, err := s.db.ExecContext(ctx, `INSERT INTO knowledge_articles (scope, title, question, summary, link) VALUES (?, ?, ?, ?, ?)`, scope, article.Title, article.Question, article.Summary, article.Link); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *StoreDB) seedDevCenter(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM devcenter_sections`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	type seed struct {
		key, title, subtitle, description, code string
		bullets []string
		examples []string
		checklist []ChecklistData
		metadata map[string]string
	}
	seeds := []seed{
		{key: "architecture", title: "Architecture", subtitle: "Understand how DevServer is structured end to end.", description: "A single control plane coordinates server bootstrap, project metadata, tasks, and automation.", code: "const task = await createTask({\n  name: 'Install nginx',\n  progress: 72,\n  rollback: 'Remove package and restore previous config',\n})", bullets: []string{"Backend services stay in internal/ and expose testable interfaces.", "The web app consumes mock facades that can be swapped for APIs later.", "Each long-running workflow is modeled as a task sequence with rollback support."}, examples: []string{"Bootstrap -> Setup -> Login -> Dashboard", "Project -> Deployments -> Logs -> Settings"}, checklist: []ChecklistData{{Label: "Shell layout", Detail: "A consistent shell wraps product screens.", Status: "passed"}, {Label: "Service boundaries", Detail: "Mock services keep the UI replaceable.", Status: "passed"}, {Label: "Task engine", Detail: "Progress and rollback hooks are represented in the UI.", Status: "warning"}}, metadata: map[string]string{"recommendation": "Keep the API shape stable while swapping mocks for real services."}},
		{key: "tasks", title: "Tasks", subtitle: "Track installation, deployment, and maintenance workflows.", description: "Every meaningful operation is represented as a task with progress, status, and rollback intent.", code: "Task: Install PostgreSQL\nProgress: 82%\nRollback: Stop service and restore previous data path", bullets: []string{"Tasks are ordered.", "Tasks report progress.", "Tasks can roll back.", "Tasks can surface warnings before execution."}, examples: []string{"Install Redis", "Generate SSL", "Deploy project release"}, checklist: []ChecklistData{{Label: "Queued tasks", Detail: "Awaiting execution or approval.", Status: "passed"}, {Label: "Running tasks", Detail: "Progress is shown in the queue.", Status: "warning"}, {Label: "Failed tasks", Detail: "Failures should capture rollback context.", Status: "failed"}}, metadata: map[string]string{"model": "Task progress should be driven by backend events, not UI timers."}},
		{key: "knowledge", title: "Knowledge", subtitle: "Teach while you build with documentation and best practices.", description: "Every product page can explain why a feature exists, not just what it does.", code: "Why Brotli?\nIt reduces transfer size for static assets and keeps page loads fast.", bullets: []string{"Include examples.", "Explain tradeoffs.", "Show practical next steps.", "Keep the copy close to the action."}, examples: []string{"What is a reverse proxy?", "Why cache at the edge?", "How does RBAC shape workflows?"}, checklist: []ChecklistData{{Label: "Guidance cards", Detail: "Best practice content sits beside the workflow.", Status: "passed"}, {Label: "Inline learning", Detail: "Contextual docs are visible on every major page.", Status: "passed"}}},
		{key: "dependencies", title: "Dependencies", subtitle: "Map module and service dependencies explicitly.", description: "Dependency visibility prevents hidden coupling between services, tasks, and project workflows.", code: "registry.Register(&modules.Nginx{})\nregistry.Register(&modules.Redis{})", bullets: []string{"Modules register in one place.", "Tasks depend on prerequisites.", "UI pages consume facades, not concrete APIs."}, examples: []string{"Nginx depends on system config", "Projects depend on platform state", "Deployments depend on healthy services"}, checklist: []ChecklistData{{Label: "Registry mapping", Detail: "Modules should be registered, not hardcoded.", Status: "passed"}, {Label: "Contract isolation", Detail: "Pages call services and hooks only.", Status: "passed"}}},
		{key: "api", title: "API", subtitle: "Plan the eventual API surface without coupling the UI to it yet.", description: "The frontend is built around service interfaces so a future API layer can slot in beneath it.", code: "export async function listProjects() {\n  return fetch('/api/projects').then((res) => res.json())\n}", bullets: []string{"Mock service now.", "API later.", "No UI refactor required.", "Keep response shapes stable."}, examples: []string{"GET /projects", "GET /projects/:slug", "POST /projects", "PATCH /settings"}, checklist: []ChecklistData{{Label: "Interface-first", Detail: "Service contracts describe the UI data needs.", Status: "passed"}, {Label: "Mock replacement", Detail: "Real API calls can live behind the same functions.", Status: "warning"}}},
		{key: "database", title: "Database", subtitle: "Understand the platform data model and how it supports the UI.", description: "A separate platform database stores users, projects, deployments, tasks, notifications, and settings.", code: "CREATE TABLE projects (\n  id UUID PRIMARY KEY,\n  slug TEXT UNIQUE NOT NULL\n);", bullets: []string{"Keep platform metadata separate.", "Store audit events centrally.", "Use the same shape from mock to real data."}, examples: []string{"users", "projects", "deployments", "tasks", "notifications", "settings"}, checklist: []ChecklistData{{Label: "Schema clarity", Detail: "Tables should mirror product concerns.", Status: "passed"}, {Label: "Audit trail", Detail: "Change history should be queryable.", Status: "passed"}}},
		{key: "project-doctor", title: "Project Doctor", subtitle: "Inspect a project and surface risks before deployment.", description: "Doctor mode is a guided diagnostic surface for repository, environment, service, and deployment issues.", code: "Doctor checks:\n- repo connectivity\n- env completeness\n- service health\n- SSL expiry", bullets: []string{"Scan for missing env vars.", "Check runtime readiness.", "Validate deploy target health.", "Show actionable fixes."}, examples: []string{"Repository connected", "Database reachable", "SSL expiring soon"}, checklist: []ChecklistData{{Label: "Repository", Detail: "Connected and fetchable.", Status: "passed"}, {Label: "Environment", Detail: "One secret is missing from the current mock snapshot.", Status: "warning"}, {Label: "Deployment target", Detail: "Healthy and ready to accept rollout.", Status: "passed"}}},
	}
	for _, sct := range seeds {
		b, _ := json.Marshal(sct.bullets)
		e, _ := json.Marshal(sct.examples)
		c, _ := json.Marshal(sct.checklist)
		k, _ := json.Marshal([]KnowledgeArticleData{})
		m, _ := json.Marshal(sct.metadata)
		if _, err := s.db.ExecContext(ctx, `INSERT INTO devcenter_sections ("key", title, subtitle, description, bullets_json, examples_json, checklist_json, knowledge_json, code_sample, metadata_json) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			sct.key, sct.title, sct.subtitle, sct.description, b, e, c, k, sct.code, m); err != nil {
			return err
		}
	}
	return nil
}

func (s *StoreDB) BootstrapDecision(ctx context.Context) (BootstrapDecisionData, error) {
	var data BootstrapDecisionData
	var setupRequired, authenticated int
	if err := s.db.QueryRowContext(ctx, `SELECT setup_required, authenticated, installation_progress, message FROM bootstrap_state WHERE id = 1`).Scan(&setupRequired, &authenticated, &data.InstallationProgress, &data.Message); err != nil {
		return BootstrapDecisionData{}, err
	}
	data.SetupRequired = setupRequired == 1
	data.Authenticated = authenticated == 1
	switch {
	case data.SetupRequired:
		data.Destination = "setup"
	case data.Authenticated:
		data.Destination = "dashboard"
	default:
		data.Destination = "login"
	}
	return data, nil
}

func (s *StoreDB) SetupWizardData(ctx context.Context) (SetupWizardDataStore, error) {
	articles, err := s.knowledgeByScope(ctx, "setup")
	if err != nil {
		return SetupWizardDataStore{}, err
	}
	return SetupWizardDataStore{
		Steps: []SetupStepData{
			{Key: "welcome", Title: "Welcome", Subtitle: "Meet DevServer and confirm the baseline"},
			{Key: "license", Title: "License", Subtitle: "Review licensing and usage terms"},
			{Key: "scan", Title: "System Scan", Subtitle: "Validate platform readiness"},
			{Key: "admin", Title: "Administrator", Subtitle: "Create the first admin account"},
			{Key: "ai", Title: "AI Providers", Subtitle: "Configure optional providers"},
			{Key: "install-type", Title: "Installation Type", Subtitle: "Choose a quick or advanced flow"},
			{Key: "summary", Title: "Summary", Subtitle: "Review before changes are applied"},
			{Key: "installing", Title: "Installing", Subtitle: "Provisioning is in progress"},
			{Key: "finished", Title: "Finished", Subtitle: "The platform is ready"},
		},
		Checks: []SetupChecklistData{
			{Label: "Ubuntu", Detail: "Detected 24.04 LTS baseline", Status: "passed"},
			{Label: "CPU", Detail: "8 cores available", Status: "passed"},
			{Label: "RAM", Detail: "16 GB available", Status: "passed"},
			{Label: "Disk", Detail: "120 GB free on root volume", Status: "warning"},
			{Label: "Internet", Detail: "Connectivity is available", Status: "passed"},
			{Label: "DNS", Detail: "Resolver responds to lookups", Status: "passed"},
			{Label: "Swap", Detail: "Swap is enabled and healthy", Status: "passed"},
			{Label: "SSH", Detail: "SSH service is reachable", Status: "passed"},
			{Label: "Firewall", Detail: "Default policy not yet applied", Status: "warning"},
			{Label: "Git", Detail: "git binary found", Status: "passed"},
			{Label: "Node", Detail: "Node can be installed later", Status: "pending"},
			{Label: "Redis", Detail: "Not installed yet", Status: "pending"},
			{Label: "Nginx", Detail: "Not installed yet", Status: "pending"},
			{Label: "Postgres", Detail: "Not installed yet", Status: "pending"},
			{Label: "PM2", Detail: "Optional process manager unavailable", Status: "pending"},
		},
		Providers: []string{"OpenAI", "Anthropic", "Local Model", "None"},
		InstallationTypes: []SetupInstallationTypeData{{Key: "quick", Label: "Quick install", Detail: "Recommended defaults with sensible production-ready settings."}, {Key: "advanced", Label: "Advanced install", Detail: "Choose every module and inspect each generated step."}},
		KnowledgeArticles: articles,
		TaskID: "task-boot-setup",
	}, nil
}

func (s *StoreDB) LoginData(ctx context.Context) (LoginDataStore, error) {
	articles, err := s.knowledgeByScope(ctx, "login")
	if err != nil {
		return LoginDataStore{}, err
	}
	return LoginDataStore{
		Branding: "DevServer",
		Subtitle: "Control server bootstrap, deployments, and infrastructure workflows from one place.",
		SupportEmail: "support@devserver.local",
		KnowledgeArticles: articles,
	}, nil
}

func (s *StoreDB) DashboardData(ctx context.Context) (DashboardDataStore, error) {
	projects, err := s.ListProjects(ctx)
	if err != nil {
		return DashboardDataStore{}, err
	}
	tasks, err := s.Tasks(ctx)
	if err != nil {
		return DashboardDataStore{}, err
	}
	articles, err := s.knowledgeByScope(ctx, "dashboard")
	if err != nil {
		return DashboardDataStore{}, err
	}
	return DashboardDataStore{
		Server: "production-east-1",
		Notifications: 4,
		Metrics: []MetricData{{Label: "CPU", Value: "42%", Detail: "32% average over 15 min", Trend: "+4%", Tone: "success"}, {Label: "RAM", Value: "68%", Detail: "A little warm on two hosts", Trend: "+2%", Tone: "warning"}, {Label: "Disk", Value: "51%", Detail: "Plenty of headroom remains", Trend: "stable", Tone: "info"}, {Label: "Network", Value: "24 ms", Detail: "Median request path latency", Trend: "-7%", Tone: "accent"}},
		Sections: []SectionCardData{
			{Title: "Projects", Detail: "Active projects and their health.", Items: []SectionItemData{{Label: projects[0].Name, Value: "healthy", Tone: "success"}, {Label: projects[1].Name, Value: "deploying", Tone: "warning"}, {Label: projects[2].Name, Value: "stable", Tone: "info"}}},
			{Title: "Domains", Detail: "DNS and SSL coverage.", Items: []SectionItemData{{Label: "app.devserver.local", Value: "valid", Tone: "success"}, {Label: "api.devserver.local", Value: "renew in 22 days", Tone: "warning"}, {Label: "staging.devserver.local", Value: "healthy", Tone: "success"}}},
			{Title: "Services", Detail: "Runtime services and systemd state.", Items: []SectionItemData{{Label: "nginx", Value: "running", Tone: "success"}, {Label: "postgres", Value: "running", Tone: "success"}, {Label: "redis", Value: "degraded", Tone: "warning"}}},
			{Title: "Deployments", Detail: "Release history and rollout state.", Items: []SectionItemData{{Label: "Production", Value: "2 hours ago", Tone: "success"}, {Label: "Staging", Value: "12 minutes ago", Tone: "accent"}, {Label: "Preview", Value: "queued", Tone: "warning"}}},
		},
		Activity: []ActivityItemData{{Title: "Redis task queued", Detail: "Prepare memory policy and reload service", When: "2m ago", Tone: "warning"}, {Title: "Deploy approved", Detail: "Production deployment passed validation", When: "11m ago", Tone: "success"}, {Title: "SSL renewal warning", Detail: "One certificate will expire soon", When: "25m ago", Tone: "info"}},
		Tasks: tasks,
		Recommendations: []RecommendationItemData{{Title: "Enable SSL renewal", Detail: "One domain is nearing expiry and should be automated.", Reason: "Improves uptime and reduces manual maintenance."}, {Title: "Review redis memory limits", Detail: "A warning was detected in the monitoring snapshot.", Reason: "Keeps cache pressure from affecting the stack."}, {Title: "Add project metadata", Detail: "Projects can be richer if owners and environments are registered.", Reason: "Improves filtering and auditability."}},
		Knowledge: articles,
	}, nil
}

func (s *StoreDB) ListProjects(ctx context.Context) ([]ProjectData, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT slug, name, description, owner, repository, environment, status, updated_at, progress FROM projects ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ProjectData
	for rows.Next() {
		var item ProjectData
		if err := rows.Scan(&item.Slug, &item.Name, &item.Description, &item.Owner, &item.Repository, &item.Environment, &item.Status, &item.UpdatedAt, &item.Progress); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *StoreDB) ProjectDetail(ctx context.Context, slug string) (ProjectDetailDataStore, error) {
	var project ProjectData
	if err := s.db.QueryRowContext(ctx, `SELECT slug, name, description, owner, repository, environment, status, updated_at, progress FROM projects WHERE slug = ?`, slug).Scan(&project.Slug, &project.Name, &project.Description, &project.Owner, &project.Repository, &project.Environment, &project.Status, &project.UpdatedAt, &project.Progress); err != nil {
		return ProjectDetailDataStore{}, err
	}
	env, _ := s.projectEnvironment(ctx, slug)
	deployments, _ := s.projectDeployments(ctx, slug)
	logs, _ := s.projectLogs(ctx, slug)
	knowledge, _ := s.knowledgeByScope(ctx, "dashboard")
	return ProjectDetailDataStore{
		Project: project,
		Repository: map[string]string{"url": project.Repository, "branch": branchFor(project.Environment), "buildCommand": "pnpm build", "startCommand": "pnpm start"},
		Environment: env,
		Deployments: deployments,
		Domains: []ProjectDomainData{{Host: project.Slug + ".devserver.local", SSL: "Valid", Target: "nginx / root"}, {Host: "api." + project.Slug + ".devserver.local", SSL: "Renewing", Target: "internal api / upstream"}, {Host: "cdn." + project.Slug + ".devserver.local", SSL: "Pending", Target: "storage / assets"}},
		Logs: logs,
		Database: ProjectDatabaseData{Name: project.Slug + "-platform", Engine: "PostgreSQL 16", Status: dbStatusFor(project.Status)},
		Overview: []SectionCardData{{Title: "Ownership", Detail: "Project metadata and lifecycle status.", Items: []SectionItemData{{Label: "Owner", Value: project.Owner, Tone: "accent"}, {Label: "Environment", Value: project.Environment, Tone: "info"}, {Label: "Status", Value: project.Status, Tone: project.Status}}}, {Title: "Build", Detail: "Repository and release commands.", Items: []SectionItemData{{Label: "Repository", Value: project.Repository, Tone: "neutral"}, {Label: "Branch", Value: branchFor(project.Environment), Tone: "accent"}, {Label: "Deploy target", Value: "Primary fleet", Tone: "success"}}}},
		Activity: []ActivityItemData{{Title: "Project synced", Detail: "Repository state was refreshed from the latest snapshot.", When: "2m ago", Tone: "success"}, {Title: "Task queued", Detail: "A deployment workflow is waiting for approval.", When: "15m ago", Tone: "warning"}},
		Tasks: []TaskItemData{{Title: "Build project index", Progress: 100, State: "Done", Detail: "Metadata and environments are indexed."}, {Title: "Validate secrets", Progress: 68, State: "Running", Detail: "Checking secret references and mounts."}, {Title: "Prepare terminal session", Progress: 28, State: "Queued", Detail: "Waiting for a workspace slot."}},
		Recommendations: []RecommendationItemData{{Title: "Promote preview to staging", Detail: "This project has a solid deployment history and could reuse the same release path.", Reason: "Keeps the release train consistent."}, {Title: "Review SSL expiry windows", Detail: "One domain is renewing and should be validated before traffic shifts.", Reason: "Protects the project edge surface."}},
		Knowledge: knowledge,
	}, nil
}

func (s *StoreDB) ProjectFormData(ctx context.Context, slug string) (ProjectFormDataStore, error) {
	if slug == "" {
		return ProjectFormDataStore{
			Name:         "",
			Slug:         "",
			Description:  "",
			Repository:   "git@github.com:devserver/new-project.git",
			Branch:       "main",
			Environment:  "production",
			Owner:        "",
			DeployTarget: "Primary fleet",
			HealthCheck:  "/health",
			AutoDeploy:   true,
			Domains:      "app.devserver.local, api.devserver.local",
		}, nil
	}
	project, err := s.ProjectDetail(ctx, slug)
	if err != nil {
		return ProjectFormDataStore{}, err
	}
	return formFromProject(project.Project), nil
}

func (s *StoreDB) SaveProject(ctx context.Context, form ProjectFormDataStore) (ProjectData, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.db.ExecContext(ctx, `INSERT INTO projects (slug, name, description, owner, repository, environment, status, updated_at, progress) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET name = excluded.name, description = excluded.description, owner = excluded.owner, repository = excluded.repository, environment = excluded.environment, status = excluded.status, updated_at = excluded.updated_at, progress = excluded.progress`,
		form.Slug, form.Name, form.Description, form.Owner, form.Repository, form.Environment, "healthy", now, 100)
	if err != nil {
		return ProjectData{}, err
	}
	_, _ = s.db.ExecContext(ctx, `DELETE FROM project_env WHERE project_slug = ?`, form.Slug)
	for _, entry := range strings.Split(form.Domains, ",") {
		if trimmed := strings.TrimSpace(entry); trimmed != "" {
			_, _ = s.db.ExecContext(ctx, `INSERT INTO project_env (project_slug, "key", value, visibility) VALUES (?, ?, ?, ?)`, form.Slug, "DOMAIN", trimmed, "public")
		}
	}
	return ProjectData{Slug: form.Slug, Name: form.Name, Description: form.Description, Owner: form.Owner, Repository: form.Repository, Environment: form.Environment, Status: "healthy", UpdatedAt: "just now", Progress: 100}, nil
}

func (s *StoreDB) Settings(ctx context.Context) (SettingsDataStore, error) {
	return SettingsDataStore{
		Sections: []SettingsSectionData{{Title: "Platform defaults", Description: "Control global behavior for all projects and new deployments.", Items: []SettingsItemData{{Label: "Default environment", Value: "production", Detail: "New projects inherit this environment unless overridden."}, {Label: "Deployment mode", Value: "Rolling", Detail: "Updates roll through the fleet with health checks."}, {Label: "Time zone", Value: "Asia/Kolkata", Detail: "Used for audit logs, tasks, and scheduled jobs."}}}, {Title: "Notifications", Description: "Route operational updates to the right people.", Items: []SettingsItemData{{Label: "Email alerts", Value: "Enabled", Detail: "Warnings and failures go to the admin team."}, {Label: "Task summaries", Value: "Every 6 hours", Detail: "Digest mode keeps noise manageable."}, {Label: "Slack bridge", Value: "Connected", Detail: "Workspace bridge currently points to #devserver-ops."}}}},
		Preferences: []SettingsItemData{{Label: "Theme", Value: "Dark / Graphite", Detail: "The product ships in a calm high-contrast palette."}, {Label: "Default shell", Value: "/bin/bash", Detail: "Used when launching terminal sessions from projects."}, {Label: "Editor format", Value: "Prettier + ESLint", Detail: "Formatting presets are surfaced in project configs."}},
		Security: []SettingsItemData{{Label: "Session TTL", Value: "12 hours", Detail: "Admins can reduce this later without changing the UI."}, {Label: "MFA enforcement", Value: "Recommended", Detail: "Enable for admin accounts before production use."}, {Label: "Audit retention", Value: "90 days", Detail: "Audit logs are retained on the platform database."}},
		AIProviders: []SettingsItemData{{Label: "OpenAI", Value: "Connected", Detail: "Configured and ready"}, {Label: "Anthropic", Value: "Available", Detail: "Can be enabled from settings"}, {Label: "Local Model", Value: "Available", Detail: "Local provider stub for offline testing"}, {Label: "None", Value: "Disabled", Detail: "Skip AI for now"}},
		MCPIntegrations: []SettingsItemData{{Label: "GitHub", Value: "Connected", Detail: "Repository access and pull request metadata"}, {Label: "Slack", Value: "Connected", Detail: "Notifications and incident routing"}, {Label: "Linear", Value: "Available", Detail: "Issue tracking integration stub"}, {Label: "Grafana", Value: "Available", Detail: "Observability integration stub"}},
	}, nil
}

func (s *StoreDB) DevCenterSection(ctx context.Context, key string) (DevCenterSectionDataStore, error) {
	row := s.db.QueryRowContext(ctx, `SELECT title, subtitle, description, bullets_json, examples_json, checklist_json, knowledge_json, code_sample, metadata_json FROM devcenter_sections WHERE "key" = ?`, key)
	var section DevCenterSectionDataStore
	var bulletsJSON, examplesJSON, checklistJSON, knowledgeJSON, metadataJSON string
	if err := row.Scan(&section.Title, &section.Subtitle, &section.Description, &bulletsJSON, &examplesJSON, &checklistJSON, &knowledgeJSON, &section.CodeSample, &metadataJSON); err != nil {
		return DevCenterSectionDataStore{}, err
	}
	section.Key = key
	_ = json.Unmarshal([]byte(bulletsJSON), &section.Bullets)
	_ = json.Unmarshal([]byte(examplesJSON), &section.Examples)
	_ = json.Unmarshal([]byte(checklistJSON), &section.Checklist)
	_ = json.Unmarshal([]byte(knowledgeJSON), &section.Knowledge)
	_ = json.Unmarshal([]byte(metadataJSON), &section.Metadata)
	return section, nil
}

func (s *StoreDB) DevCenterSections(ctx context.Context) ([]DevCenterSectionSummaryData, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT "key", title, subtitle FROM devcenter_sections ORDER BY title`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []DevCenterSectionSummaryData
	for rows.Next() {
		var item DevCenterSectionSummaryData
		if err := rows.Scan(&item.Key, &item.Title, &item.Subtitle); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *StoreDB) knowledgeByScope(ctx context.Context, scope string) ([]KnowledgeArticleData, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT title, question, summary, link FROM knowledge_articles WHERE scope = ? ORDER BY title`, scope)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []KnowledgeArticleData
	for rows.Next() {
		var item KnowledgeArticleData
		if err := rows.Scan(&item.Title, &item.Question, &item.Summary, &item.Link); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *StoreDB) KnowledgeArticles(ctx context.Context, scope string) ([]KnowledgeArticleData, error) {
	if strings.TrimSpace(scope) == "" {
		scope = "dashboard"
	}
	return s.knowledgeByScope(ctx, scope)
}

func (s *StoreDB) Tasks(ctx context.Context) ([]TaskItemData, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT name, progress, state, detail FROM tasks ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []TaskItemData
	for rows.Next() {
		var item TaskItemData
		if err := rows.Scan(&item.Title, &item.Progress, &item.State, &item.Detail); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *StoreDB) CreateTask(ctx context.Context, scope, name, detail string) (TaskRecordData, error) {
	id := fmt.Sprintf("task-%d", time.Now().UnixNano())
	task := TaskRecordData{ID: id, Name: name, Scope: scope, Progress: 0, State: "Running", Detail: detail}
	if _, err := s.db.ExecContext(ctx, `INSERT INTO tasks (id, scope, name, progress, state, detail, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, task.ID, task.Scope, task.Name, task.Progress, task.State, task.Detail, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return TaskRecordData{}, err
	}
	return task, nil
}

func (s *StoreDB) UpdateTask(ctx context.Context, taskID string, progress int, state, detail string) error {
	_, err := s.db.ExecContext(ctx, `UPDATE tasks SET progress = ?, state = ?, detail = ?, updated_at = ? WHERE id = ?`, progress, state, detail, time.Now().UTC().Format(time.RFC3339), taskID)
	return err
}

func (s *StoreDB) GetTask(ctx context.Context, taskID string) (TaskRecordData, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id, scope, name, progress, state, detail FROM tasks WHERE id = ?`, taskID)
	var task TaskRecordData
	if err := row.Scan(&task.ID, &task.Scope, &task.Name, &task.Progress, &task.State, &task.Detail); err != nil {
		return TaskRecordData{}, err
	}
	return task, nil
}

func (s *StoreDB) CompleteSetup(ctx context.Context, req SetupRequestData) error {
	_, err := s.db.ExecContext(ctx, `UPDATE bootstrap_state SET setup_required = 0, authenticated = 1, installation_progress = 100, message = ? WHERE id = 1`, "Setup complete. Ready for login.")
	if err != nil {
		return err
	}
	email := strings.TrimSpace(req.AdminEmail)
	password := req.AdminPassword
	if email == "" {
		email = "admin@admin.com"
	}
	if password == "" {
		password = "12345678"
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO users (email, password, role, remember, created_at) VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET password = excluded.password, role = excluded.role, remember = excluded.remember`,
		email, password, "superadmin", 1, time.Now().UTC().Format(time.RFC3339))
	return err
}

func (s *StoreDB) Auth(ctx context.Context, req LoginRequestData) (AuthResponseData, error) {
	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return AuthResponseData{}, errors.New("email and password are required")
	}
	if !strings.Contains(req.Email, "@") {
		return AuthResponseData{}, errors.New("enter a valid email address")
	}
	if len(req.Password) < 8 {
		return AuthResponseData{}, errors.New("password must be at least 8 characters")
	}
	var user UserRecordData
	if err := s.db.QueryRowContext(ctx, `SELECT email, password, role, remember FROM users WHERE email = ?`, strings.TrimSpace(req.Email)).Scan(&user.Email, &user.Password, &user.Role, &user.Remember); err != nil {
		return AuthResponseData{}, errors.New("invalid email or password")
	}
	if user.Password != req.Password {
		return AuthResponseData{}, errors.New("invalid email or password")
	}
	_, err := s.db.ExecContext(ctx, `UPDATE bootstrap_state SET authenticated = 1, setup_required = 0 WHERE id = 1`)
	if err != nil {
		return AuthResponseData{}, err
	}
	_, _ = s.db.ExecContext(ctx, `UPDATE users SET remember = ? WHERE email = ?`, boolToInt(req.Remember), strings.TrimSpace(req.Email))
	return AuthResponseData{Token: "devserver-token", Role: strings.Title(user.Role)}, nil
}

func (s *StoreDB) currentSetupTask(ctx context.Context) (string, error) {
	row := s.db.QueryRowContext(ctx, `SELECT id FROM tasks WHERE scope = 'setup' ORDER BY updated_at DESC LIMIT 1`)
	var id string
	if err := row.Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func (s *StoreDB) projectEnvironment(ctx context.Context, slug string) ([]ProjectEnvironmentData, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT "key", value, visibility FROM project_env WHERE project_slug = ? ORDER BY "key"`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ProjectEnvironmentData
	for rows.Next() {
		var item ProjectEnvironmentData
		if err := rows.Scan(&item.Key, &item.Value, &item.Visibility); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *StoreDB) projectDeployments(ctx context.Context, slug string) ([]ProjectDeploymentData, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, version, status, environment, deployed_at, note FROM deployments WHERE project_slug = ? ORDER BY deployed_at DESC`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ProjectDeploymentData
	for rows.Next() {
		var item ProjectDeploymentData
		if err := rows.Scan(&item.ID, &item.Version, &item.Status, &item.Environment, &item.DeployedAt, &item.Note); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *StoreDB) projectLogs(ctx context.Context, slug string) ([]ProjectLogData, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, title, detail, "time", tone FROM project_logs WHERE project_slug = ? ORDER BY "time" DESC`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ProjectLogData
	for rows.Next() {
		var item ProjectLogData
		if err := rows.Scan(&item.ID, &item.Title, &item.Detail, &item.Time, &item.Tone); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func formFromProject(project ProjectData) ProjectFormDataStore {
	return ProjectFormDataStore{
		Name: project.Name,
		Slug: project.Slug,
		Description: project.Description,
		Repository: project.Repository,
		Branch: branchFor(project.Environment),
		Environment: project.Environment,
		Owner: project.Owner,
		DeployTarget: "Primary fleet",
		HealthCheck: "/health",
		AutoDeploy: true,
		Domains: project.Slug + ".devserver.local, api." + project.Slug + ".devserver.local",
	}
}

func branchFor(environment string) string {
	if environment == "production" {
		return "main"
	}
	return "develop"
}

func dbStatusFor(status string) string {
	if status == "blocked" {
		return "syncing"
	}
	return "ready"
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
