export type Tone = 'neutral' | 'accent' | 'success' | 'warning' | 'danger' | 'info'

export interface NavItem {
  label: string
  href: string
  subtitle: string
  badge?: string
}

export interface NavGroup {
  title: string
  items: NavItem[]
}

export interface Metric {
  label: string
  value: string
  delta: string
  helper: string
  tone: Tone
}

export interface ChecklistItem {
  label: string
  detail: string
  status: 'healthy' | 'warning' | 'offline' | 'pending' | 'passed' | 'failed'
}

export interface TimelineItem {
  title: string
  detail: string
  time: string
  tone: Tone
}

export interface ModuleCard {
  name: string
  description: string
  version: string
  status: 'ready' | 'blocked' | 'pending'
  nextStep: string
}

export interface ServiceRow {
  name: string
  status: 'running' | 'stopped' | 'degraded'
  host: string
  port: string
  mode: string
  action: string
}

export interface TargetRow {
  name: string
  region: string
  host: string
  status: 'healthy' | 'warning' | 'critical'
  latency: string
  lastDeploy: string
}

export interface BackupJob {
  name: string
  schedule: string
  retention: string
  lastRun: string
  status: 'success' | 'warning' | 'failed'
  destination: string
}

export interface ConfigSection {
  title: string
  detail: string
  path: string
  value: string
}

export interface AlertItem {
  title: string
  severity: 'info' | 'warning' | 'critical'
  detail: string
  timestamp: string
}

export interface DashboardData {
  metrics: Metric[]
  checks: ChecklistItem[]
  timeline: TimelineItem[]
  modules: ModuleCard[]
}

export interface DoctorData {
  checks: ChecklistItem[]
  recommendations: string[]
}

export interface InstallData {
  modules: ModuleCard[]
  steps: TimelineItem[]
}

export interface ServiceData {
  services: ServiceRow[]
}

export interface DeployData {
  targets: TargetRow[]
  timeline: TimelineItem[]
}

export interface BackupData {
  jobs: BackupJob[]
  schedule: TimelineItem[]
}

export interface ConfigData {
  sections: ConfigSection[]
  env: ConfigSection[]
  preview: string
}

export interface MonitorData {
  alerts: AlertItem[]
  services: ChecklistItem[]
  hostHealth: Metric[]
}

export type BootstrapDestination = 'setup' | 'login' | 'dashboard'

export interface BootstrapDecision {
  destination: BootstrapDestination
  setupRequired: boolean
  authenticated: boolean
  installationProgress: number
  message: string
}

export interface WizardStep {
  key: string
  title: string
  subtitle: string
}

export interface SetupChecklistItem {
  label: string
  detail: string
  status: 'pending' | 'running' | 'passed' | 'warning' | 'failed'
}

export interface SetupWizardData {
  steps: WizardStep[]
  checks: SetupChecklistItem[]
  providers: string[]
  installationTypes: Array<{
    key: string
    label: string
    detail: string
  }>
  knowledgeArticles: KnowledgeArticle[]
  taskId?: string
}

export interface LoginData {
  branding: string
  subtitle: string
  supportEmail: string
  knowledgeArticles: KnowledgeArticle[]
}

export interface DashboardMetric {
  label: string
  value: string
  detail: string
  trend: string
  tone: Tone
}

export interface DashboardSectionCard {
  title: string
  detail: string
  items: Array<{
    label: string
    value: string
    tone?: Tone
  }>
}

export interface ActivityItem {
  title: string
  detail: string
  when: string
  tone: Tone
}

export interface TaskItem {
  title: string
  progress: number
  state: 'Queued' | 'Running' | 'Done' | 'Waiting'
  detail: string
}

export interface RecommendationItem {
  title: string
  detail: string
  reason: string
}

export interface KnowledgeArticle {
  title: string
  question: string
  summary: string
  link: string
}

export interface TaskEvent {
  taskId: string
  name: string
  progress: number
  state: string
  detail: string
  scope: string
}

export interface CommandRecord {
  id: string
  name: string
  workspaceId?: string
  capability: string
  provider?: string
  target?: string
  status: 'queued' | 'running' | 'completed' | 'failed' | 'cancelled' | 'rolled_back'
  taskId?: string
  progress?: number
  detail?: string
  error?: string
  result?: unknown
  parameters?: Record<string, unknown>
  metadata?: Record<string, string>
  createdAt: string
  startedAt?: string
  completedAt?: string
}

export interface DashboardDataV2 {
  metrics: DashboardMetric[]
  sections: DashboardSectionCard[]
  activity: ActivityItem[]
  tasks: TaskItem[]
  recommendations: RecommendationItem[]
  knowledge: KnowledgeArticle[]
  server: string
  notifications: number
}

export type ProjectStatus = 'healthy' | 'warning' | 'maintenance' | 'blocked'

export interface ProjectListItem {
  name: string
  slug: string
  description: string
  owner: string
  repository: string
  environment: string
  status: ProjectStatus
  updatedAt: string
  progress: number
}

export interface ProjectDeployment {
  id: string
  version: string
  status: 'succeeded' | 'running' | 'queued' | 'failed'
  environment: string
  deployedAt: string
  note: string
}

export interface ProjectDomain {
  host: string
  ssl: 'Valid' | 'Renewing' | 'Pending'
  target: string
}

export interface ProjectLogEntry {
  id: string
  title: string
  detail: string
  time: string
  tone: Tone
}

export interface ProjectEnvironmentItem {
  key: string
  value: string
  visibility: 'public' | 'secret'
}

