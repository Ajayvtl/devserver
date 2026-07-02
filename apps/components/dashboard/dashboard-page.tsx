import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import { LearnCard } from '../ui/learn-card'
import { MetricCard } from '../ui/metric-card'
import { Progress } from '../ui/progress'

import type { DashboardDataV2, KnowledgeArticle } from '@/lib/types'

interface Props {
  data: DashboardDataV2
  knowledge: KnowledgeArticle[]
}

export function DashboardPageContent({ data, knowledge }: Props) {
  return (
    <div className="dashboard-layout">
      <section id="overview" className="dashboard-hero">
        <div className="dashboard-hero__copy">
          <Badge tone="success">Live.</Badge>
          <h1>GitHub clarity, Linear motion, Railway speed.</h1>
          <p>
            The platform is responsive, mock-driven, and ready to wire into real APIs without changing page layout.
          </p>
          <div className="dashboard-hero__actions">
            <Button variant="primary" href="#deployments">
              New deployment
            </Button>
            <Button variant="secondary" href="#terminal">
              Open terminal
            </Button>
          </div>
        </div>

        <Card className="dashboard-hero__panel">
          <div className="card__eyebrow">Current server</div>
          <h3>{data.server}</h3>
          <p>All metrics below are backed by service interfaces and mock data today.</p>
          <div className="dashboard-hero__pulse">
            <span />
            <span />
            <span />
          </div>
        </Card>
      </section>

      <section className="section-block">
        <div className="section-block__header">
          <div>
            <div className="card__eyebrow">Section 1</div>
            <h2>CPU, RAM, Disk, Network</h2>
          </div>
          <Badge tone="accent">Live.</Badge>
        </div>
        <div className="metric-grid">
          {data.metrics.map((metric) => (
            <MetricCard
              key={metric.label}
              label={metric.label}
              value={metric.value}
              detail={metric.detail}
              trend={metric.trend}
              tone={metric.tone}
            />
          ))}
        </div>
      </section>

      <div className="dashboard-grid">
        <section id="projects" className="dashboard-panel">
          <PanelCard title="Projects" detail="Active workloads and release targets." items={data.sections[0].items} />
        </section>

        <section id="domains" className="dashboard-panel">
          <PanelCard title="Domains" detail="DNS, certificates, and edge exposure." items={data.sections[1].items} />
        </section>

        <section id="services" className="dashboard-panel">
          <PanelCard title="Services" detail="System and app service state." items={data.sections[2].items} />
        </section>

        <section id="deployments" className="dashboard-panel">
          <PanelCard title="Deployments" detail="Most recent rollout activity." items={data.sections[3].items} />
        </section>
      </div>

      <div className="dashboard-grid dashboard-grid--wide">
        <Card id="activity" className="dashboard-section">
          <SectionHeader title="Recent Activity" eyebrow="Section 3" />
          <Timeline items={data.activity} />
        </Card>

        <Card id="tasks" className="dashboard-section">
          <SectionHeader title="Task Queue" eyebrow="Section 4" />
          {data.tasks.length ? (
            <div className="task-list">
              {data.tasks.map((task) => (
                <div key={task.title} className="task-row">
                  <div className="task-row__top">
                    <strong>{task.title}</strong>
                    <Badge tone={task.state === 'Done' ? 'success' : task.state === 'Running' ? 'accent' : 'warning'}>
                      {task.state}
                    </Badge>
                  </div>
                  <p>{task.detail}</p>
                  <Progress value={task.progress} tone={task.state === 'Done' ? 'success' : 'accent'} />
                </div>
              ))}
            </div>
          ) : (
            <EmptyState title="No queued tasks" description="New operations will appear here once you start deploying or installing modules." />
          )}
        </Card>
      </div>

      <section className="section-block">
        <div className="section-block__header">
          <div>
            <div className="card__eyebrow">Workspace</div>
            <h2>Servers, SSL, Databases, Storage, Terminal, Monitoring, Logs, Users, Settings</h2>
          </div>
          <Badge tone="info">Navigation mapped</Badge>
        </div>

        <div className="system-grid">
          {systemTiles.map((tile) => (
            <Card key={tile.id} id={tile.id} className="system-tile">
              <div className="section-header-lite">
                <div className="card__eyebrow">{tile.eyebrow}</div>
                <h3>{tile.title}</h3>
                <p>{tile.detail}</p>
              </div>
              <Badge tone={tile.tone}>{tile.value}</Badge>
            </Card>
          ))}
        </div>
      </section>

      <div className="dashboard-grid dashboard-grid--wide">
        <Card id="ai" className="dashboard-section">
          <SectionHeader title="AI Recommendations" eyebrow="Section 5" />
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
            <EmptyState title="No recommendations yet" description="AI suggestions will appear here once the platform has context." />
          )}
        </Card>

        <aside className="dashboard-rail">
          <LearnCard title="Learn the platform" articles={knowledge} />
          <Card>
            <SectionHeader title="Monitoring" eyebrow="Sections" />
            <div className="rail-links">
              <a href="#servers">Servers</a>
              <a href="#ssl">SSL</a>
              <a href="#databases">Databases</a>
              <a href="#storage">Storage</a>
              <a href="#terminal">Terminal</a>
              <a href="#monitoring">Monitoring</a>
              <a href="#logs">Logs</a>
              <a href="#users">Users</a>
              <a href="#settings">Settings</a>
            </div>
          </Card>
        </aside>
      </div>
    </div>
  )
}

