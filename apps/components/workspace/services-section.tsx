'use client'

import React, { useState, useMemo } from 'react'
import { EmptyState } from '../ui/empty-state'
import { ProviderCard } from '../ui/provider-card'
import { Badge } from '../ui/badge'
import type { WorkspaceProviderInfo } from '@/lib/types'

interface Props { data: WorkspaceProviderInfo[] }

export function ServicesSection({ data }: Props) {
  const [search, setSearch] = useState('')
  const [filterStatus, setFilterStatus] = useState<string>('all')
  const [selectedProvider, setSelectedProvider] = useState<WorkspaceProviderInfo | null>(null)

  const filteredData = useMemo(() => {
    if (!data) return []
    return data.filter(s => {
      const name = s.name || ''
      const status = s.state?.status || ''
      const matchSearch = name.toLowerCase().includes(search.toLowerCase())
      const matchStatus = filterStatus === 'all' || status === filterStatus
      return matchSearch && matchStatus
    })
  }, [data, search, filterStatus])

  return (
    <div className="ws-section" style={{ display: 'flex', flexWrap: 'wrap', gap: '24px' }}>
      <div style={{ flex: '1 1 min(100%, 800px)' }}>
        <div className="ws-section__header">
          <div>
            <div className="card__eyebrow">Provider Registry</div>
            <h2 className="ws-section__title">Services</h2>
            <p className="ws-section__subtitle">{data?.length || 0} active providers detected</p>
          </div>
        </div>

        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '12px', marginBottom: '24px' }}>
          <input 
            type="text" 
            placeholder="Search providers..." 
            value={search}
            onChange={e => setSearch(e.target.value)}
            style={{ flex: '1 1 200px', padding: '8px 12px', borderRadius: '6px', border: '1px solid rgba(255,255,255,0.1)', background: 'rgba(0,0,0,0.2)', color: 'var(--text)' }}
          />
          <select 
            value={filterStatus}
            onChange={e => setFilterStatus(e.target.value)}
            style={{ flex: '0 0 auto', padding: '8px 12px', borderRadius: '6px', border: '1px solid rgba(255,255,255,0.1)', background: 'rgba(0,0,0,0.2)', color: 'var(--text)' }}
          >
            <option value="all">All Statuses</option>
            <option value="running">Running</option>
            <option value="stopped">Stopped</option>
            <option value="not_installed">Not Installed</option>
            <option value="failed">Failed</option>
          </select>
        </div>

        {filteredData.length ? (
          <div className="ws-service-grid" style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
            gap: '16px'
          }}>
            {filteredData.map((s) => (
              <div key={s.name} onClick={() => setSelectedProvider(s)} style={{ cursor: 'pointer', height: '100%' }}>
                <ProviderCard service={s} />
              </div>
            ))}
          </div>
        ) : (
          <EmptyState title="No providers match" description="Try adjusting your search or filters." />
        )}
      </div>

      {selectedProvider && (
        <div style={{ flex: '1 1 300px', minWidth: '300px', borderLeft: '1px solid rgba(255,255,255,0.1)', paddingLeft: '24px', paddingTop: '16px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
            <h3 style={{ margin: 0 }}>Inspector</h3>
            <button 
              onClick={() => setSelectedProvider(null)} 
              style={{ background: 'transparent', border: 'none', color: 'var(--text-muted)', cursor: 'pointer', fontSize: '1.2rem' }}
            >
              &times;
            </button>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <span style={{ color: 'var(--text-muted)' }}>Name</span>
            <div>{selectedProvider.name}</div>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <span style={{ color: 'var(--text-muted)' }}>Version</span>
            <div>{selectedProvider.version || 'Unknown'}</div>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <span style={{ color: 'var(--text-muted)' }}>Status</span>
            <div><Badge tone={selectedProvider.state?.status === 'running' ? 'success' : 'neutral'}>{selectedProvider.state?.status || 'unknown'}</Badge></div>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <span style={{ color: 'var(--text-muted)' }}>Health</span>
            <div><Badge tone={selectedProvider.state?.health === 'healthy' ? 'success' : 'warning'}>{selectedProvider.state?.health || 'unknown'}</Badge></div>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <span style={{ color: 'var(--text-muted)' }}>PID</span>
            <div>{(selectedProvider.metrics?.pid || 0) > 0 ? selectedProvider.metrics.pid : 'N/A'}</div>
          </div>
          <div style={{ marginBottom: '12px' }}>
            <span style={{ color: 'var(--text-muted)' }}>Uptime</span>
            <div>{(selectedProvider.metrics?.uptime || '') !== '' ? selectedProvider.metrics.uptime : 'N/A'}</div>
          </div>
          
          <h4 style={{ marginTop: '24px', marginBottom: '8px' }}>Capabilities</h4>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
            {Object.entries(selectedProvider.capabilities || {}).filter(([_, v]) => v).map(([k]) => (
              <Badge key={k} tone="neutral">{k}</Badge>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}
