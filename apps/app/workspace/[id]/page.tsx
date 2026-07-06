import {
  getWorkspaceOverview,
  getWorkspaceFiles,
  getWorkspaceGit,
  getWorkspaceEnvironment,
  getWorkspaceInfrastructure,
  getWorkspaceServices,
  getWorkspaceDeployments,
  getWorkspaceDomains,
  getWorkspaceLogs,
  getWorkspaceAI,
  getWorkspaceMCP,
  getWorkspaceDoctor,
  getWorkspaceSettings,
  getWorkspaceKnowledge,
} from '@/lib/services/workspace'

import { AppShell } from '@/components/app-shell'
import { WorkspacePageClient } from '@/components/workspace/workspace-page-client'

interface PageProps {
  params: Promise<{ id: string }>
  searchParams: Promise<{ section?: string }>
}

export default async function WorkspacePage({ params, searchParams }: PageProps) {
  const { id } = await params
  const sp = await searchParams
  const section = sp.section || 'overview'

  const [overview, filesRes, git, env, infra, services, deployments, domains, logs, ai, knowledge, mcp, doctor, settings] = await Promise.all([
    getWorkspaceOverview(id),
    getWorkspaceFiles(id),
    getWorkspaceGit(id),
    getWorkspaceEnvironment(id),
    getWorkspaceInfrastructure(id),
    getWorkspaceServices(id),
    getWorkspaceDeployments(id),
    getWorkspaceDomains(id),
    getWorkspaceLogs(id),
    getWorkspaceAI(id),
    getWorkspaceKnowledge(id),
    getWorkspaceMCP(id),
    getWorkspaceDoctor(id),
    getWorkspaceSettings(id),
  ])

  return (
    <AppShell server={overview.workspace?.name || 'DevServer'} notifications={0}>
      <WorkspacePageClient
        id={id}
        section={section}
        overview={overview}
        files={filesRes}
        git={git}
        environment={env}
        infrastructure={infra}
        services={services}
        deployments={deployments}
        database={{ kind: overview.project?.framework || '', files: [], migrations: [], environment: [] }}
        domains={domains}
        logs={logs}
        ai={ai}
        knowledge={knowledge}
        mcp={mcp}
        doctor={doctor}
        settings={settings}
      />
    </AppShell>
  )
}
