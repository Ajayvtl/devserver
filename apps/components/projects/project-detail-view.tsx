'use client'

import { useMemo, useState } from 'react'

import type { ProjectDetailData } from '@/lib/types'

import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { CodeBlock } from '../ui/code-block'
import { DataTable } from '../ui/table'
import { Dialog } from '../ui/dialog'
import { EmptyState } from '../ui/empty-state'
import { LearnCard } from '../ui/learn-card'
import { Progress } from '../ui/progress'
import { useToast } from '../toast'

interface Props {
  data: ProjectDetailData
}

const workspaceSections = [
  { id: 'overview', label: 'Overview' },
  { id: 'repository', label: 'Repository' },
  { id: 'environment', label: 'Environment' },
  { id: 'infrastructure', label: 'Infrastructure' },
  { id: 'files', label: 'Files' },
  { id: 'database', label: 'Database' },
  { id: 'domains', label: 'Domains' },
  { id: 'deployment', label: 'Deployment' },
  { id: 'monitoring', label: 'Monitoring' },
  { id: 'logs', label: 'Logs' },
  { id: 'ai', label: 'AI' },
  { id: 'knowledge', label: 'Knowledge' },
  { id: 'doctor', label: 'Doctor' },
  { id: 'settings', label: 'Settings' },
] as const

