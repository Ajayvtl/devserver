'use client'

import { useState, useEffect } from 'react'
import { 
  CheckCircle2, 
  AlertTriangle, 
  Play, 
  Shield, 
  RefreshCw, 
  Star, 
  Info, 
  Trash2, 
  Globe, 
  Cpu, 
  Terminal, 
  Laptop, 
  Activity,
  Layers,
  ChevronDown,
  ChevronUp
} from 'lucide-react'

import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import { Progress } from '../ui/progress'

import type { DashboardDataV2, KnowledgeArticle } from '@/lib/types'
import { useAuth } from '@/components/auth/auth-context'
import { requestOrFallback } from '@/lib/api/client'

interface Props {
  data: DashboardDataV2
  knowledge: KnowledgeArticle[]
}

export function DashboardPageContent({ data, knowledge }: Props) {
  const { user, currentOrgId } = useAuth()
  const [systemOrgs, setSystemOrgs] = useState<any[]>([])
  const [loadingOrgs, setLoadingOrgs] = useState(false)

  const isSuperAdmin = user?.role === 'admin' || user?.role === 'super admin' || user?.role === 'owner'

  useEffect(() => {
    if (isSuperAdmin && currentOrgId) {
      setLoadingOrgs(true)
      requestOrFallback<any>('/api/v1/superadmin/organizations', { data: [] }, {
        headers: { 'X-Org-ID': currentOrgId }
      })
      .then(res => {
        const list = Array.isArray(res) ? res : (res?.data || [])
        setSystemOrgs(list)
      })
      .catch(err => {
        console.error('Failed to load system organizations', err)
      })
      .finally(() => setLoadingOrgs(false))
    }
  }, [isSuperAdmin, currentOrgId])

  // Pinned/Favorites workspaces
  const [pinnedWorkspaces, setPinnedWorkspaces] = useState<string[]>(['devserver'])
  
  // Expand/Collapse status deck pills
  const [expandedPill, setExpandedPill] = useState<'health' | 'ai' | 'pipelines' | null>(null)
  
  // Toggles for sections
  const [showResourceMetrics, setShowResourceMetrics] = useState(false)

  // Security widget sessions data
  const [sessions, setSessions] = useState([
    { id: 'sess_1', device: 'Chrome 126.0 (Windows 11)', ip: '10.0.0.1', geo: 'Bengaluru, India', lastActive: 'Current Session' },
    { id: 'sess_2', device: 'Safari (Apple iPad)', ip: '10.0.0.42', geo: 'Bengaluru, India', lastActive: '2h ago' },
    { id: 'sess_3', device: 'VSCode Client', ip: '192.168.1.15', geo: 'Local Network', lastActive: '1d ago' }
  ])

  // Mock workspace detailed metadata for "Continue Working"
  const activeWorkspaceMeta = {
    name: 'devserver',
    repo: 'https://github.com/Ajayvtl/devserver.git',
    branch: 'feature/indexer-stabilization',
    gitState: '3 files modified, 0 untracked',
    buildStatus: 'Passing',
    testResults: '48 / 48 passed',
    coverage: '92.4%',
    runtime: 'Go 1.22.4, Next.js 15.5.20'
  }

  const togglePin = (name: string) => {
    setPinnedWorkspaces(prev => 
      prev.includes(name) ? prev.filter(w => w !== name) : [...prev, name]
    )
  }

  const togglePill = (pill: 'health' | 'ai' | 'pipelines') => {
    setExpandedPill(prev => prev === pill ? null : pill)
  }

  const revokeSession = (id: string) => {
    setSessions(prev => prev.filter(s => s.id !== id))
  }

  return (
    <div className="dashboard-layout" style={{ display: 'grid', gap: '16px' }}>
      
      {/* 1. COMPACT STATUS PILLS (REPLACES OVERSIZED CARDS) */}
      <section className="section-block">
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(240px, 1fr))', gap: '10px' }}>
          
          {/* Health Pill */}
          <div>
            <div className="status-deck-pill" onClick={() => togglePill('health')}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <CheckCircle2 size={16} style={{ color: 'var(--success)' }} />
                <span style={{ fontWeight: 600, fontSize: '0.85rem' }}>Platform: 98% nominal</span>
              </div>
              {expandedPill === 'health' ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
            </div>
            {expandedPill === 'health' && (
              <div className="status-deck-expanded">
                <strong>Fleet Status Details:</strong>
                <div style={{ marginTop: '6px', lineHeight: '1.4' }}>
                  • 8/8 executors connected & online.<br />
                  • Docker provider: Healthy (v26.1.1)<br />
                  • SSH adapters: 2 active connections.
                </div>
              </div>
            )}
          </div>

          {/* AI Pill */}
          <div>
            <div className="status-deck-pill" onClick={() => togglePill('ai')}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <Cpu size={16} style={{ color: 'var(--accent)' }} />
                <span style={{ fontWeight: 600, fontSize: '0.85rem' }}>AI Copilot: Ollama (codellama)</span>
              </div>
              {expandedPill === 'ai' ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
            </div>
            {expandedPill === 'ai' && (
              <div className="status-deck-expanded">
                <strong>Ollama Local Facade:</strong>
                <div style={{ marginTop: '6px', lineHeight: '1.4' }}>
                  • Endpoint: `http://localhost:11434`<br />
                  • Response Latency: 45ms (cached)<br />
                  • Embeddings index status: Indexed (1,482 blocks).
                </div>
              </div>
            )}
          </div>

          {/* Pipelines Pill */}
          <div>
            <div className="status-deck-pill" onClick={() => togglePill('pipelines')}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <Layers size={16} style={{ color: 'var(--info)' }} />
                <span style={{ fontWeight: 600, fontSize: '0.85rem' }}>Pipelines: 0 blocked</span>
              </div>
              {expandedPill === 'pipelines' ? <ChevronUp size={14} /> : <ChevronDown size={14} />}
            </div>
            {expandedPill === 'pipelines' && (
              <div className="status-deck-expanded">
                <strong>Active Deployments status:</strong>
                <div style={{ marginTop: '6px', lineHeight: '1.4' }}>
                  • Production: v1.2.3 deployed successfully 4h ago.<br />
                  • Staging Promotion: Waiting for manual trigger.<br />
                  • Auto-rollback rules: Enabled.
                </div>
              </div>
            )}
          </div>

        </div>
      </section>

      {/* 2. DEVELOPER WORKFLOWS PRIORITIZED */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: '16px' }}>
        
        {/* Continue Working & Workspace Info Summary */}
        <Card style={{ padding: '14px', display: 'flex', flexDirection: 'column', justifyContent: 'space-between' }}>
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
              <span className="card__eyebrow" style={{ color: 'var(--accent)' }}>Continue Working</span>
              <Badge tone="accent">Active Workspace</Badge>
            </div>
            <h3 style={{ margin: '0 0 8px', fontSize: '1.1rem', fontFamily: 'var(--font-heading)' }}>
              {activeWorkspaceMeta.name}
            </h3>
            
            <div style={{ display: 'grid', gap: '6px', fontSize: '0.8rem', color: 'var(--muted-strong)' }}>
              <div><strong>Repo:</strong> <code style={{ fontSize: '0.75rem' }}>{activeWorkspaceMeta.repo}</code></div>
              <div><strong>Branch:</strong> <code>{activeWorkspaceMeta.branch}</code></div>
              <div><strong>Git State:</strong> <span style={{ color: 'var(--warning)' }}>{activeWorkspaceMeta.gitState}</span></div>
              <div style={{ display: 'flex', gap: '14px', marginTop: '4px' }}>
                <span><strong>Build:</strong> <span style={{ color: 'var(--success)' }}>{activeWorkspaceMeta.buildStatus}</span></span>
                <span><strong>Tests:</strong> {activeWorkspaceMeta.testResults}</span>
                <span><strong>Coverage:</strong> {activeWorkspaceMeta.coverage}</span>
              </div>
              <div><strong>Runtime env:</strong> <span style={{ color: 'var(--info)' }}>{activeWorkspaceMeta.runtime}</span></div>
            </div>
          </div>

          <div style={{ marginTop: '14px', width: '100%' }}>
            <Button variant="primary" href={`/workspace/${activeWorkspaceMeta.name}`}>
              <span style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '6px', width: '100%', fontSize: '0.85rem' }}>
                <Play size={12} /> Resume Workspace IDE
              </span>
            </Button>
          </div>
        </Card>

        {/* Workspaces & Git status list */}
        <Card style={{ padding: '14px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
            <h3 style={{ margin: 0, fontSize: '0.95rem' }}>Workspaces</h3>
            <span style={{ fontSize: '0.72rem', color: 'var(--muted)' }}>Click star to pin</span>
          </div>

          <div style={{ display: 'grid', gap: '8px' }}>
            {data.sections[0].items.map(ws => {
              const isPinned = pinnedWorkspaces.includes(ws.label)
              return (
                <div key={ws.label} style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '8px 12px',
                  borderRadius: '8px',
                  background: isPinned ? 'rgba(87, 212, 255, 0.04)' : 'rgba(255,255,255,0.01)',
                  border: isPinned ? '1px solid rgba(87, 212, 255, 0.15)' : '1px solid var(--border)'
                }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                    <button 
                      onClick={() => togglePin(ws.label)} 
                      style={{ background: 'none', border: 'none', color: isPinned ? 'var(--warning)' : 'var(--muted)', cursor: 'pointer', fontSize: '1rem', padding: 0 }}
                    >
                      ★
                    </button>
                    <div>
                      <strong style={{ fontSize: '0.85rem' }}>{ws.label}</strong>
                      <div style={{ fontSize: '0.7rem', color: 'var(--muted)' }}>branch: main • coverage 92%</div>
                    </div>
                  </div>
                  <Badge tone={ws.tone ?? 'neutral'}>{ws.value}</Badge>
                </div>
              )
            })}
          </div>
        </Card>

      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))', gap: '16px' }}>
        
        {/* Security Session widget */}
        <Card style={{ padding: '14px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
              <Shield size={14} style={{ color: 'var(--success)' }} />
              <h3 style={{ margin: 0, fontSize: '0.95rem' }}>Access Security Widget</h3>
            </div>
            <span style={{ fontSize: '0.7rem', color: 'var(--muted)' }}>Last Login: 19:45:12</span>
          </div>

          <div style={{ display: 'grid', gap: '8px' }}>
            {sessions.map(s => (
              <div key={s.id} className="security-session-row">
                <div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                    <Laptop size={12} style={{ color: 'var(--muted)' }} />
                    <span style={{ fontWeight: 600 }}>{s.device}</span>
                  </div>
                  <div style={{ fontSize: '0.72rem', color: 'var(--muted)', marginTop: '2px' }}>
                    IP: {s.ip} • Geolocation: {s.geo}
                  </div>
                </div>
                <div>
                  {s.lastActive === 'Current Session' ? (
                    <Badge tone="success">Active</Badge>
                  ) : (
                    <button 
                      onClick={() => revokeSession(s.id)}
                      style={{
                        background: 'rgba(240, 138, 138, 0.1)',
                        border: '1px solid rgba(240, 138, 138, 0.2)',
                        color: 'var(--danger)',
                        borderRadius: '4px',
                        padding: '2px 6px',
                        fontSize: '0.72rem',
                        cursor: 'pointer'
                      }}
                    >
                      Revoke
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        </Card>

        {/* AI status and token usage tracker */}
        <Card style={{ padding: '14px' }}>
          <h3 style={{ margin: '0 0 10px', fontSize: '0.95rem' }}>AI Copilot Costs & Usage</h3>
          <div style={{ display: 'grid', gap: '8px', fontSize: '0.8rem', color: 'var(--muted-strong)' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', padding: '6px 0', borderBottom: '1px solid var(--border)' }}>
              <span>Total Prompt Tokens:</span>
              <strong>142,852 tokens</strong>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', padding: '6px 0', borderBottom: '1px solid var(--border)' }}>
              <span>Completion Tokens:</span>
              <strong>38,591 tokens</strong>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', padding: '6px 0', borderBottom: '1px solid var(--border)' }}>
              <span>Estimated Cost today:</span>
              <strong style={{ color: 'var(--success)' }}>$0.00 (Local LLM)</strong>
            </div>
            <div style={{ display: 'flex', justifyContent: 'space-between', padding: '6px 0' }}>
              <span>Context Cache Efficiency:</span>
              <strong>94.2%</strong>
            </div>
          </div>
        </Card>

      </div>

      {/* 3. COLLAPSIBLE SYSTEM RESOURCE METRICS */}
      <section className="section-block">
        <Card style={{ padding: '10px 14px' }}>
          <button 
            onClick={() => setShowResourceMetrics(!showResourceMetrics)}
            style={{
              width: '100%',
              background: 'none',
              border: 'none',
              color: 'var(--text)',
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              fontWeight: 600,
              fontSize: '0.85rem',
              cursor: 'pointer',
              padding: 0
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
              <Activity size={14} style={{ color: 'var(--accent)' }} />
              <span>Fleet Resource Metrics (CPU, RAM, Disk)</span>
            </div>
            <span>{showResourceMetrics ? 'Hide Metrics' : 'Expand Metrics'}</span>
          </button>

          {showResourceMetrics && (
            <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '10px', marginTop: '12px' }}>
              {data.metrics.map(metric => (
                <div key={metric.label} style={{
                  padding: '10px',
                  borderRadius: '8px',
                  background: 'rgba(255,255,255,0.01)',
                  border: '1px solid var(--border)'
                }}>
                  <div style={{ fontSize: '0.72rem', color: 'var(--muted)', textTransform: 'uppercase' }}>{metric.label}</div>
                  <div style={{ fontSize: '1.2rem', fontWeight: 'bold', margin: '4px 0' }}>{metric.value}</div>
                  <div style={{ fontSize: '0.72rem', color: 'var(--muted-strong)' }}>{metric.detail}</div>
                </div>
              ))}
            </div>
          )}
        </Card>
      </section>

      {/* 4. NOTIFICATIONS & ATTENTION ITEMS LOG */}
      <Card style={{ padding: '14px' }}>
        <h3 style={{ margin: '0 0 10px', fontSize: '0.95rem' }}>Incidents & Attention Log</h3>
        <div style={{ display: 'grid', gap: '8px' }}>
          {data.activity.slice(0, 3).map((item, idx) => (
            <div key={item.title} style={{ 
              display: 'flex', 
              justifyContent: 'space-between', 
              alignItems: 'center',
              padding: '8px 12px',
              borderRadius: '8px',
              background: 'rgba(255, 255, 255, 0.01)',
              border: '1px solid var(--border)'
            }}>
              <div>
                <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                  <span style={{
                    width: '6px',
                    height: '6px',
                    borderRadius: '50%',
                    background: item.tone === 'danger' ? 'var(--danger)' : item.tone === 'warning' ? 'var(--warning)' : 'var(--success)'
                  }} />
                  <strong style={{ fontSize: '0.82rem' }}>{item.title}</strong>
                </div>
                <div style={{ fontSize: '0.72rem', color: 'var(--muted)', marginTop: '2px' }}>
                  {item.detail}
                </div>
              </div>
              <span style={{ fontSize: '0.72rem', color: 'var(--muted)', fontFamily: 'monospace' }}>
                TR-{(1000 + idx)}
              </span>
            </div>
          ))}
        </div>
      </Card>

      {/* 5. SUPER ADMIN ORGANIZATIONS LIST */}
      {isSuperAdmin && (
        <Card style={{ padding: '14px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
            <h3 style={{ margin: 0, fontSize: '0.95rem', display: 'flex', alignItems: 'center', gap: '6px' }}>
              <Layers size={14} style={{ color: 'var(--accent)' }} />
              <span>Platform Tenant Index (Super Admin view)</span>
            </h3>
            <Badge tone="accent">Total: {systemOrgs.length}</Badge>
          </div>
          
          {loadingOrgs ? (
            <div className="skeleton-group" style={{ padding: '10px' }}>
              <div className="skeleton-line" style={{ height: '24px' }} />
              <div className="skeleton-line" style={{ height: '24px' }} />
            </div>
          ) : (
            <div className="table-wrap">
              <table className="data-table" style={{ fontSize: '0.8rem', width: '100%' }}>
                <thead>
                  <tr>
                    <th style={{ textAlign: 'left', padding: '8px' }}>Organization Name</th>
                    <th style={{ textAlign: 'left', padding: '8px' }}>Slug / ID</th>
                    <th style={{ textAlign: 'left', padding: '8px' }}>Registration Date</th>
                    <th style={{ textAlign: 'left', padding: '8px' }}>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {systemOrgs.map(org => (
                    <tr key={org.id} style={{ borderBottom: '1px solid var(--border)' }}>
                      <td style={{ padding: '8px', fontWeight: '500' }}>{org.name}</td>
                      <td style={{ padding: '8px' }}>
                        <code style={{ fontSize: '0.72rem', color: 'var(--muted)', background: 'rgba(255,255,255,0.03)', padding: '2px 6px', borderRadius: '4px' }}>
                          {org.slug || org.id.substring(0, 8)}
                        </code>
                      </td>
                      <td style={{ padding: '8px', color: 'var(--muted)' }}>
                        {new Date(org.createdAt).toLocaleDateString()}
                      </td>
                      <td style={{ padding: '8px' }}>
                        <Badge tone="success">Active</Badge>
                      </td>
                    </tr>
                  ))}
                  {systemOrgs.length === 0 && (
                    <tr>
                      <td colSpan={4} style={{ textAlign: 'center', color: 'var(--muted)', padding: '20px' }}>
                        No organizations registered.
                      </td>
                    </tr>
                  )}
                </tbody>
              </table>
            </div>
          )}
        </Card>
      )}

    </div>
  )
}
