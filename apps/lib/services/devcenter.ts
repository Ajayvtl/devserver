import { devcenterFallback } from './fallbacks'
import type { DevCenterSectionData } from '../types'

export async function getDevCenterSection(slug: string): Promise<DevCenterSectionData> {
  return devcenterFallback(slug)
}

export async function listDevCenterSections(): Promise<Array<{ key: string; title: string; subtitle: string }>> {
  return [
    { key: 'architecture', title: 'Architecture', subtitle: 'Understand how DevServer is structured end to end.' },
    { key: 'tasks', title: 'Tasks', subtitle: 'Track installation, deployment, and maintenance workflows.' },
    { key: 'knowledge', title: 'Knowledge', subtitle: 'Teach while you build with documentation and best practices.' },
    { key: 'dependencies', title: 'Dependencies', subtitle: 'Map module and service dependencies explicitly.' },
    { key: 'api', title: 'API', subtitle: 'Plan the eventual API surface without coupling the UI to it yet.' },
    { key: 'database', title: 'Database', subtitle: 'Understand the platform data model and how it supports the UI.' },
    { key: 'project-doctor', title: 'Project Doctor', subtitle: 'Inspect a project and surface risks before deployment.' },
  ]
}
