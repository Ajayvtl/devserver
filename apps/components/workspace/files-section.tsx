'use client'

import { useState, useMemo } from 'react'
import { Badge } from '../ui/badge'
import { Input } from '../ui/input'
import type { WorkspaceFileInfo, WorkspaceSession } from '@/lib/types'
import { submitCommand } from '@/lib/api/client'

interface Props { 
  files: WorkspaceFileInfo[]; 
  total: number; 
  onOpenFile?: (path: string) => void;
  session?: WorkspaceSession;
  onSessionChange?: (updates: Partial<WorkspaceSession>) => void;
}

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

interface TreeNode {
  name: string
  path: string
  isDirectory: boolean
  children?: TreeNode[]
  file?: WorkspaceFileInfo
}

function buildTree(files: WorkspaceFileInfo[]): TreeNode[] {
  const root: TreeNode = { name: 'root', path: '', isDirectory: true, children: [] }

  for (const f of files) {
    const parts = f.path.split('/')
    let current = root
    let currentPath = ''

    for (let i = 0; i < parts.length; i++) {
      const part = parts[i]
      const isFile = i === parts.length - 1
      currentPath = currentPath ? `${currentPath}/${part}` : part

      let child = current.children!.find(c => c.name === part)
      if (!child) {
        child = {
          name: part,
          path: currentPath,
          isDirectory: !isFile,
          children: isFile ? undefined : [],
          file: isFile ? f : undefined
        }
        current.children!.push(child)
      }
      current = child
    }
  }

  // Sort directories first, then alphabetical
  const sortNode = (node: TreeNode) => {
    if (node.children) {
      node.children.sort((a, b) => {
        if (a.isDirectory && !b.isDirectory) return -1
        if (!a.isDirectory && b.isDirectory) return 1
        return a.name.localeCompare(b.name)
      })
      node.children.forEach(sortNode)
    }
  }
  sortNode(root)

  return root.children || []
}

