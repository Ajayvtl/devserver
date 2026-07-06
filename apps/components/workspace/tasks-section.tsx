'use client'

import { useEffect, useState } from 'react'
import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { request, connectSocket, submitCommand } from '@/lib/api/client'

interface TaskRecord {
  id: string
  name: string
  capability: string
  target?: string
  status: 'queued' | 'running' | 'completed' | 'failed' | 'cancelled' | 'rolled_back'
  progress?: number
  error?: string
  createdAt: string
  startedAt?: string
  completedAt?: string
}

export function TasksSection({ workspaceId }: { workspaceId: string }) {
  const [tasks, setTasks] = useState<TaskRecord[]>([])
  
  useEffect(() => {
    // Initial fetch
    request<TaskRecord[]>('/api/commands')
      .then(data => {
        if (Array.isArray(data)) {
          setTasks(data.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()))
        }
      })
      .catch(err => console.error("Failed to fetch commands", err))

    // Subscribe to events
    const socket = connectSocket<{ payload?: any }>('/ws/events', {
      onMessage: () => {
        // Just refetch on any event to keep it simple and live
        request<TaskRecord[]>('/api/commands')
          .then(data => {
            if (Array.isArray(data)) {
              setTasks(data.sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()))
            }
          })
      }
    })

    return () => socket.close()
  }, [])

  return (
    <div className="ws-section" style={{ maxWidth: '100%' }}>
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Command Center</div>
          <h2 className="ws-section__title">Tasks</h2>
          <p className="ws-section__subtitle">Live execution engine status</p>
        </div>
      </div>

      <div style={{ display: 'grid', gap: '16px' }}>
        {tasks.length === 0 ? (
          <div style={{ padding: '40px', textAlign: 'center', color: 'var(--muted)', background: 'rgba(255,255,255,0.02)', borderRadius: '12px' }}>
            No tasks have been executed yet.
          </div>
        ) : (
          tasks.map(task => (
            <Card key={task.id} style={{ display: 'flex', flexDirection: 'column', padding: '16px' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', marginBottom: '12px' }}>
                <div>
                  <h3 style={{ margin: '0 0 4px', fontSize: '1.1rem' }}>{task.name || task.capability}</h3>
                  <div style={{ fontSize: '0.85rem', color: 'var(--muted)' }}>
                    Capability: {task.capability} {task.target && `• Target: ${task.target}`}
                  </div>
                </div>
                <Badge tone={
                  task.status === 'completed' ? 'success' : 
                  task.status === 'running' ? 'info' : 
                  task.status === 'failed' ? 'danger' : 'neutral'
                }>
                  {task.status.toUpperCase()}
                </Badge>
              </div>

              {(task.status === 'running' || task.status === 'queued') && (
                <div style={{ marginBottom: '16px' }}>
                  <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.8rem', marginBottom: '4px', color: 'var(--muted-strong)' }}>
                    <span>Progress</span>
                    <span>{task.progress || 0}%</span>
                  </div>
                  <div style={{ height: '6px', background: 'rgba(255,255,255,0.1)', borderRadius: '3px', overflow: 'hidden' }}>
                    <div style={{ height: '100%', width: `${task.progress || 0}%`, background: 'var(--accent)', transition: 'width 0.3s ease' }} />
                  </div>
                </div>
              )}

              {task.error && (
                <div style={{ padding: '8px 12px', background: 'rgba(240, 138, 138, 0.1)', border: '1px solid rgba(240, 138, 138, 0.2)', borderRadius: '8px', color: 'var(--danger)', fontSize: '0.85rem', marginBottom: '16px' }}>
                  <strong>Error: </strong> {task.error}
                </div>
              )}

              <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))', gap: '8px', borderTop: '1px solid rgba(148, 163, 184, 0.08)', paddingTop: '12px', marginTop: 'auto' }}>
                <div style={{ fontSize: '0.8rem', color: 'var(--muted)' }}>
                  <div style={{ marginBottom: '2px' }}>Started At</div>
                  <div style={{ color: 'var(--text)' }}>{task.startedAt ? new Date(task.startedAt).toLocaleTimeString() : '-'}</div>
                </div>
                {task.completedAt && (
                  <div style={{ fontSize: '0.8rem', color: 'var(--muted)' }}>
                    <div style={{ marginBottom: '2px' }}>Completed At</div>
                    <div style={{ color: 'var(--text)' }}>{new Date(task.completedAt).toLocaleTimeString()}</div>
                  </div>
                )}
                
                <div style={{ display: 'flex', gap: '8px', justifyContent: 'flex-end', alignItems: 'center' }}>
                  {(task.status === 'running' || task.status === 'queued') && (
                    <button className="ws-btn" style={{ background: 'rgba(255, 255, 255, 0.05)', color: 'var(--text)', border: '1px solid rgba(255,255,255,0.1)', padding: '6px 12px', borderRadius: '6px', fontSize: '0.8rem', cursor: 'pointer' }}>
                      Cancel
                    </button>
                  )}
                  {task.status === 'failed' && (
                    <button className="ws-btn" onClick={() => submitCommand({ capability: task.capability, target: task.target, name: task.name })} style={{ background: 'rgba(87, 212, 255, 0.1)', color: 'var(--accent)', border: '1px solid rgba(87, 212, 255, 0.2)', padding: '6px 12px', borderRadius: '6px', fontSize: '0.8rem', cursor: 'pointer' }}>
                      Retry
                    </button>
                  )}
                  <button className="ws-btn" style={{ background: 'rgba(255, 255, 255, 0.05)', color: 'var(--text)', border: '1px solid rgba(255,255,255,0.1)', padding: '6px 12px', borderRadius: '6px', fontSize: '0.8rem', cursor: 'pointer' }}>
                    Logs
                  </button>
                </div>
              </div>
            </Card>
          ))
        )}
      </div>
    </div>
  )
}
