import { AppShell } from '@/components/app-shell'
import { ProjectsPage } from '@/components/projects/projects-page'
import { getProjects } from '@/lib/services/projects'

export default async function ProjectsIndexPage() {
  const projects = await getProjects()

  return (
    <AppShell server="production-east-1" notifications={4}>
      <ProjectsPage projects={projects} />
    </AppShell>
  )
}