export function FilesSection({ files, total, onOpenFile, session, onSessionChange }: Props) {
  const [query, setQuery] = useState('')
  const [contextMenu, setContextMenu] = useState<{x: number, y: number, path: string} | null>(null)
  
  // Local state for non-session UI
  const [localExpanded, setLocalExpanded] = useState<string[]>([])
  const expanded = session?.expandedFolders || localExpanded
  
  const toggleFolder = (path: string) => {
    const isExpanded = expanded.includes(path)
    const next = isExpanded ? expanded.filter(p => p !== path) : [...expanded, path]
    if (onSessionChange) onSessionChange({ expandedFolders: next })
    else setLocalExpanded(next)
  }

  const selectedFiles = session?.selectedFiles || []
  const selectFile = (path: string, multi: boolean = false) => {
    if (multi) {
      const next = selectedFiles.includes(path) ? selectedFiles.filter(p => p !== path) : [...selectedFiles, path]
      if (onSessionChange) onSessionChange({ selectedFiles: next })
    } else {
      if (onSessionChange) onSessionChange({ selectedFiles: [path] })
    }
  }

  const filtered = query
    ? files.filter((f) => f.path.toLowerCase().includes(query.toLowerCase()))
    : files

  const tree = useMemo(() => buildTree(filtered), [filtered])

  const handleAction = async (action: string, target?: string) => {
    if (action === 'reindex') {
      try {
        await submitCommand({ capability: 'workspace.index', name: 'Reindex workspace' })
      } catch (e) {
        console.error(e)
      }
      return
    }
    if (action === 'open_editor' && target && onOpenFile) {
      onOpenFile(target)
      setContextMenu(null)
      return
    }
    
    // Simulate tasks for context menu actions
    try {
      if (action !== 'open_editor') {
        await submitCommand({ capability: `file.${action}`, target: target || 'workspace', name: `File: ${action}` })
      }
    } catch (e) {
      console.error(e)
    }
    setContextMenu(null)
  }

  const handleContextMenu = (e: React.MouseEvent, path: string) => {
    e.preventDefault()
    setContextMenu({ x: e.clientX, y: e.clientY, path })
    selectFile(path, false)
  }

  // A basic mock of git status for visual flair
  const getGitStatus = (path: string) => {
    if (path.endsWith('.go')) return { label: 'M', color: '#eab308' } // Modified
    if (path.endsWith('.tsx')) return { label: 'U', color: '#22c55e' } // Untracked
    if (path.includes('delete')) return { label: 'D', color: '#ef4444' } // Deleted
    return null
  }

  const renderTree = (nodes: TreeNode[], depth: number = 0) => {
    return nodes.map(node => {
      const isExpanded = expanded.includes(node.path)
      const isSelected = selectedFiles.includes(node.path)
      const git = getGitStatus(node.path)

      if (node.isDirectory) {
        return (
          <div key={node.path}>
            <div 
              onClick={() => toggleFolder(node.path)}
              onContextMenu={(e) => handleContextMenu(e, node.path)}
              style={{
                display: 'flex', alignItems: 'center', gap: '6px',
                padding: `4px 8px 4px ${8 + depth * 12}px`,
                cursor: 'pointer',
                background: isSelected ? 'rgba(255,255,255,0.08)' : 'transparent',
                color: 'var(--text)',
                fontSize: '0.85rem',
                userSelect: 'none'
              }}
            >
              <span style={{ fontSize: '0.75rem', width: '12px', display: 'inline-block', opacity: 0.7 }}>
                {isExpanded ? '▼' : '▶'}
              </span>
              <span style={{ color: 'var(--accent)', opacity: 0.9 }}>📂</span>
              <span>{node.name}</span>
            </div>
            {isExpanded && node.children && (
              <div>{renderTree(node.children, depth + 1)}</div>
            )}
          </div>
        )
      }

      return (
        <div 
          key={node.path}
          onClick={(e) => selectFile(node.path, e.metaKey || e.ctrlKey)}
          onDoubleClick={() => handleAction('open_editor', node.path)}
          onContextMenu={(e) => handleContextMenu(e, node.path)}
          style={{
            display: 'flex', alignItems: 'center', justifyContent: 'space-between',
            padding: `4px 8px 4px ${26 + depth * 12}px`,
            cursor: 'pointer',
            background: isSelected ? 'rgba(87, 212, 255, 0.15)' : 'transparent',
            color: isSelected ? 'var(--text)' : 'var(--muted-strong)',
            fontSize: '0.85rem',
            userSelect: 'none',
            border: '1px solid transparent',
            borderColor: isSelected ? 'rgba(87, 212, 255, 0.2)' : 'transparent',
          }}
        >
          <div style={{ display: 'flex', alignItems: 'center', gap: '6px', overflow: 'hidden' }}>
            <span style={{ opacity: 0.6 }}>📄</span>
            <span style={{ whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>{node.name}</span>
          </div>
          {git && (
            <span style={{ color: git.color, fontSize: '0.75rem', fontWeight: 'bold', marginLeft: '8px' }}>
              {git.label}
            </span>
          )}
        </div>
      )
    })
  }

  const activeFile = files.find(f => f.path === selectedFiles[0])

  return (
    <div className="ws-section" style={{ maxWidth: '100%', height: 'calc(100vh - 120px)', display: 'flex', flexDirection: 'column' }} onClick={() => setContextMenu(null)}>
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Workspace Explorer v2</div>
          <h2 className="ws-section__title">Explorer</h2>
          <p className="ws-section__subtitle">Interactive file navigation and code intelligence</p>
        </div>
        <div style={{ display: 'flex', gap: '8px' }}>
          <button className="ws-btn" onClick={() => handleAction('reindex')} style={{ background: 'rgba(87, 212, 255, 0.1)', color: 'var(--accent)', border: '1px solid rgba(87,212,255,0.2)', padding: '6px 12px', borderRadius: '8px', cursor: 'pointer' }}>
            ⟳ Reindex
          </button>
        </div>
      </div>

      <div style={{
        display: 'grid',
        gridTemplateColumns: '260px 1fr',
        gap: '1px',
        background: 'rgba(148, 163, 184, 0.1)',
        border: '1px solid rgba(148, 163, 184, 0.15)',
        borderRadius: '12px',
        flex: 1,
        overflow: 'hidden',
        position: 'relative'
      }}>
        {/* Left Pane: Tree & Search */}
        <div style={{ background: 'rgba(5, 10, 19, 0.95)', display: 'flex', flexDirection: 'column' }}>
          <div style={{ padding: '10px 14px', borderBottom: '1px solid rgba(148, 163, 184, 0.08)' }}>
            <Input
              placeholder="Search files (Ctrl+P)..."
              value={query}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setQuery(e.target.value)}
              style={{ background: 'rgba(255,255,255,0.02)', border: 'none', height: '32px' }}
            />
          </div>
          
          {session && session.openTabs.length > 0 && (
             <div style={{ borderBottom: '1px solid rgba(148, 163, 184, 0.08)' }}>
                <div style={{ padding: '8px 14px', fontSize: '0.75rem', fontWeight: 600, color: 'var(--muted)', textTransform: 'uppercase' }}>Open Editors</div>
                <div style={{ padding: '4px 0' }}>
                   {session.openTabs.map(t => (
                      <div key={t.path} onClick={() => onOpenFile && onOpenFile(t.path)} style={{ padding: '4px 14px 4px 24px', fontSize: '0.8rem', cursor: 'pointer', color: session.activeTab === t.path ? 'var(--accent)' : 'var(--muted)', display: 'flex', gap: '6px' }}>
                        <span style={{ opacity: 0.6 }}>📝</span>
                        {t.path.split('/').pop()}
                        {t.dirty && <span style={{ width: '6px', height: '6px', borderRadius: '50%', background: 'var(--accent)', alignSelf: 'center', marginLeft: 'auto' }} />}
                      </div>
                   ))}
                </div>
             </div>
          )}

          <div style={{ flex: 1, overflowY: 'auto', padding: '8px 0' }}>
             <div style={{ padding: '4px 14px', fontSize: '0.75rem', fontWeight: 600, color: 'var(--muted)', textTransform: 'uppercase', marginBottom: '4px' }}>Workspace</div>
            {filtered.length === 0 ? (
              <div style={{ padding: '20px', textAlign: 'center', color: 'var(--muted)' }}>No files found</div>
            ) : (
              <div className="ws-file-tree" style={{ paddingBottom: '20px' }}>
                {renderTree(tree)}
              </div>
            )}
          </div>
        </div>

        {/* Right Pane: Inspector */}
        <div style={{ background: 'rgba(7, 14, 25, 0.95)', display: 'flex', flexDirection: 'column' }}>
          <div style={{ padding: '12px 14px', borderBottom: '1px solid rgba(148, 163, 184, 0.08)', fontSize: '0.8rem', fontWeight: 600, color: 'var(--muted)', textTransform: 'uppercase', letterSpacing: '0.05em' }}>
            Inspector
          </div>
          
          <div style={{ flex: 1, overflowY: 'auto', padding: '16px' }}>
            {!activeFile ? (
              <div style={{ height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'var(--muted)', fontSize: '0.9rem', textAlign: 'center', padding: '20px' }}>
                Select a file to view properties, git status, and AI insights.
              </div>
            ) : (
              <div style={{ display: 'grid', gap: '20px' }}>
                <div>
                  <h3 style={{ fontSize: '1rem', margin: '0 0 4px', wordBreak: 'break-all' }}>{activeFile.path.split('/').pop()}</h3>
                  <div style={{ color: 'var(--muted)', fontSize: '0.85rem', marginBottom: '12px' }}>{activeFile.path}</div>
                  <div style={{ display: 'flex', gap: '8px' }}>
                    <Badge tone="info">{activeFile.kind}</Badge>
                    <Badge tone="neutral">{formatSize(activeFile.size)}</Badge>
                  </div>
                </div>

                <div style={{ background: 'rgba(255,255,255,0.02)', border: '1px solid rgba(255,255,255,0.05)', padding: '12px', borderRadius: '8px' }}>
                  <div style={{ fontSize: '0.8rem', color: 'var(--muted)', marginBottom: '8px', textTransform: 'uppercase' }}>Quick Actions</div>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: '8px' }}>
                    <button onClick={() => handleAction('open_editor', activeFile.path)} style={{ background: 'rgba(255,255,255,0.1)', border: 'none', color: '#fff', padding: '6px 12px', borderRadius: '6px', fontSize: '0.85rem', cursor: 'pointer' }}>Open Editor</button>
                    <button onClick={() => handleAction('find_references', activeFile.path)} style={{ background: 'rgba(255,255,255,0.05)', border: 'none', color: 'var(--muted-strong)', padding: '6px 12px', borderRadius: '6px', fontSize: '0.85rem', cursor: 'pointer' }}>Find References</button>
                    <button onClick={() => handleAction('git_blame', activeFile.path)} style={{ background: 'rgba(255,255,255,0.05)', border: 'none', color: 'var(--muted-strong)', padding: '6px 12px', borderRadius: '6px', fontSize: '0.85rem', cursor: 'pointer' }}>Git Blame</button>
                  </div>
                </div>

                <div style={{ background: 'rgba(255,255,255,0.02)', border: '1px solid rgba(255,255,255,0.05)', padding: '12px', borderRadius: '8px' }}>
                  <div style={{ fontSize: '0.8rem', color: 'var(--muted)', marginBottom: '8px', textTransform: 'uppercase' }}>File Metadata</div>
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px', fontSize: '0.85rem' }}>
                    <div>
                      <div style={{ color: 'var(--muted)', marginBottom: '2px' }}>Created</div>
                      <div style={{ color: 'var(--text)' }}>2 days ago</div>
                    </div>
                    <div>
                      <div style={{ color: 'var(--muted)', marginBottom: '2px' }}>Modified</div>
                      <div style={{ color: 'var(--text)' }}>Just now</div>
                    </div>
                    <div>
                      <div style={{ color: 'var(--muted)', marginBottom: '2px' }}>Permissions</div>
                      <div style={{ color: 'var(--text)', fontFamily: 'var(--font-mono)' }}>-rw-r--r--</div>
                    </div>
                    <div>
                      <div style={{ color: 'var(--muted)', marginBottom: '2px' }}>Owner</div>
                      <div style={{ color: 'var(--text)' }}>devserver</div>
                    </div>
                  </div>
                </div>

                <div style={{ background: 'linear-gradient(135deg, rgba(87, 212, 255, 0.05), rgba(73, 208, 142, 0.05))', border: '1px solid rgba(87, 212, 255, 0.15)', padding: '12px', borderRadius: '8px' }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '6px', fontSize: '0.8rem', color: 'var(--accent)', marginBottom: '8px', textTransform: 'uppercase', fontWeight: 600 }}>
                    <span>✦</span> AI Intelligence
                  </div>
                  <button onClick={() => handleAction('ai_explain', activeFile.path)} style={{ width: '100%', background: 'var(--accent)', border: 'none', color: '#000', fontWeight: 600, padding: '8px', borderRadius: '6px', fontSize: '0.85rem', cursor: 'pointer' }}>
                    Explain Code
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>
        
        {/* Context Menu */}
        {contextMenu && (
          <div style={{
            position: 'fixed',
            top: contextMenu.y,
            left: contextMenu.x,
            background: 'rgba(15, 23, 42, 0.98)',
            border: '1px solid rgba(148, 163, 184, 0.2)',
            borderRadius: '8px',
            padding: '4px',
            boxShadow: '0 10px 25px rgba(0,0,0,0.5)',
            zIndex: 1000,
            minWidth: '160px',
            backdropFilter: 'blur(10px)'
          }}>
            {[
              { label: 'Open', action: 'open_editor' },
              { label: 'Open to Side', action: 'open_side' },
              { divider: true },
              { label: 'Rename...', action: 'rename' },
              { label: 'Delete', action: 'delete' },
              { label: 'Duplicate', action: 'duplicate' },
              { divider: true },
              { label: 'Copy Path', action: 'copy_path' },
              { label: 'Reveal in OS', action: 'reveal' },
              { divider: true },
              { label: 'AI Explain', action: 'ai_explain' }
            ].map((item, idx) => item.divider ? (
              <div key={idx} style={{ height: '1px', background: 'rgba(148, 163, 184, 0.1)', margin: '4px 0' }} />
            ) : (
              <div 
                key={item.label}
                onClick={(e) => { e.stopPropagation(); handleAction(item.action!, contextMenu.path) }}
                style={{
                  padding: '6px 12px',
                  fontSize: '0.85rem',
                  color: 'var(--text)',
                  cursor: 'pointer',
                  borderRadius: '4px',
                }}
                onMouseEnter={e => e.currentTarget.style.background = 'rgba(87, 212, 255, 0.15)'}
                onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
              >
                {item.label}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
