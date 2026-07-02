import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { Progress } from '../ui/progress'
import { EmptyState } from '../ui/empty-state'
import type { WorkspaceOverview } from '@/lib/types'

interface Props { data: WorkspaceOverview }

export function OverviewSection({ data }: Props) {
  const { workspace: ws, project, health, git, plugins } = data
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Workspace</div>
          <h2 className="ws-section__title">{ws.name}</h2>
          <p className="ws-section__subtitle">{ws.framework} · {ws.runtime} · {ws.languages?.join(', ')}</p>
        </div>
        <Badge tone={health.score >= 60 ? 'success' : health.score >= 30 ? 'warning' : 'danger'}>
          Health {health.score}
        </Badge>
      </div>

      <div className="metric-grid--three" style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: 14 }}>
        <Card>
          <div className="card__eyebrow">Framework</div>
          <strong>{ws.framework || 'Unknown'}</strong>
          <p style={{ margin: '6px 0 0', color: 'var(--muted)' }}>{ws.kind}</p>
        </Card>
        <Card>
          <div className="card__eyebrow">Runtime</div>
          <strong>{ws.runtime || 'Unknown'}</strong>
          <p style={{ margin: '6px 0 0', color: 'var(--muted)' }}>{ws.packageManager}</p>
        </Card>
        <Card>
          <div className="card__eyebrow">Repository</div>
          <strong>{project.repository || 'None'}</strong>
          <p style={{ margin: '6px 0 0', color: 'var(--muted)' }}>{git.branch || 'No branch'}</p>
        </Card>
      </div>

      <div className="ws-grid-2">
        <Card>
          <div className="section-header-lite"><h3>Git Status</h3></div>
          <div className="ws-detail-list">
            <div className="ws-detail-row"><span>Branch</span><Badge tone="accent">{git.branch || '—'}</Badge></div>
            <div className="ws-detail-row"><span>Commit</span><code className="ws-mono">{git.commit ? git.commit.slice(0, 8) : '—'}</code></div>
            <div className="ws-detail-row"><span>State</span><Badge tone={git.clean ? 'success' : 'warning'}>{git.clean ? 'Clean' : 'Dirty'}</Badge></div>
            <div className="ws-detail-row"><span>Ahead / Behind</span><span>{git.ahead} / {git.behind}</span></div>
          </div>
        </Card>
        <Card>
          <div className="section-header-lite"><h3>Health</h3></div>
          <div className="ws-detail-list">
            <div className="ws-detail-row"><span>Score</span><strong>{health.score}</strong></div>
            <div className="ws-detail-row"><span>Build</span><Badge tone="neutral">{health.build}</Badge></div>
            <div className="ws-detail-row"><span>Tests</span><Badge tone="neutral">{health.tests}</Badge></div>
            <div className="ws-detail-row"><span>Lint</span><Badge tone="neutral">{health.lint}</Badge></div>
          </div>
          <Progress value={health.score} tone={health.score >= 60 ? 'success' : 'accent'} />
        </Card>
      </div>

      {plugins.detected?.length > 0 && (
        <Card>
          <div className="section-header-lite"><h3>Detected Plugins</h3></div>
          <div className="ws-tag-list">
            {plugins.detected.map((p) => <Badge key={p} tone="info">{p}</Badge>)}
          </div>
          {plugins.capabilities?.length > 0 && (
            <div className="ws-tag-list" style={{ marginTop: 10 }}>
              {plugins.capabilities.map((c) => <Badge key={c} tone="neutral">{c}</Badge>)}
            </div>
          )}
        </Card>
      )}
    </div>
  )
}
