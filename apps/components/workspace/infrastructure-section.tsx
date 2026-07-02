import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import type { InfrastructureInfo } from '@/lib/types'

interface Props { data: InfrastructureInfo }

export function InfrastructureSection({ data }: Props) {
  const installed = data.tools?.filter((t) => t.installed) || []
  const missing = data.tools?.filter((t) => !t.installed) || []

  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Infrastructure Overview</div>
          <h2 className="ws-section__title">Infrastructure</h2>
          <p className="ws-section__subtitle">{data.hostname} · {data.os}/{data.arch} · {installed.length} tools installed</p>
        </div>
        <Badge tone="info">Read-only</Badge>
      </div>

      <div className="ws-infra-grid">
        {(data.tools || []).map((tool) => (
          <Card key={tool.name} className={`ws-infra-card${tool.installed ? ' ws-infra-card--installed' : ''}`}>
            <div className="ws-infra-card__top">
              <strong>{tool.name}</strong>
              <Badge tone={tool.installed ? (tool.healthy ? 'success' : 'warning') : 'neutral'}>
                {tool.installed ? (tool.healthy ? 'Healthy' : 'Unhealthy') : 'Not installed'}
              </Badge>
            </div>
            {tool.installed && (
              <div className="ws-infra-card__details">
                <code className="ws-mono" style={{ fontSize: '0.82rem' }}>{tool.version}</code>
                <div className="ws-detail-row" style={{ marginTop: 8 }}>
                  <span>Configured</span>
                  <Badge tone={tool.configured ? 'success' : 'warning'}>{tool.configured ? 'Yes' : 'No'}</Badge>
                </div>
              </div>
            )}
          </Card>
        ))}
      </div>

      <div className="ws-grid-2">
        <Card>
          <div className="section-header-lite"><h3>Installed ({installed.length})</h3></div>
          <div className="ws-tag-list">
            {installed.map((t) => <Badge key={t.name} tone="success">{t.name}</Badge>)}
          </div>
        </Card>
        <Card>
          <div className="section-header-lite"><h3>Not Available ({missing.length})</h3></div>
          <div className="ws-tag-list">
            {missing.length ? missing.map((t) => <Badge key={t.name} tone="neutral">{t.name}</Badge>) : <span style={{ color: 'var(--muted)' }}>All tools detected</span>}
          </div>
        </Card>
      </div>
    </div>
  )
}
