import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import type { WorkspaceServiceInfo } from '@/lib/types'

interface Props { data: WorkspaceServiceInfo[] }

export function ServicesSection({ data }: Props) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Service Registry</div>
          <h2 className="ws-section__title">Services</h2>
          <p className="ws-section__subtitle">{data?.length || 0} detected services</p>
        </div>
      </div>
      {data?.length ? (
        <div className="ws-service-grid">
          {data.map((s) => (
            <Card key={s.name}>
              <div className="ws-infra-card__top">
                <strong>{s.name}</strong>
                <Badge tone={s.status === 'running' ? 'success' : s.status === 'installed' ? 'info' : 'danger'}>{s.status}</Badge>
              </div>
              {s.port && <div className="ws-detail-row" style={{ marginTop: 8 }}><span>Port</span><code className="ws-mono">{s.port}</code></div>}
              {s.pid && <div className="ws-detail-row"><span>PID</span><code className="ws-mono">{s.pid}</code></div>}
            </Card>
          ))}
        </div>
      ) : (
        <EmptyState title="No services detected" description="Service detection runs when infrastructure tools like Redis, PostgreSQL, or Nginx are installed." />
      )}
    </div>
  )
}
