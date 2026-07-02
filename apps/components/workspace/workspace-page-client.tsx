'use client'

import type {
  WorkspaceOverview, WorkspaceFilesResponse, WorkspaceGitInfo,
  WorkspaceEnvironmentInfo, InfrastructureInfo, WorkspaceServiceInfo,
  DeploymentInfo, WorkspaceDomainInfo, WorkspaceLogEntry,
  AIContextInfo, MCPInfo, WorkspaceDoctorData, WorkspaceSettingsData,
  WorkspaceSection,
} from '@/lib/types'

import { WorkspaceLayout } from './workspace-layout'
import { OverviewSection } from './overview-section'
import { FilesSection } from './files-section'
import { RepositorySection } from './repository-section'
import { EnvironmentSection } from './environment-section'
import { InfrastructureSection } from './infrastructure-section'
import { ServicesSection } from './services-section'
import { DeploymentsSection } from './deployments-section'
import {
  DatabaseSection, DomainsSection, LogsSection, AISection,
  KnowledgeSection, DoctorSection, SettingsSection, MCPSection,
} from './remaining-sections'

interface Props {
  id: string
  section: string
  overview: WorkspaceOverview
  files: WorkspaceFilesResponse
  git: WorkspaceGitInfo
  environment: WorkspaceEnvironmentInfo
  infrastructure: InfrastructureInfo
  services: WorkspaceServiceInfo[]
  deployments: DeploymentInfo
  database: { kind: string; files: string[]; migrations: string[]; environment: string[] }
  domains: WorkspaceDomainInfo[]
  logs: WorkspaceLogEntry[]
  ai: AIContextInfo
  knowledge: { files: string[]; topics: string[] }
  mcp: MCPInfo
  doctor: WorkspaceDoctorData
  settings: WorkspaceSettingsData
}

export function WorkspacePageClient(props: Props) {
  const active = (props.section || 'overview') as WorkspaceSection

  function renderSection() {
    switch (active) {
      case 'overview':
        return <OverviewSection data={props.overview} />
      case 'files':
        return <FilesSection files={props.files.files} total={props.files.total} />
      case 'repository':
        return <RepositorySection data={props.git} />
      case 'environment':
        return <EnvironmentSection data={props.environment} />
      case 'infrastructure':
        return <InfrastructureSection data={props.infrastructure} />
      case 'services':
        return <ServicesSection data={props.services} />
      case 'deployments':
        return <DeploymentsSection data={props.deployments} />
      case 'database':
        return <DatabaseSection data={props.database} />
      case 'domains':
        return <DomainsSection data={props.domains} />
      case 'logs':
        return <LogsSection data={props.logs} />
      case 'ai':
        return <AISection data={props.ai} />
      case 'knowledge':
        return <KnowledgeSection data={props.knowledge} />
      case 'doctor':
        return <DoctorSection data={props.doctor} />
      case 'settings':
        return <SettingsSection data={props.settings} />
      case 'mcp':
        return <MCPSection data={props.mcp} />
      default:
        return <OverviewSection data={props.overview} />
    }
  }

  return (
    <WorkspaceLayout id={props.id} active={active}>
      {renderSection()}
    </WorkspaceLayout>
  )
}