function PanelCard({
  title,
  detail,
  items,
}: {
  title: string
  detail: string
  items: Array<{ label: string; value: string; tone?: 'neutral' | 'accent' | 'success' | 'warning' | 'danger' | 'info' }>
}) {
  return (
    <Card>
      <SectionHeader title={title} description={detail} />
      <div className="panel-list">
        {items.length ? (
          items.map((item) => (
            <div key={item.label} className="panel-list__item">
              <span>{item.label}</span>
              <Badge tone={item.tone ?? 'neutral'}>{item.value}</Badge>
            </div>
          ))
        ) : (
          <EmptyState title="Nothing here yet" description="This area will populate when real APIs are connected." />
        )}
      </div>
    </Card>
  )
}

function Timeline({
  items,
}: {
  items: Array<{ title: string; detail: string; when: string; tone: 'neutral' | 'accent' | 'success' | 'warning' | 'danger' | 'info' }>
}) {
  return (
    <div className="timeline">
      {items.map((item) => (
        <div key={item.title} className="timeline__item">
          <div className={`timeline__dot timeline__dot--${item.tone}`} />
          <div className="timeline__body">
            <div className="timeline__top">
              <h4>{item.title}</h4>
              <span>{item.when}</span>
            </div>
            <p>{item.detail}</p>
          </div>
        </div>
      ))}
    </div>
  )
}

function SectionHeader({ title, eyebrow, description }: { title: string; eyebrow?: string; description?: string }) {
  return (
    <div className="section-header-lite">
      {eyebrow ? <div className="card__eyebrow">{eyebrow}</div> : null}
      <h3>{title}</h3>
      {description ? <p>{description}</p> : null}
    </div>
  )
}

const systemTiles = [
  { id: 'servers', eyebrow: 'Servers', title: 'Fleet overview', detail: 'Current hosts, roles, and health.', value: '8 online', tone: 'success' as const },
  { id: 'ssl', eyebrow: 'SSL', title: 'Certificates', detail: 'Renewal windows and trust coverage.', value: '4 expiring soon', tone: 'warning' as const },
  { id: 'databases', eyebrow: 'Databases', title: 'Platform data', detail: 'Users, projects, plugins, and tasks.', value: 'Ready', tone: 'accent' as const },
  { id: 'storage', eyebrow: 'Storage', title: 'Snapshots', detail: 'Volumes and retention policies.', value: 'Healthy', tone: 'info' as const },
  { id: 'terminal', eyebrow: 'Terminal', title: 'Command access', detail: 'Local and remote execution surface.', value: 'Available', tone: 'neutral' as const },
  { id: 'monitoring', eyebrow: 'Monitoring', title: 'Live telemetry', detail: 'Metrics, alerts, and thresholds.', value: 'Watching', tone: 'success' as const },
  { id: 'logs', eyebrow: 'Logs', title: 'Audit trail', detail: 'User and system activity snapshots.', value: 'Indexed', tone: 'info' as const },
  { id: 'users', eyebrow: 'Users', title: 'Access control', detail: 'Roles, sessions, and invitations.', value: 'RBAC', tone: 'accent' as const },
  { id: 'settings', eyebrow: 'Settings', title: 'Platform config', detail: 'Defaults, preferences, and policies.', value: 'Ready', tone: 'neutral' as const },
]
