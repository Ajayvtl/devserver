'use client'

import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import type { InfrastructureInfo } from '@/lib/types'
import { submitCommand } from '@/lib/api/client'

interface Props { data: InfrastructureInfo }

export function InfrastructureSection({ data }: Props) {
  const installed = data.tools?.filter((t) => t.installed) || []
  const missing = data.tools?.filter((t) => !t.installed) || []

  const handleAction = async (action: string, target: string) => {
    try {
      await submitCommand({
        capability: `service.${action}`,
        target: target,
        name: `${action} ${target}`,
      })
    } catch (err) {
      console.error(err)
    }
  }

  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Infrastructure Overview</div>
          <h2 className="ws-section__title">Infrastructure</h2>
          <p className="ws-section__subtitle">{data.hostname} · {data.os}/{data.arch} · {installed.length} tools installed</p>
        </div>
      </div>

      <div className="ws-infra-grid">
        {(data.tools || []).map((tool) => (
          <Card key={tool.name} className={`ws-infra-card${tool.installed ? ' ws-infra-card--installed' : ''}`} style={{ display: 'flex', flexDirection: 'column' }}>
            <div className="ws-infra-card__top">
              <strong style={{ fontSize: '1.1rem' }}>{tool.name}</strong>
              <Badge tone={tool.installed ? (tool.healthy ? 'success' : 'warning') : 'neutral'}>
                {tool.installed ? (tool.healthy ? 'Healthy' : 'Unhealthy') : 'Not installed'}
              </Badge>
            </div>
            
            <div className="ws-infra-card__details" style={{ flex: 1 }}>
              {tool.installed ? (
                <div style={{ marginBottom: '16px' }}>
                  <code className="ws-mono" style={{ fontSize: '0.85rem', display: 'block', marginBottom: '8px' }}>Version {tool.version}</code>
                  <div className="ws-detail-row">
                    <span>Configured</span>
                    <Badge tone={tool.configured ? 'success' : 'warning'}>{tool.configured ? 'Yes' : 'No'}</Badge>
                  </div>
                </div>
              ) : (
                <div style={{ marginBottom: '16px', color: 'var(--muted)', fontSize: '0.85rem' }}>
                  This tool is not installed on the system.
                </div>
              )}
            </div>
            
            <div style={{ borderTop: '1px solid rgba(148, 163, 184, 0.08)', paddingTop: '12px', display: 'flex', gap: '8px', flexWrap: 'wrap' }}>
              {tool.installed ? (
                <>
                  <button className="ws-btn" onClick={() => handleAction('configure', tool.name)} style={{ background: 'rgba(255, 255, 255, 0.05)', color: 'var(--text)', border: '1px solid rgba(255,255,255,0.1)', padding: '6px 12px', borderRadius: '6px', fontSize: '0.8rem', cursor: 'pointer', flex: 1 }}>
                    Configure
                  </button>
                  <button className="ws-btn" onClick={() => handleAction('update', tool.name)} style={{ background: 'rgba(255, 255, 255, 0.05)', color: 'var(--text)', border: '1px solid rgba(255,255,255,0.1)', padding: '6px 12px', borderRadius: '6px', fontSize: '0.8rem', cursor: 'pointer', flex: 1 }}>
                    Update
                  </button>
                  <button className="ws-btn" onClick={() => handleAction('restart', tool.name)} style={{ background: 'rgba(255, 255, 255, 0.05)', color: 'var(--text)', border: '1px solid rgba(255,255,255,0.1)', padding: '6px 12px', borderRadius: '6px', fontSize: '0.8rem', cursor: 'pointer', flex: 1 }}>
                    Restart
                  </button>
                </>
              ) : (
                <button className="ws-btn" onClick={() => handleAction('install', tool.name)} style={{ background: 'var(--accent)', color: '#000', border: 'none', padding: '6px 12px', borderRadius: '6px', fontSize: '0.8rem', cursor: 'pointer', fontWeight: 600, width: '100%' }}>
                  Install {tool.name}
                </button>
              )}
            </div>
          </Card>
        ))}
      </div>
    </div>
  )
}
