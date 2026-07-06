import { request, submitCommand } from '../api/client'
import type {
  WorkspaceOverview,
  WorkspaceFilesResponse,
  WorkspaceGitInfo,
  WorkspaceEnvironmentInfo,
  InfrastructureInfo,
  WorkspaceProviderInfo,
  DeploymentInfo,
  WorkspaceDomainInfo,
  WorkspaceLogEntry,
  AIContextInfo,
  MCPInfo,
  WorkspaceDoctorData,
  WorkspaceSettingsData,
} from '../types'

export async function getWorkspaceOverview(
  id: string,
): Promise<WorkspaceOverview> {
  return request(`/api/workspaces/${id}/overview`)
}

export async function getWorkspaceFiles(
  id: string,
  query = '',
): Promise<WorkspaceFilesResponse> {
  const suffix = query ? `?q=${encodeURIComponent(query)}` : ''
  return request(`/api/workspaces/${id}/files${suffix}`)
}

export async function getWorkspaceGit(
  id: string,
): Promise<WorkspaceGitInfo> {
  return request(`/api/workspaces/${id}/git`)
}

export async function getWorkspaceEnvironment(
  id: string,
): Promise<WorkspaceEnvironmentInfo> {
  return request(`/api/workspaces/${id}/environment`)
}

export async function getWorkspaceInfrastructure(
  id: string,
): Promise<InfrastructureInfo> {
  return request(`/api/workspaces/${id}/infrastructure`)
}

export async function getWorkspaceServices(
  id: string,
): Promise<WorkspaceProviderInfo[]> {
  return request(`/api/workspaces/${id}/services`)
}

export async function getWorkspaceDeployments(
  id: string,
): Promise<DeploymentInfo> {
  return request(`/api/workspaces/${id}/deployments`)
}

export async function getWorkspaceDomains(
  id: string,
): Promise<WorkspaceDomainInfo[]> {
  return request(`/api/workspaces/${id}/domains`)
}

export async function getWorkspaceLogs(
  id: string,
): Promise<WorkspaceLogEntry[]> {
  return request(`/api/workspaces/${id}/logs`)
}

export async function getWorkspaceAI(
  id: string,
): Promise<AIContextInfo> {
  return request(`/api/workspaces/${id}/ai`)
}

export async function getWorkspaceKnowledge(id: string): Promise<any> {
  return request(`/api/workspaces/${id}/knowledge`)
}

export async function getWorkspaceMCP(
  id: string,
): Promise<MCPInfo> {
  return request(`/api/workspaces/${id}/mcp`)
}

export async function getWorkspaceDoctor(
  id: string,
): Promise<WorkspaceDoctorData> {
  return request(`/api/workspaces/${id}/doctor`)
}

export async function getWorkspaceSettings(
  id: string,
): Promise<WorkspaceSettingsData> {
  return request(`/api/workspaces/${id}/settings`)
}

export async function triggerWorkspaceReindex(): Promise<{ status: string }> {
  return submitCommand<{ status: string }>({
    capability: 'workspace.index',
    workspaceId: 'devserver',
    name: 'Reindex workspace',
  })
}
