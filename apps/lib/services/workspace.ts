import { submitCommand } from '../api/client'
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
} from '../types'

export async function getWorkspaceOverview(): Promise<WorkspaceOverview> {
  return {
    workspace: { id: 'devserver', name: 'Devserver', root: '.', kind: 'Go', framework: 'Go', languages: ['Go', 'TypeScript'], runtime: 'go', packageManager: 'npm', generatedAt: new Date().toISOString() },
    project: { id: 'devserver', name: 'Devserver', root: '.', framework: 'Go', languages: ['Go', 'TypeScript'], runtime: 'go', packageManager: 'npm', repository: 'git' },
    health: { score: 60, build: 'unknown', tests: 'unknown', lint: 'unknown' },
    git: { branch: 'main', commit: '', clean: true, ahead: 0, behind: 0, recentCommits: [], changedFiles: [] },
    plugins: { detected: [], capabilities: [] },
  }
}

export async function getWorkspaceFiles(query = ''): Promise<WorkspaceFilesResponse> {
  return { files: [], total: 0 }
}

export async function getWorkspaceGit(): Promise<WorkspaceGitInfo> {
  return {
    branch: '', commit: '', clean: true, ahead: 0, behind: 0, recentCommits: [], changedFiles: [],
  }
}

export async function getWorkspaceEnvironment(): Promise<WorkspaceEnvironmentInfo> {
  return {
    files: [], keys: [], secrets: [], preview: {},
  }
}

export async function getWorkspaceInfrastructure(): Promise<InfrastructureInfo> {
  return {
    tools: [], os: '', arch: '', hostname: '',
  }
}

export async function getWorkspaceServices(): Promise<WorkspaceServiceInfo[]> {
  return []
}

export async function getWorkspaceDeployments(): Promise<DeploymentInfo> {
  return { entries: [], current: '' }
}

export async function getWorkspaceDomains(): Promise<WorkspaceDomainInfo[]> {
  return []
}

export async function getWorkspaceLogs(): Promise<WorkspaceLogEntry[]> {
  return []
}

export async function getWorkspaceAI(): Promise<AIContextInfo> {
  return {
    workspace: { id: '', name: '', root: '', kind: '', framework: '', languages: [], runtime: '', packageManager: '', generatedAt: '' },
    git: { branch: '', commit: '', clean: true, ahead: 0, behind: 0, recentCommits: [], changedFiles: [] },
    architecture: { kind: '', framework: '', languages: [], runtime: '', packageManager: '', entryPoints: [], notes: [] },
    tasks: { suggested: [] },
    dependencies: { frontend: [], backend: [], database: [], tools: [] },
    routes: { next: [], go: [], api: [] },
    database: { kind: '', files: [], migrations: [], environment: [] },
    generatedAt: '',
  }
}

export async function getWorkspaceMCP(): Promise<MCPInfo> {
  return { providers: [] }
}

export async function getWorkspaceDoctor(): Promise<WorkspaceDoctorData> {
  return {
    health: { score: 0, build: 'unknown', tests: 'unknown', lint: 'unknown' },
    project: { id: '', name: '', root: '', framework: '', languages: [], runtime: '', packageManager: '', repository: '' },
    git: { branch: '', commit: '', clean: true, ahead: 0, behind: 0, recentCommits: [], changedFiles: [] },
    plugins: { detected: [], capabilities: [] },
    checks: [],
    recommendations: [],
  }
}

export async function getWorkspaceSettings(): Promise<WorkspaceSettingsData> {
  return {
    workspace: { id: '', name: '', root: '', kind: '', framework: '', languages: [], runtime: '', packageManager: '', generatedAt: '' },
    index: { version: 0, generatedAt: '', workspace: '', files: {} },
    cache: { updatedAt: '', fingerprint: '', changed: [] },
  }
}

export async function triggerWorkspaceReindex(): Promise<{ status: string }> {
  return submitCommand<{ status: string }>({
    capability: 'workspace.index',
    workspaceId: 'devserver',
    name: 'Reindex workspace',
  })
}
