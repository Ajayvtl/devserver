'use client'

import { useState } from 'react'
import { Badge } from '../ui/badge'
import { Input } from '../ui/input'
import { submitCommand } from '@/lib/api/client'

export function KnowledgeSection({
  data,
  onJumpToDefinition,
}: {
  data: {
    files: string[]
    topics: string[]
    symbols?: {
      name: string
      kind: string
      package?: string
      exported?: boolean
      position?: {
        file: string
        line: number
      }
    }[]
  }
  onJumpToDefinition?: (file: string, line: number) => void
}) {
  const [search, setSearch] = useState('')
  const [activeSymbol, setActiveSymbol] = useState<any>(null)
  const [activeTab, setActiveTab] = useState<'references' | 'hierarchy' | 'graph'>('references')

  const kinds = Array.from(new Set(data.symbols?.map((s) => s.kind) || [])).sort()
  const pkgs = Array.from(new Set(data.symbols?.map((s) => s.package).filter(Boolean) || [])).sort()

  const filteredSymbols = data.symbols?.filter((s) => {
    if (search && !s.name.toLowerCase().includes(search.toLowerCase())) return false
    return true
  }) || []

  // Outline of symbols in the same file as the active symbol
  const fileOutline = activeSymbol
    ? data.symbols?.filter(s => s.position?.file === activeSymbol.position?.file) || []
    : []

  const triggerTask = async (capability: string, target: string, name: string) => {
    try {
      await submitCommand({ capability, target, name: `${name} for ${target}` })
      alert(`Queued Task: ${name} for ${target}`)
    } catch (e) {
      console.error(e)
    }
  }

  return (
    <div className="ws-section" style={{ maxWidth: '100%', height: 'calc(100vh - 120px)', display: 'flex', flexDirection: 'column' }}>
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Knowledge Engine</div>
          <h2 className="ws-section__title">Code Intelligence</h2>
          <p className="ws-section__subtitle">Explore {data.symbols?.length || 0} indexed symbols, references, and dependencies.</p>
        </div>
        <button className="ws-btn" onClick={() => triggerTask('workspace.analyze', 'repo', 'Analyze Repository')} style={{ background: 'rgba(87, 212, 255, 0.1)', color: 'var(--accent)', border: '1px solid rgba(87,212,255,0.2)', padding: '8px 16px', borderRadius: '8px', cursor: 'pointer', fontWeight: 600 }}>
          ✦ Analyze Repository
        </button>
      </div>

      <div style={{
        display: 'grid',
        gridTemplateColumns: '320px 1fr 340px',
        gap: '1px',
        background: 'rgba(148, 163, 184, 0.1)',
        border: '1px solid rgba(148, 163, 184, 0.15)',
        borderRadius: '12px',
        flex: 1,
        overflow: 'hidden'
      }}>
        
        {/* Left Pane: Search & Symbol List */}
        <div style={{ background: 'rgba(5, 10, 19, 0.95)', display: 'flex', flexDirection: 'column' }}>
          <div style={{ padding: '12px 14px', borderBottom: '1px solid rgba(148, 163, 184, 0.08)' }}>
            <Input
              placeholder="Search symbols (Ctrl+T)..."
              value={search}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setSearch(e.target.value)}
              style={{ background: 'rgba(255,255,255,0.02)', border: 'none', height: '32px' }}
            />
          </div>
          <div style={{ flex: 1, overflowY: 'auto', padding: '12px' }}>
            {filteredSymbols.length === 0 ? (
              <div style={{ padding: '40px', textAlign: 'center', color: 'var(--muted)' }}>No symbols found</div>
            ) : (
              <div style={{ display: 'grid', gap: '2px' }}>
                {filteredSymbols.map(s => {
                  const isActive = activeSymbol === s
                  return (
                    <div 
                      key={`${s.position?.file}:${s.position?.line}:${s.name}`}
                      onClick={() => setActiveSymbol(s)}
                      style={{
                        display: 'flex', alignItems: 'center', gap: '8px', padding: '8px 10px',
                        borderRadius: '6px', cursor: 'pointer',
                        background: isActive ? 'rgba(87, 212, 255, 0.1)' : 'transparent',
                        border: '1px solid transparent',
                        borderColor: isActive ? 'rgba(87, 212, 255, 0.2)' : 'transparent',
                        color: isActive ? 'var(--text)' : 'var(--muted-strong)',
                      }}
                    >
                      <span style={{ fontSize: '0.85rem', color: 'var(--accent)', opacity: 0.8, minWidth: '40px' }}>{s.kind}</span>
                      <span style={{ fontFamily: 'var(--font-mono)', fontSize: '0.85rem', flex: 1, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                        {s.name}
                      </span>
                      {s.exported && <span style={{ width: '6px', height: '6px', borderRadius: '50%', background: 'var(--success)' }} title="Exported" />}
                    </div>
                  )
                })}
              </div>
            )}
          </div>
        </div>

        {/* Center Pane: Intelligence Details */}
        <div style={{ background: 'rgba(7, 14, 25, 0.95)', display: 'flex', flexDirection: 'column' }}>
          {!activeSymbol ? (
             <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--muted)', fontSize: '0.9rem' }}>
               Select a symbol to explore references and hierarchy.
             </div>
          ) : (
            <>
              {/* Symbol Header */}
              <div style={{ padding: '20px', borderBottom: '1px solid rgba(148, 163, 184, 0.08)' }}>
                <div style={{ display: 'flex', gap: '8px', marginBottom: '8px' }}>
                  <Badge tone="accent">{activeSymbol.kind}</Badge>
                  {activeSymbol.package && <Badge tone="info">{activeSymbol.package}</Badge>}
                  {activeSymbol.exported && <Badge tone="success">Exported</Badge>}
                </div>
                <h3 style={{ fontSize: '1.4rem', margin: '0 0 8px', fontFamily: 'var(--font-mono)' }}>{activeSymbol.name}</h3>
                <div style={{ color: 'var(--muted)', fontSize: '0.85rem', display: 'flex', gap: '16px' }}>
                  <span>{activeSymbol.position?.file} : {activeSymbol.position?.line}</span>
                  <a href="#" onClick={(e) => { 
                    e.preventDefault(); 
                    if (onJumpToDefinition && activeSymbol.position) {
                      onJumpToDefinition(activeSymbol.position.file, activeSymbol.position.line)
                    } else {
                      triggerTask('code.jump', activeSymbol.name, 'Jump to Definition') 
                    }
                  }} style={{ color: 'var(--accent)', textDecoration: 'none' }}>
                    Jump to Definition ↗
                  </a>
                </div>
              </div>

              {/* Tabs */}
              <div style={{ display: 'flex', padding: '0 20px', borderBottom: '1px solid rgba(148, 163, 184, 0.08)', gap: '20px' }}>
                {(['references', 'hierarchy', 'graph'] as const).map(tab => (
                  <div
                    key={tab}
                    onClick={() => setActiveTab(tab)}
                    style={{
                      padding: '12px 0',
                      cursor: 'pointer',
                      fontSize: '0.9rem',
                      fontWeight: activeTab === tab ? 600 : 400,
                      color: activeTab === tab ? 'var(--text)' : 'var(--muted)',
                      borderBottom: `2px solid ${activeTab === tab ? 'var(--accent)' : 'transparent'}`
                    }}
                  >
                    {tab === 'references' ? 'Find References' : tab === 'hierarchy' ? 'Call Hierarchy' : 'Dependency Graph'}
                  </div>
                ))}
              </div>

              {/* Tab Content */}
              <div style={{ flex: 1, padding: '20px', overflowY: 'auto' }}>
                <div style={{ padding: '24px', background: 'rgba(255,255,255,0.02)', border: '1px dashed rgba(255,255,255,0.1)', borderRadius: '12px', textAlign: 'center' }}>
                  <div style={{ marginBottom: '16px', color: 'var(--muted)' }}>
                    Trigger a background task to compute {activeTab === 'references' ? 'all references' : activeTab === 'hierarchy' ? 'the call hierarchy' : 'the dependency graph'} across the workspace.
                  </div>
                  <button 
                    onClick={() => triggerTask(`code.${activeTab}`, activeSymbol.name, `Compute ${activeTab}`)}
                    style={{ background: 'rgba(255,255,255,0.1)', border: '1px solid rgba(255,255,255,0.1)', color: '#fff', padding: '8px 16px', borderRadius: '8px', cursor: 'pointer' }}
                  >
                    Run {activeTab === 'references' ? 'Find References' : activeTab === 'hierarchy' ? 'Call Hierarchy' : 'Dependency Graph'} Task
                  </button>
                </div>
              </div>
            </>
          )}
        </div>

        {/* Right Pane: File Outline & AI Explanations */}
        <div style={{ background: 'rgba(5, 10, 19, 0.95)', display: 'flex', flexDirection: 'column' }}>
          
          <div style={{ padding: '12px 14px', borderBottom: '1px solid rgba(148, 163, 184, 0.08)', fontSize: '0.8rem', fontWeight: 600, color: 'var(--muted)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            AI Explanation
          </div>
          <div style={{ padding: '16px', borderBottom: '1px solid rgba(148, 163, 184, 0.08)' }}>
            {!activeSymbol ? (
               <div style={{ color: 'var(--muted)', fontSize: '0.85rem' }}>Select a symbol to generate an explanation based on indexed context.</div>
            ) : (
              <div style={{ background: 'linear-gradient(135deg, rgba(87, 212, 255, 0.05), rgba(73, 208, 142, 0.05))', border: '1px solid rgba(87, 212, 255, 0.15)', padding: '12px', borderRadius: '8px' }}>
                <p style={{ margin: '0 0 12px', fontSize: '0.85rem', color: 'var(--muted-strong)', lineHeight: 1.5 }}>
                  This symbol is part of the `{activeSymbol.package || 'core'}` package. AI can synthesize its usages to explain its role in the system.
                </p>
                <button onClick={() => triggerTask('ai.explain', activeSymbol.name, 'AI Explanation')} style={{ width: '100%', background: 'var(--accent)', border: 'none', color: '#000', fontWeight: 600, padding: '8px', borderRadius: '6px', fontSize: '0.85rem', cursor: 'pointer' }}>
                  Explain Symbol Context
                </button>
              </div>
            )}
          </div>

          <div style={{ padding: '12px 14px', borderBottom: '1px solid rgba(148, 163, 184, 0.08)', fontSize: '0.8rem', fontWeight: 600, color: 'var(--muted)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            File Outline
          </div>
          <div style={{ flex: 1, overflowY: 'auto', padding: '12px' }}>
             {!activeSymbol ? (
               <div style={{ color: 'var(--muted)', fontSize: '0.85rem', textAlign: 'center', marginTop: '20px' }}>No active file.</div>
             ) : (
               <div style={{ display: 'grid', gap: '4px' }}>
                 <div style={{ fontSize: '0.8rem', color: 'var(--muted)', marginBottom: '8px', textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }} title={activeSymbol.position?.file}>
                   {activeSymbol.position?.file.split('/').pop()}
                 </div>
                 {fileOutline.map(s => (
                    <div 
                      key={s.name}
                      onClick={() => setActiveSymbol(s)}
                      style={{ 
                        display: 'flex', alignItems: 'center', gap: '8px', fontSize: '0.85rem', padding: '4px 8px', borderRadius: '4px', cursor: 'pointer',
                        background: activeSymbol.name === s.name ? 'rgba(255,255,255,0.05)' : 'transparent',
                        color: activeSymbol.name === s.name ? 'var(--text)' : 'var(--muted-strong)'
                      }}
                    >
                      <span style={{ color: 'var(--accent)', opacity: 0.8, minWidth: '36px' }}>{s.kind}</span>
                      <span style={{ fontFamily: 'var(--font-mono)' }}>{s.name}</span>
                    </div>
                 ))}
               </div>
             )}
          </div>

        </div>
      </div>
    </div>
  )
}
