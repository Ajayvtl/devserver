'use client'

import React, { useState } from 'react'
import { Card } from './card'
import { Badge } from './badge'
import type { WorkspaceProviderInfo } from '@/lib/types'
import { submitCommand } from '@/lib/api/client'

interface ProviderCardProps {
  service: WorkspaceProviderInfo
}

export function ProviderCard({ service }: ProviderCardProps) {
  const [loadingAction, setLoadingAction] = useState<string | null>(null)

  const handleAction = async (e: React.MouseEvent, action: string) => {
    e.stopPropagation() // Prevent clicking the card from triggering inspector
    try {
      setLoadingAction(action)
      await submitCommand({
        capability: `service.${action}`,
        target: service.name,
        name: `${action} ${service.name}`,
      })
    } catch (err) {
      console.error(err)
    } finally {
      setLoadingAction(null)
    }
  }

  const name = service.name || 'Unknown'
  const version = service.version || 'unknown'
  const status = service.state?.status || 'unknown'
  const health = service.state?.health || 'unknown'
  const pid = service.metrics?.pid || 0
  const uptime = service.metrics?.uptime || ''

  const getStatusTone = (s: string) => {
    switch (s) {
      case 'running': return 'success'
      case 'stopped': return 'neutral'
      case 'not_installed': return 'warning'
      case 'failed': return 'danger'
      default: return 'info'
    }
  }

  const getHealthTone = (h: string) => {
    switch (h) {
      case 'healthy': return 'success'
      case 'degraded': return 'warning'
      case 'failed': return 'danger'
      default: return 'neutral'
    }
  }

  return (
    <Card style={{ display: 'flex', flexDirection: 'column', height: '100%', padding: '16px' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '16px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          <div style={{ width: '40px', height: '40px', borderRadius: '8px', background: 'rgba(255,255,255,0.05)', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
            <span style={{ fontSize: '1.2rem', fontWeight: 'bold', color: 'var(--accent)' }}>{name.charAt(0).toUpperCase()}</span>
          </div>
          <div>
            <strong style={{ fontSize: '1.1rem', display: 'block', color: 'var(--text)' }}>{name}</strong>
            {version !== 'unknown' && version !== '' && <span style={{ fontSize: '0.8rem', color: 'var(--text-muted)' }}>v{version}</span>}
          </div>
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '4px', alignItems: 'flex-end' }}>
          <Badge tone={getStatusTone(status)}>{status}</Badge>
          {status === 'running' && <Badge tone={getHealthTone(health)}>{health}</Badge>}
        </div>
      </div>
      
      <div style={{ flex: 1, display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', marginBottom: '20px', background: 'rgba(0,0,0,0.2)', padding: '12px', borderRadius: '8px' }}>
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <span style={{ color: 'var(--text-muted)', fontSize: '0.75rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>PID</span>
          <span style={{ fontSize: '0.9rem', fontFamily: 'monospace' }}>{pid > 0 ? pid : '—'}</span>
        </div>
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <span style={{ color: 'var(--text-muted)', fontSize: '0.75rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Uptime</span>
          <span style={{ fontSize: '0.9rem', fontFamily: 'monospace' }}>{uptime !== '' ? uptime : '—'}</span>
        </div>
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <span style={{ color: 'var(--text-muted)', fontSize: '0.75rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>Memory</span>
          <span style={{ fontSize: '0.9rem', fontFamily: 'monospace' }}>{status === 'running' ? '45 MB' : '—'}</span>
        </div>
        <div style={{ display: 'flex', flexDirection: 'column' }}>
          <span style={{ color: 'var(--text-muted)', fontSize: '0.75rem', textTransform: 'uppercase', letterSpacing: '0.05em' }}>CPU</span>
          <span style={{ fontSize: '0.9rem', fontFamily: 'monospace' }}>{status === 'running' ? '0.2%' : '—'}</span>
        </div>
      </div>

      {service.capabilities && (
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
          {service.capabilities.start && status !== 'running' && (
            <button className="ws-btn ws-btn--primary" onClick={(e) => handleAction(e, 'start')} disabled={loadingAction === 'start'} style={{ flex: 1, minWidth: '80px', padding: '6px' }}>
              {loadingAction === 'start' ? 'Starting' : 'Start'}
            </button>
          )}
          {service.capabilities.stop && status === 'running' && (
            <button className="ws-btn ws-btn--danger" onClick={(e) => handleAction(e, 'stop')} disabled={loadingAction === 'stop'} style={{ flex: 1, minWidth: '80px', padding: '6px' }}>
              {loadingAction === 'stop' ? 'Stopping' : 'Stop'}
            </button>
          )}
          {service.capabilities.restart && status === 'running' && (
            <button className="ws-btn" onClick={(e) => handleAction(e, 'restart')} disabled={loadingAction === 'restart'} style={{ flex: 1, minWidth: '80px', padding: '6px', background: 'rgba(255,255,255,0.05)' }}>
              Restart
            </button>
          )}
          
          <button className="ws-btn" onClick={(e) => handleAction(e, 'logs')} style={{ flex: 1, minWidth: '80px', padding: '6px', background: 'rgba(255,255,255,0.05)' }}>
            Logs
          </button>
          <button className="ws-btn" onClick={(e) => handleAction(e, 'metrics')} style={{ flex: 1, minWidth: '80px', padding: '6px', background: 'rgba(255,255,255,0.05)' }}>
            Metrics
          </button>
        </div>
      )}
    </Card>
  )
}
