import { request, requestOrFallback } from '../api/client'
import { projectDetailFallback, projectFallbacks } from './fallbacks'
import type { ProjectDetailData, ProjectFormData, ProjectFormOptions, ProjectListItem } from '../types'

export async function getProjects(): Promise<ProjectListItem[]> {
  return requestOrFallback<ProjectListItem[]>('/api/projects', projectFallbacks())
}

export async function getProjectDetail(slug: string): Promise<ProjectDetailData> {
  return requestOrFallback<ProjectDetailData>(`/api/projects/${encodeURIComponent(slug)}`, projectDetailFallback(slug))
}

export async function getProjectFormData(slug?: string): Promise<ProjectFormData> {
  const path = slug ? `/api/projects/${encodeURIComponent(slug)}/form` : '/api/projects/form'
  return request<ProjectFormData>(path)
}

export async function getProjectFormOptions(mode: 'create' | 'edit'): Promise<ProjectFormOptions> {
  return {
    title: mode === 'edit' ? 'Edit project' : 'Create project',
    subtitle: mode === 'edit' ? 'Update the project metadata and deployment settings.' : 'Create a new project and prepare its deploy surface.',
    submitLabel: mode === 'edit' ? 'Save project' : 'Create project',
  }
}

export async function saveProject(data: ProjectFormData, mode: 'create' | 'edit'): Promise<ProjectListItem> {
  const method = mode === 'edit' ? 'PATCH' : 'POST'
  const path = mode === 'edit' ? `/api/projects/${encodeURIComponent(data.slug)}` : '/api/projects'
  return request<ProjectListItem>(path, {
    method,
    body: data,
  })
}
