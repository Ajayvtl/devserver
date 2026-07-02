import { requestOrFallback } from '../api/client'
import type {
  WorkspaceOverview,
  WorkspaceFilesResponse,
  WorkspaceGitInfo,
  WorkspaceEnvironmentInfo,
  InfrastructureInfo,
  WorkspaceServiceInfo,
  DeploymentInfo,
  WorkspaceDomainInfo,
  WorkspaceLogEntry,
  AIContextInfo,
  MCPInfo,
  WorkspaceDoctorData,
  WorkspaceSettingsData,
  WorkspaceSection,
} from '../types'

const WORKSPACE_ID = 'devserver'

function wsPath(section: WorkspaceSection | string, query = '') {
  const base = `/api/workspaces/${WORKSPACE_ID}/${section}`
  return query ? `${base}?${query}` : base
}

export async function getWorkspaceOverview(): Promise<WorkspaceOverview> {
  return requestOrFallback<WorkspaceOverview>(wsPath('overview'), {
    workspace: { id: 'devserver', name: 'Devserver', root: '.', kind: 'Go', framework: 'Go', languages: ['Go', 'TypeScript'], runtime: 'go', packageManager: 'npm', generatedAt: new Date().toISOString() },
    project: { id: 'devserver', name: 'Devserver', root: '.', framework: 'Go', languages: ['Go', 'TypeScript'], runtime: 'go', packageManager: 'npm', repository: 'git' },
    health: { score: 60, build: 'unknown', tests: 'unknown', lint: 'unknown' },
    git: { branch: 'main', commit: '', clean: true, ahead: 0, behind: 0, recentCommits: [], changedFiles: [] },
    plugins: { detected: [], capabilities: [] },
  })
}

export async function getWorkspaceFiles(query = ''): Promise<WorkspaceFilesResponse> {
  return requestOrFallback<WorkspaceFilesResponse>(wsPath('files', query ? `q=${encodeURIComponent(query)}` : ''), { files: [], total: 0 })
}

export async function getWorkspaceGit(): Promise<WorkspaceGitInfo> {
  return requestOrFallback<WorkspaceGitInfo>(wsPath('repository'), {
    branch: '', commit: '', clean: true, ahead: 0, behind: 0, recentCommits: [], changedFiles: [],
  })
}

export async function getWorkspaceEnvironment(): Promise<WorkspaceEnvironmentInfo> {
  return requestOrFallback<WorkspaceEnvironmentInfo>(wsPath('environment'), {
    files: [], keys: [], secrets: [], preview: {},
  })
}

export async function getWorkspaceInfrastructure(): Promise<InfrastructureInfo> {
  return requestOrFallback<InfrastructureInfo>(wsPath('infrastructure'), {
    tools: [], os: '', arch: '', hostname: '',
  })
}

export async function getWorkspaceServices(): Promise<WorkspaceServiceInfo[]> {
  return requestOrFallback<WorkspaceServiceInfo[]>(wsPath('services'), [])
}

export async function getWorkspaceDeployments(): Promise<DeploymentInfo> {
  return requestOrFallback<DeploymentInfo>(wsPath('deployments'), { entries: [], current: '' })
}

export async function getWorkspaceDomains(): Promise<WorkspaceDomainInfo[]> {
  return requestOrFallback<WorkspaceDomainInfo[]>(wsPath('domains'), [])
}

export async function getWorkspaceLogs(): Promise<WorkspaceLogEntry[]> {
  return requestOrFallback<WorkspaceLogEntry[]>(wsPath('logs'), [])
}

export async function getWorkspaceAI(): Promise<AIContextInfo> {
  return requestOrFallback<AIContextInfo>(wsPath('ai'), {
    workspace: { id: '', name: '', root: '', kind: '', framework: '', languages: [], runtime: '', packageManager: '', generatedAt: '' },
    git: { branch: '', commit: '', clean: true, ahead: 0, behind: 0, recentCommits: [], changedFiles: [] },
    architecture: { kind: '', framework: '', languages: [], runtime: '', packageManager: '', entryPoints: [], notes: [] },
    tasks: { suggested: [] },
    dependencies: { frontend: [], backend: [], database: [], tools: [] },
    routes: { next: [], go: [], api: [] },
    database: { kind: '', files: [], migrations: [], environment: [] },
    generatedAt: '',
  })
}

export async function getWorkspaceMCP(): Promise<MCPInfo> {
  return requestOrFallback<MCPInfo>(wsPath('mcp'), { providers: [] })
}

export async function getWorkspaceDoctor(): Promise<WorkspaceDoctorData> {
  return requestOrFallback<WorkspaceDoctorData>(wsPath('doctor'), {
    health: { score: 0, build: 'unknown', tests: 'unknown', lint: 'unknown' },
    project: { id: '', name: '', root: '', framework: '', languages: [], runtime: '', packageManager: '', repository: '' },
    git: { branch: '', commit: '', clean: true, ahead: 0, behind: 0, recentCommits: [], changedFiles: [] },
    plugins: { detected: [], capabilities: [] },
    checks: [],
    recommendations: [],
  })
}

export async function getWorkspaceSettings(): Promise<WorkspaceSettingsData> {
  return requestOrFallback<WorkspaceSettingsData>(wsPath('settings'), {
    workspace: { id: '', name: '', root: '', kind: '', framework: '', languages: [], runtime: '', packageManager: '', generatedAt: '' },
    index: { version: 0, generatedAt: '', workspace: '', files: {} },
    cache: { updatedAt: '', fingerprint: '', changed: [] },
  })
}

export async function triggerWorkspaceReindex(): Promise<{ status: string }> {
  return requestOrFallback<{ status: string }>(`/api/workspaces/${WORKSPACE_ID}/index`, { status: 'offline' }, { method: 'POST' })
}
