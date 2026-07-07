'use client'

import { useState } from 'react'
import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import { MetricCard } from '../ui/metric-card'
import { Progress } from '../ui/progress'

import type { DashboardDataV2, KnowledgeArticle } from '@/lib/types'

interface Props {
  data: DashboardDataV2
  knowledge: KnowledgeArticle[]
}

export function DashboardPageContent({ data, knowledge }: Props) {
  // Personalization settings
  const [pinnedWorkspaces, setPinnedWorkspaces] = useState<string[]>(['devserver'])
  const [visibleWidgets, setVisibleWidgets] = useState({
    telemetry: true,
    timeline: true,
    tasks: true,
    personalization: true
  })

  // Mock list of attention today items
  const attentionItems = [
    { id: 'att-1', type: 'danger', text: 'Executor SSH-Remote is unreachable.', action: 'Reconnect' },
    { id: 'att-2', type: 'warning', text: 'Staging environment is missing DB_PASSWORD secret.', action: 'Resolve' }
  ]

  const togglePin = (name: string) => {
    setPinnedWorkspaces(prev => 
      prev.includes(name) ? prev.filter(w => w !== name) : [...prev, name]
    )
  }

  return (
    <div className="dashboard-layout">
      {/* 1. PRODUCT IDENTITY: PRIMARY PLATFORM STATUS DECK */}
      <section className="section-block" style={{ marginBottom: '24px' }}>
        <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: '16px' }}>
          
          {/* Question: Is my platform healthy? */}
          <Card style={{ padding: '20px', borderLeft: '4px solid var(--success)' }}>
            <div className="card__eyebrow">Platform Health</div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '16px', marginTop: '12px' }}>
              <div style={{
                width: '56px',
                height: '56px',
                borderRadius: '50%',
                border: '4px solid var(--success)',
                display: 'grid',
                placeItems: 'center',
                fontWeight: '700',
                fontSize: '1.1rem',
                color: 'var(--success)'
              }}>
                98%
              </div>
              <div>
                <h3 style={{ margin: 0 }}>All Systems Nominal</h3>
                <p style={{ margin: '4px 0 0', fontSize: '0.82rem', color: 'var(--muted)' }}>
                  8 executors connected & online.
                </p>
              </div>
            </div>
          </Card>

          {/* Question: Is AI configured? */}
          <Card style={{ padding: '20px', borderLeft: '4px solid var(--accent)' }}>
            <div className="card__eyebrow">AI Copilot Config</div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '14px', marginTop: '12px' }}>
              <div style={{ fontSize: '1.8rem' }}>🤖</div>
              <div>
                <h3 style={{ margin: 0 }}>Ollama (Localhost)</h3>
                <p style={{ margin: '4px 0 0', fontSize: '0.82rem', color: 'var(--muted)' }}>
                  Model: `codellama` (Latency 45ms)
                </p>
              </div>
            </div>
            <div style={{ marginTop: '10px' }}>
              <Badge tone="success">Active Connection</Badge>
            </div>
          </Card>

          {/* Question: Are deployments blocked? */}
          <Card style={{ padding: '20px', borderLeft: '4px solid var(--info)' }}>
            <div className="card__eyebrow">Release Pipelines</div>
            <div style={{ display: 'flex', alignItems: 'center', gap: '14px', marginTop: '12px' }}>
              <div style={{ fontSize: '1.8rem' }}>🚀</div>
              <div>
                <h3 style={{ margin: 0 }}>No Blocked Deployments</h3>
                <p style={{ margin: '4px 0 0', fontSize: '0.82rem', color: 'var(--muted)' }}>
                  Production deployment completed v1.2.3.
                </p>
              </div>
            </div>
            <div style={{ marginTop: '10px' }}>
              <Badge tone="neutral">0 blockers in queue</Badge>
            </div>
          </Card>
        </div>
      </section>

      {/* 2. PRODUCT IDENTITY: CONTEXTUAL NEXT ACTIONS & ATTENTION TODAY */}
      <section className="section-block" style={{ marginBottom: '24px' }}>
        <div style={{ display: 'grid', gridTemplateColumns: '2fr 1fr', gap: '16px', alignItems: 'stretch' }}>
          
          {/* Action: What should I do next? */}
          <Card style={{ padding: '22px', background: 'linear-gradient(135deg, rgba(87, 212, 255, 0.1), rgba(73, 208, 142, 0.05))', border: '1px solid rgba(87, 212, 255, 0.25)' }}>
            <div className="card__eyebrow" style={{ color: 'var(--accent)' }}>Recommended Next Action</div>
            <h2 style={{ margin: '8px 0 4px', fontSize: '1.4rem', fontFamily: 'var(--font-heading)' }}>
              Open active workspace & sync files
            </h2>
            <p style={{ margin: '0 0 16px', color: 'var(--muted-strong)', fontSize: '0.9rem' }}>
              Workspace `devserver` has 4 uncommitted files on branch `feature/indexer-stabilization`.
            </p>
            <div style={{ display: 'flex', gap: '10px' }}>
              <Button variant="primary" href="/workspace/devserver">
                Launch Workspace IDE
              </Button>
              <Button variant="secondary" href="/projects">
                View All Workspaces
              </Button>
            </div>
          </Card>

          {/* Action: What needs attention today? */}
          <Card style={{ padding: '20px' }}>
            <div className="card__eyebrow" style={{ color: 'var(--danger)' }}>Needs Attention Today</div>
            <div style={{ display: 'grid', gap: '12px', marginTop: '12px' }}>
              {attentionItems.map(item => (
                <div key={item.id} style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'start', gap: '10px', fontSize: '0.85rem' }}>
                  <p style={{ margin: 0, color: 'var(--text)' }}>
                    <span style={{ color: item.type === 'danger' ? 'var(--danger)' : 'var(--warning)', marginRight: '6px' }}>●</span>
                    {item.text}
                  </p>
                  <button style={{
                    background: 'none',
                    border: 'none',
                    color: 'var(--accent)',
                    fontWeight: 'bold',
                    cursor: 'pointer',
                    padding: 0
                  }}>
                    {item.action}
                  </button>
                </div>
              ))}
            </div>
          </Card>
        </div>
      </section>

      {/* 3. CONTROL DECK GRID */}
      <div className="dashboard-grid" style={{ marginBottom: '24px' }}>
        
        {/* Workspace Quick List Card */}
        <Card style={{ padding: '20px' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '14px' }}>
            <h3 style={{ margin: 0 }}>Workspaces Context</h3>
            <span style={{ fontSize: '0.75rem', color: 'var(--muted)' }}>Click star to pin</span>
          </div>

          <div style={{ display: 'grid', gap: '10px' }}>
            {data.sections[0].items.map(ws => {
              const isPinned = pinnedWorkspaces.includes(ws.label)
              return (
                <div key={ws.label} style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '12px 16px',
                  borderRadius: '12px',
                  background: isPinned ? 'rgba(87, 212, 255, 0.06)' : 'rgba(255,255,255,0.02)',
                  border: isPinned ? '1px solid rgba(87, 212, 255, 0.2)' : '1px solid var(--border)'
                }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                    <button 
                      onClick={() => togglePin(ws.label)} 
                      style={{ background: 'none', border: 'none', color: isPinned ? 'var(--warning)' : 'var(--muted)', cursor: 'pointer', fontSize: '1.1rem' }}
                    >
                      {isPinned ? '★' : '☆'}
                    </button>
                    <div>
                      <strong style={{ fontSize: '0.9rem' }}>{ws.label}</strong>
                      <div style={{ fontSize: '0.75rem', color: 'var(--muted)' }}>Branch: main</div>
                    </div>
                  </div>
                  <Badge tone={ws.tone ?? 'neutral'}>{ws.value}</Badge>
                </div>
              )
            })}
          </div>
        </Card>

        {/* Running Jobs & Task Engine Queue */}
        <Card style={{ padding: '20px' }}>
          <h3 style={{ margin: '0 0 14px' }}>Task Execution Queue</h3>
          {data.tasks.length ? (
            <div style={{ display: 'grid', gap: '14px' }}>
              {data.tasks.map((task) => (
                <div key={task.title} style={{ padding: '10px', borderRadius: '10px', background: 'rgba(255, 255, 255, 0.02)', border: '1px solid var(--border)' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
                    <strong style={{ fontSize: '0.85rem' }}>{task.title}</strong>
                    <Badge tone={task.state === 'Done' ? 'success' : task.state === 'Running' ? 'accent' : 'warning'}>
                      {task.state}
                    </Badge>
                  </div>
                  <Progress value={task.progress} tone={task.state === 'Done' ? 'success' : 'accent'} />
                </div>
              ))}
            </div>
          ) : (
            <EmptyState title="No active jobs" description="Task execution queue is currently empty." />
          )}
        </Card>

        {/* Telemetry Metrics Widget */}
        {visibleWidgets.telemetry && (
          <Card style={{ padding: '20px', gridColumn: 'span 2' }}>
            <h3 style={{ margin: '0 0 14px' }}>System Resource Telemetry</h3>
            <div className="metric-grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '14px' }}>
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
          </Card>
        )}
      </div>

      {/* 4. PERSONALIZATION CONTROLS PANEL */}
      {visibleWidgets.personalization && (
        <Card style={{ padding: '20px', marginBottom: '24px' }}>
          <h3 style={{ margin: '0 0 14px' }}>Control Deck Personalization</h3>
          <div style={{ display: 'flex', flexWrap: 'wrap', gap: '18px' }}>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '0.85rem', cursor: 'pointer' }}>
              <input 
                type="checkbox" 
                checked={visibleWidgets.telemetry} 
                onChange={(e) => setVisibleWidgets(prev => ({ ...prev, telemetry: e.target.checked }))}
              />
              Show Telemetry Metrics
            </label>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '0.85rem', cursor: 'pointer' }}>
              <input 
                type="checkbox" 
                checked={visibleWidgets.timeline} 
                onChange={(e) => setVisibleWidgets(prev => ({ ...prev, timeline: e.target.checked }))}
              />
              Show Recent Activity
            </label>
            <label style={{ display: 'flex', alignItems: 'center', gap: '8px', fontSize: '0.85rem', cursor: 'pointer' }}>
              <input 
                type="checkbox" 
                checked={visibleWidgets.tasks} 
                onChange={(e) => setVisibleWidgets(prev => ({ ...prev, tasks: e.target.checked }))}
              />
              Show Task Queue
            </label>
          </div>
        </Card>
      )}

      {/* 5. RECENT ACTIVITY TIMELINE */}
      {visibleWidgets.timeline && (
        <Card style={{ padding: '20px' }}>
          <h3 style={{ margin: '0 0 14px' }}>Audited Activity Timeline</h3>
          <div className="timeline">
            {data.activity.map((item) => (
              <div key={item.title} className="timeline__item">
                <div className={`timeline__dot timeline__dot--${item.tone}`} />
                <div className="timeline__body">
                  <div className="timeline__top" style={{ display: 'flex', justifyContent: 'space-between' }}>
                    <h4 style={{ margin: 0, fontSize: '0.9rem' }}>{item.title}</h4>
                    <span style={{ fontSize: '0.75rem', color: 'var(--muted)' }}>{item.when}</span>
                  </div>
                  <p style={{ margin: '4px 0 0', fontSize: '0.82rem', color: 'var(--muted-strong)' }}>{item.detail}</p>
                </div>
              </div>
            ))}
          </div>
        </Card>
      )}
    </div>
  )
}
