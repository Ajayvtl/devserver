'use client'

import type {
  WorkspaceOverview, WorkspaceFilesResponse, WorkspaceGitInfo,
  WorkspaceEnvironmentInfo, InfrastructureInfo, WorkspaceProviderInfo,
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
import { TasksSection } from './tasks-section'
import { KnowledgeSection } from './knowledge-section'
import {
  DatabaseSection, DomainsSection, LogsSection, AISection,
  DoctorSection, SettingsSection, MCPSection,
} from './remaining-sections'

interface Props {
  id: string
  section: string
  overview: WorkspaceOverview
  files: WorkspaceFilesResponse
  git: WorkspaceGitInfo
  environment: WorkspaceEnvironmentInfo
  infrastructure: InfrastructureInfo
  services: WorkspaceProviderInfo[]
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

import { useState } from 'react'
import type { Document, WorkspaceSession } from '@/lib/types'
import { WorkspaceEditor } from './workspace-editor'

export function WorkspacePageClient(props: Props) {
  const active = (props.section || 'overview') as WorkspaceSection
  
  const [session, setSession] = useState<WorkspaceSession>({
    workspaceId: props.id,
    openTabs: [],
    activeTab: null,
    recent: [],
    favorites: [],
    expandedFolders: [],
    selectedFiles: [],
    searchHistory: []
  })

  const openDocument = (path: string, line: number = 0) => {
    setSession(prev => {
      const existing = prev.openTabs.find(d => d.path === path)
      let nextTabs = prev.openTabs
      if (existing) {
        nextTabs = prev.openTabs.map(d => d.path === path ? { ...d, cursor: { line, column: 0 } } : d)
      } else {
        nextTabs = [...prev.openTabs, { workspaceId: props.id, path, language: 'text', version: 1, dirty: false, cursor: { line, column: 0 } }]
      }
      return {
        ...prev,
        openTabs: nextTabs,
        activeTab: path,
        recent: [path, ...prev.recent.filter(p => p !== path)].slice(0, 10)
      }
    })
  }

  const closeDocument = (path: string) => {
    setSession(prev => {
      const nextTabs = prev.openTabs.filter(d => d.path !== path)
      let nextActive = prev.activeTab
      if (prev.activeTab === path) {
        nextActive = nextTabs.length > 0 ? nextTabs[nextTabs.length - 1].path : null
      }
      return { ...prev, openTabs: nextTabs, activeTab: nextActive }
    })
  }

  const updateSession = (updates: Partial<WorkspaceSession>) => {
    setSession(prev => ({ ...prev, ...updates }))
  }

  function renderSection() {
    switch (active) {
      case 'overview':
        return <OverviewSection data={props.overview} />
      case 'files':
        return <FilesSection 
          files={props.files.files} 
          total={props.files.total} 
          onOpenFile={(path) => openDocument(path, 0)} 
          session={session} 
          onSessionChange={updateSession} 
        />
      case 'repository':
        return <RepositorySection data={props.git} />
      case 'environment':
        return <EnvironmentSection data={props.environment} />
      case 'infrastructure':
        return <InfrastructureSection data={props.infrastructure} />
      case 'services':
        return <ServicesSection data={props.services} />
      case 'tasks':
        return <TasksSection workspaceId={props.id} />
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
        return <KnowledgeSection data={props.knowledge} onJumpToDefinition={(file, line) => openDocument(file, line)} />
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
      {session.openTabs.length > 0 ? (
        <div style={{ display: 'grid', gridTemplateColumns: 'minmax(350px, 450px) 1fr', gap: '16px', height: '100%', padding: '16px' }}>
          <div style={{ overflowY: 'auto', background: 'rgba(5, 10, 19, 0.4)', borderRadius: '12px' }}>
            {renderSection()}
          </div>
          <WorkspaceEditor 
            workspaceId={props.id}
            documents={session.openTabs}
            activeDocPath={session.activeTab}
            onTabClick={(path) => updateSession({ activeTab: path })}
            onCloseTab={closeDocument}
          />
        </div>
      ) : (
        renderSection()
      )}
    </WorkspaceLayout>
  )
}
