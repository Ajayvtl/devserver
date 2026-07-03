import { submitCommand } from '../api/client'
import { projectDetailFallback, projectFallbacks } from './fallbacks'
import type { CommandRecord, ProjectDetailData, ProjectFormData, ProjectFormOptions, ProjectListItem } from '../types'

export async function getProjects(): Promise<ProjectListItem[]> {
  return projectFallbacks()
}

export async function getProjectDetail(slug: string): Promise<ProjectDetailData> {
  return projectDetailFallback(slug)
}

export async function getProjectFormData(slug?: string): Promise<ProjectFormData> {
  const project = projectDetailFallback(slug ?? 'aurora').project
  return {
    name: project.name,
    slug: project.slug,
    description: project.description,
    repository: project.repository,
    branch: project.environment === 'production' ? 'main' : 'develop',
    environment: project.environment,
    owner: project.owner,
    deployTarget: 'Primary fleet',
    healthCheck: '/health',
    autoDeploy: project.status === 'healthy',
    domains: `${project.slug}.devserver.local, api.${project.slug}.devserver.local`,
  }
}

export async function getProjectFormOptions(mode: 'create' | 'edit'): Promise<ProjectFormOptions> {
  return {
    title: mode === 'edit' ? 'Edit project' : 'Create project',
    subtitle: mode === 'edit' ? 'Update the project metadata and deployment settings.' : 'Create a new project and prepare its deploy surface.',
    submitLabel: mode === 'edit' ? 'Save project' : 'Create project',
  }
}

export async function saveProject(data: ProjectFormData, mode: 'create' | 'edit'): Promise<ProjectListItem> {
  const response = await submitCommand<CommandRecord>({
    capability: 'project.save',
    parameters: data as unknown as Record<string, unknown>,
    name: mode === 'edit' ? 'Save project' : 'Create project',
  })

  const result = response.result as ProjectListItem | undefined
  if (!result) {
    throw new Error('Project save did not return a project record.')
  }

  return result
}