export function ProjectDetailView({ data }: Props) {
  const [dialogOpen, setDialogOpen] = useState(false)
  const { push } = useToast()

  const fileTree = useMemo(
    () => [
      { path: 'src/app/page.tsx', kind: 'file' },
      { path: 'src/app/layout.tsx', kind: 'file' },
      { path: 'src/lib/api.ts', kind: 'file' },
      { path: 'src/components/', kind: 'folder' },
      { path: 'src/components/dashboard/', kind: 'folder' },
      { path: 'package.json', kind: 'file' },
      { path: 'README.md', kind: 'file' },
    ],
    [],
  )

  return (
    <div className="workspace-layout">
      <aside className="workspace-layout__nav">
        <Card className="workspace-nav-card">
          <div className="brand-lockup brand-lockup--compact">
            <div className="brand-lockup__mark">D</div>
            <div>
              <div className="brand-lockup__title">{data.project.name}</div>
              <div className="brand-lockup__subtitle">Workspace</div>
            </div>
          </div>
          <div className="workspace-nav__meta">
            <Badge tone="accent">{data.project.environment}</Badge>
            <Badge tone={data.project.status === 'healthy' ? 'success' : 'warning'}>{data.project.status}</Badge>
          </div>
        </Card>

        <nav className="workspace-nav" aria-label="Workspace navigation">
          {workspaceSections.map((section) => (
            <a key={section.id} className="workspace-nav__item" href={`#${section.id}`}>
              <span>{section.label}</span>
            </a>
          ))}
        </nav>
      </aside>

      <div className="workspace-layout__content">
        <section className="workspace-hero" id="overview">
          <div>
            <div className="card__eyebrow">Workspace</div>
            <h1 className="page-title">{data.project.name}</h1>
            <p className="page-subtitle">{data.project.description}</p>
            <div className="form-badges">
              <Badge tone="accent">{data.project.environment}</Badge>
              <Badge tone={data.project.status === 'healthy' ? 'success' : 'warning'}>{data.project.status}</Badge>
              <Badge tone="info">{data.project.owner}</Badge>
            </div>
          </div>
          <div className="project-hero__actions">
            <Button href={`/projects/${data.project.slug}/edit`} variant="secondary">
              Edit workspace
            </Button>
            <Button href="#deployment" variant="primary">
              New deployment
            </Button>
          </div>
        </section>

        <section className="workspace-grid">
          <Card>
            <div className="card__eyebrow">Health</div>
            <h3>{data.project.progress}% ready</h3>
            <Progress value={data.project.progress} tone={data.project.status === 'healthy' ? 'success' : 'warning'} />
          </Card>
          <Card>
            <div className="card__eyebrow">Repository</div>
            <h3>{data.repository.url}</h3>
            <p>Branch {data.repository.branch}</p>
          </Card>
        </section>

        <section id="repository" className="workspace-section">
          <Card>
            <SectionTitle title="Repository" description="Repository source, build, and runtime commands." />
            <div className="detail-grid">
              <Item label="URL" value={data.repository.url} />
              <Item label="Branch" value={data.repository.branch} />
              <Item label="Build command" value={data.repository.buildCommand} />
              <Item label="Start command" value={data.repository.startCommand} />
            </div>
            <CodeBlock title="Repository config">{`repo: ${data.repository.url}
branch: ${data.repository.branch}
build: ${data.repository.buildCommand}
start: ${data.repository.startCommand}`}</CodeBlock>
          </Card>
        </section>

        <section id="environment" className="workspace-section">
          <Card>
            <SectionTitle title="Environment" description="Variables and visibility rules." />
            <DataTable
              columns={[{ header: 'Key' }, { header: 'Value' }, { header: 'Visibility' }]}
              rows={data.environment.map((item) => [
                item.key,
                item.visibility === 'secret' ? '••••••••' : item.value,
                <Badge key={item.key} tone={item.visibility === 'secret' ? 'warning' : 'info'}>
                  {item.visibility}
                </Badge>,
              ])}
            />
          </Card>
        </section>

        <section id="infrastructure" className="workspace-section">
          <Card>
            <SectionTitle title="Infrastructure" description="Installed and healthy workspace dependencies." />
            <div className="metric-grid metric-grid--three">
              <Metric label="Node" value="Installed" tone="success" />
              <Metric label="Nginx" value="Healthy" tone="success" />
              <Metric label="Postgres" value={data.database.status} tone="info" />
            </div>
          </Card>
        </section>

        <section id="files" className="workspace-section">
          <Card>
            <SectionTitle title="Files" description="Read-only explorer preview for the workspace." />
            <div className="file-explorer">
              {fileTree.map((item) => (
                <div key={item.path} className={`file-explorer__row file-explorer__row--${item.kind}`}>
                  <span>{item.kind === 'folder' ? 'Folder' : 'File'}</span>
                  <strong>{item.path}</strong>
                </div>
              ))}
            </div>
          </Card>
        </section>

        <section id="database" className="workspace-section">
          <Card>
            <SectionTitle title="Database" description="Workspace database connection and state." />
            <div className="detail-grid">
              <Item label="Database" value={data.database.name} />
              <Item label="Engine" value={data.database.engine} />
              <Item label="Status" value={data.database.status} />
            </div>
            <EmptyState title="Database management" description="Schema, migrations, and backups will appear here next." />
          </Card>
        </section>

        <section id="domains" className="workspace-section">
          <Card>
            <SectionTitle title="Domains" description="Public endpoints and SSL trust state." />
            <DataTable
              columns={[{ header: 'Host' }, { header: 'SSL' }, { header: 'Target' }]}
              rows={data.domains.map((domain) => [
                domain.host,
                <Badge key={domain.host} tone={domain.ssl === 'Valid' ? 'success' : domain.ssl === 'Renewing' ? 'warning' : 'info'}>
                  {domain.ssl}
                </Badge>,
                domain.target,
              ])}
            />
          </Card>
        </section>

        <section id="deployment" className="workspace-section">
          <Card>
            <SectionTitle title="Deployment" description="Recent release history for this workspace." />
            <DataTable
              columns={[{ header: 'Version' }, { header: 'Environment' }, { header: 'Status' }, { header: 'When' }, { header: 'Note' }]}
              rows={data.deployments.map((deployment) => [
                deployment.version,
                deployment.environment,
                <Badge key={deployment.id} tone={deployment.status === 'succeeded' ? 'success' : deployment.status === 'failed' ? 'danger' : 'warning'}>
                  {deployment.status}
                </Badge>,
                deployment.deployedAt,
                deployment.note,
              ])}
            />
          </Card>
        </section>

        <section id="monitoring" className="workspace-section">
          <Card>
            <SectionTitle title="Monitoring" description="Live readiness snapshot for the workspace." />
            <div className="metric-grid metric-grid--three">
              <Metric label="CPU" value="42%" tone="success" />
              <Metric label="RAM" value="68%" tone="warning" />
              <Metric label="Disk" value="51%" tone="info" />
            </div>
          </Card>
        </section>

        <section id="logs" className="workspace-section">
          <Card>
            <SectionTitle title="Logs" description="Recent project activity snapshots." />
            <div className="timeline">
              {data.logs.map((log) => (
                <div key={log.id} className="timeline__item">
                  <div className={`timeline__dot timeline__dot--${log.tone}`} />
                  <div className="timeline__body">
                    <div className="timeline__top">
                      <h4>{log.title}</h4>
                      <span>{log.time}</span>
                    </div>
                    <p>{log.detail}</p>
                  </div>
                </div>
              ))}
            </div>
          </Card>
        </section>

        <section id="ai" className="workspace-section">
          <Card>
            <SectionTitle title="AI" description="Workspace recommendations derived from metadata." />
            {data.recommendations.length ? (
              <div className="recommendation-list">
                {data.recommendations.map((item) => (
                  <div key={item.title} className="recommendation-card">
                    <div className="recommendation-card__title">{item.title}</div>
                    <p>{item.detail}</p>
                    <span>{item.reason}</span>
                  </div>
                ))}
              </div>
            ) : (
              <EmptyState title="No recommendations yet" description="AI suggestions will appear here once the workspace has more context." />
            )}
          </Card>
        </section>

        <section id="knowledge" className="workspace-section">
          <Card>
            <SectionTitle title="Knowledge" description="Learning cards and workspace guidance." />
            <LearnCard title="Workspace learning" articles={data.knowledge} />
          </Card>
        </section>

        <section id="doctor" className="workspace-section">
          <Card>
            <SectionTitle title="Doctor" description="Quick health and readiness checks." />
            <div className="check-list">
              <Check label="Repository" status="passed" detail="Connected and fetchable." />
              <Check label="Environment" status="warning" detail="One secret is missing from the snapshot." />
              <Check label="Deployment target" status="passed" detail="Healthy and ready to accept rollout." />
            </div>
          </Card>
        </section>

        <section id="settings" className="workspace-section">
          <Card>
            <SectionTitle title="Settings" description="Control deploy behavior, access, and safety." />
            <div className="stack">
              <div className="detail-grid">
                <Item label="Auto deploy" value="Enabled" />
                <Item label="Protection" value="Main branch only" />
                <Item label="Rollback policy" value="Latest successful release" />
                <Item label="Owner" value={data.project.owner} />
              </div>
              <Button variant="primary" type="button" onClick={() => setDialogOpen(true)}>
                Archive workspace
              </Button>
            </div>
          </Card>
        </section>
      </div>

      <Dialog
        open={dialogOpen}
        title="Archive workspace"
        description="This confirmation demonstrates the shared dialog system."
        onClose={() => setDialogOpen(false)}
      >
        <div className="stack">
          <p>This action is not wired to a real backend yet. It is safe to keep exploring the UI.</p>
          <Button
            variant="primary"
            onClick={() => {
              push({
                title: 'Workspace archived',
                message: 'The confirmation flow was demonstrated.',
                tone: 'info',
              })
              setDialogOpen(false)
            }}
            type="button"
          >
            Confirm archive
          </Button>
        </div>
      </Dialog>
    </div>
  )
}

function SectionTitle({ title, description }: { title: string; description: string }) {
  return (
    <div className="section-header-lite">
      <div className="card__eyebrow">Workspace</div>
      <h3>{title}</h3>
      <p>{description}</p>
    </div>
  )
}

function Item({ label, value }: { label: string; value: string }) {
  return (
    <div className="detail-item">
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  )
}

function Metric({ label, value, tone }: { label: string; value: string; tone: 'neutral' | 'accent' | 'success' | 'warning' | 'danger' | 'info' }) {
  return (
    <Card>
      <div className="card__eyebrow">{label}</div>
      <h3>{value}</h3>
      <Badge tone={tone}>{value}</Badge>
    </Card>
  )
}

function Check({ label, detail, status }: { label: string; detail: string; status: 'passed' | 'warning' | 'failed' }) {
  return (
    <div className="check-item">
      <div>
        <div className="check-item__label">
          <strong>{label}</strong>
          <Badge tone={status === 'passed' ? 'success' : status === 'warning' ? 'warning' : 'danger'}>{status}</Badge>
        </div>
        <p className="check-item__detail">{detail}</p>
      </div>
    </div>
  )
}