export interface ProjectDatabaseItem {
  name: string
  engine: string
  status: 'ready' | 'syncing' | 'archived'
}

export interface ProjectDetailData {
  project: ProjectListItem
  repository: {
    url: string
    branch: string
    buildCommand: string
    startCommand: string
  }
  environment: ProjectEnvironmentItem[]
  deployments: ProjectDeployment[]
  domains: ProjectDomain[]
  logs: ProjectLogEntry[]
  database: ProjectDatabaseItem
  overview: DashboardSectionCard[]
  activity: ActivityItem[]
  tasks: TaskItem[]
  recommendations: RecommendationItem[]
  knowledge: KnowledgeArticle[]
}

export interface ProjectFormData {
  name: string
  slug: string
  description: string
  repository: string
  branch: string
  environment: string
  owner: string
  deployTarget: string
  healthCheck: string
  autoDeploy: boolean
  domains: string
}

export interface ProjectFormOptions {
  title: string
  subtitle: string
  submitLabel: string
}

export interface ProjectSettingsItem {
  label: string
  value: string
  detail: string
}

export interface SettingsSection {
  title: string
  description: string
  items: ProjectSettingsItem[]
}

export interface SettingsData {
  sections: SettingsSection[]
  preferences: Array<{
    label: string
    value: string
    detail: string
  }>
  security: Array<{
    label: string
    value: string
    detail: string
  }>
  aiProviders?: Array<{
    label: string
    value: string
    detail: string
  }>
  mcpIntegrations?: Array<{
    label: string
    value: string
    detail: string
  }>
}

export interface DevCenterSectionData {
  key: string
  title: string
  subtitle: string
  description: string
  bullets: string[]
  examples: string[]
  checklist: ChecklistItem[]
  knowledge: KnowledgeArticle[]
  codeSample: string
  metadata?: Record<string, string>
}

// Sprint 4 — Workspace Runtime types

export interface InfraToolInfo {
  name: string
  installed: boolean
  version: string
  healthy: boolean
  configured: boolean
}

export interface InfrastructureInfo {
  tools: InfraToolInfo[]
  os: string
  arch: string
  hostname: string
}

export interface WorkspaceServiceInfo {
  name: string
  status: string
  port: string
  pid: string
}

export interface DeploymentEntry {
  id: string
  branch: string
  commit: string
  status: string
  timestamp: string
  environment: string
}

export interface DeploymentInfo {
  entries: DeploymentEntry[]
  current: string
}

export interface WorkspaceDomainInfo {
  host: string
  ssl: string
  target: string
}

export interface WorkspaceLogEntry {
  timestamp: string
  level: string
  message: string
  source: string
}

export interface WorkspaceGitInfo {
  branch: string
  commit: string
  clean: boolean
  ahead: number
  behind: number
  recentCommits: string[]
  changedFiles: string[]
}

export interface WorkspaceEnvironmentInfo {
  files: string[]
  keys: string[]
  secrets: string[]
  preview: Record<string, string>
}

export interface WorkspaceFileInfo {
  path: string
  kind: string
  size: number
  modifiedAt: string
}

export interface WorkspaceOverview {
  workspace: {
    id: string
    name: string
    root: string
    kind: string
    framework: string
    languages: string[]
    runtime: string
    packageManager: string
    generatedAt: string
  }
  project: {
    id: string
    name: string
    root: string
    framework: string
    languages: string[]
    runtime: string
    packageManager: string
    repository: string
  }
  health: {
    score: number
    build: string
    tests: string
    lint: string
  }
  git: WorkspaceGitInfo
  plugins: {
    detected: string[]
    capabilities: string[]
  }
}

export interface WorkspaceFilesResponse {
  files: WorkspaceFileInfo[]
  total: number
}

export interface AIContextInfo {
  workspace: WorkspaceOverview['workspace']
  git: WorkspaceGitInfo
  architecture: {
    kind: string
    framework: string
    languages: string[]
    runtime: string
    packageManager: string
    entryPoints: string[]
    notes: string[]
  }
  tasks: { suggested: string[] }
  dependencies: {
    frontend: string[]
    backend: string[]
    database: string[]
    tools: string[]
  }
  routes: {
    next: string[]
    go: string[]
    api: string[]
  }
  database: {
    kind: string
    files: string[]
    migrations: string[]
    environment: string[]
  }
  generatedAt: string
}

export interface MCPProvider {
  name: string
  enabled: boolean
  healthy: boolean
  endpoint: string
  protocol: string
}

export interface MCPInfo {
  providers: MCPProvider[]
}

export interface WorkspaceDoctorData {
  health: WorkspaceOverview['health']
  project: WorkspaceOverview['project']
  git: WorkspaceGitInfo
  plugins: WorkspaceOverview['plugins']
  checks: Array<{ label: string; status: string; detail: string }>
  recommendations: string[]
}

export interface WorkspaceSettingsData {
  workspace: WorkspaceOverview['workspace']
  index: {
    version: number
    generatedAt: string
    workspace: string
    files: Record<string, string>
  }
  cache: {
    updatedAt: string
    fingerprint: string
    changed: string[]
  }
}

export type WorkspaceSection =
  | 'overview'
  | 'files'
  | 'repository'
  | 'environment'
  | 'infrastructure'
  | 'services'
  | 'deployments'
  | 'database'
  | 'domains'
  | 'logs'
  | 'ai'
  | 'knowledge'
  | 'doctor'
  | 'settings'
  | 'mcp'
