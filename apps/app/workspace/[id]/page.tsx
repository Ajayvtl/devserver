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

  const [overview, filesRes, git, env, infra, services, deployments, domains, logs, ai, mcp, doctor, settings] = await Promise.all([
    getWorkspaceOverview(),
    getWorkspaceFiles(),
    getWorkspaceGit(),
    getWorkspaceEnvironment(),
    getWorkspaceInfrastructure(),
    getWorkspaceServices(),
    getWorkspaceDeployments(),
    getWorkspaceDomains(),
    getWorkspaceLogs(),
    getWorkspaceAI(),
    getWorkspaceMCP(),
    getWorkspaceDoctor(),
    getWorkspaceSettings(),
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
        knowledge={{ files: [], topics: [] }}
        mcp={mcp}
        doctor={doctor}
        settings={settings}
      />
    </AppShell>
  )
}
