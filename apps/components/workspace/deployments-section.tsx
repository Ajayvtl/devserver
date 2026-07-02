import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import type { DeploymentInfo } from '@/lib/types'

interface Props { data: DeploymentInfo }

const steps = ['Repository', 'Branch', 'Build', 'Deploy', 'Rollback']

export function DeploymentsSection({ data }: Props) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Deployment Center</div>
          <h2 className="ws-section__title">Deployments</h2>
          <p className="ws-section__subtitle">Visual workflow — execution not yet enabled</p>
        </div>
        <Badge tone="info">UI Only</Badge>
      </div>

      <Card>
        <div className="section-header-lite"><h3>Deployment Pipeline</h3></div>
        <div className="ws-pipeline">
          {steps.map((step, i) => (
            <div key={step} className="ws-pipeline__step">
              <div className={`ws-pipeline__node${i < 4 ? ' ws-pipeline__node--done' : ''}`}>{i + 1}</div>
              <span className="ws-pipeline__label">{step}</span>
              {i < steps.length - 1 && <div className="ws-pipeline__arrow">→</div>}
            </div>
          ))}
        </div>
      </Card>

      <Card>
        <div className="section-header-lite"><h3>Deployment History</h3></div>
        {data.entries?.length ? (
          <div className="ws-deploy-list">
            {data.entries.map((entry) => (
              <div key={entry.id} className="ws-deploy-row">
                <div className="ws-deploy-row__top">
                  <code className="ws-mono">{entry.commit}</code>
                  <Badge tone={entry.status === 'succeeded' ? 'success' : entry.status === 'running' ? 'accent' : 'danger'}>{entry.status}</Badge>
                </div>
                <div className="ws-deploy-row__meta">
                  <span>{entry.branch}</span>
                  <span>{entry.environment}</span>
                  <span style={{ color: 'var(--muted)' }}>{entry.timestamp}</span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <EmptyState title="No deployments" description="Deployment history will appear here once the deployment engine is enabled." />
        )}
      </Card>
    </div>
  )
}
