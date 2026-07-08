'use client'

import { useEffect, useState, useRef, useCallback } from 'react'
import type { Document } from '@/lib/types'

interface Props {
  workspaceId: string
  documents: Document[]
  activeDocPath: string | null
  onTabClick: (path: string) => void
  onCloseTab: (path: string) => void
}

export function WorkspaceEditor({ workspaceId, documents, activeDocPath, onTabClick, onCloseTab }: Props) {
  const [contentCache, setContentCache] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState<Record<string, boolean>>({})
  const contentRef = useRef<HTMLPreElement>(null)

  const activeDoc = documents.find(d => d.path === activeDocPath)
  const activePath = activeDocPath

  const loadContent = useCallback(async (path: string) => {
    const doc = documents.find(d => d.path === path)
    if (!doc) return
    setLoading(prev => ({ ...prev, [path]: true }))
    try {
      const r = await fetch(`http://127.0.0.1:8080/api/workspaces/${workspaceId}/filecontent?path=${encodeURIComponent(path)}`)
      const data = await r.json()
      setContentCache(prev => ({ ...prev, [path]: data.content || '' }))
      setLoading(prev => ({ ...prev, [path]: false }))
      setTimeout(() => scrollToCursor(doc), 50)
    } catch (err) {
      console.error(err)
      setLoading(prev => ({ ...prev, [path]: false }))
    }
  }, [workspaceId, documents])

  const activeContent = activePath ? contentCache[activePath] : undefined

  useEffect(() => {
    if (!activePath) return
    const doc = documents.find(d => d.path === activePath)
    if (!doc) return
    if (activeContent) {
      scrollToCursor(doc)
      return
    }
    loadContent(activePath)
  }, [activePath, loadContent, documents, activeContent])

  // Re-scroll if the cursor changes on the same already-loaded document
  useEffect(() => {
    if (activeDoc && activeContent) {
      scrollToCursor(activeDoc)
    }
  }, [activeDoc?.cursor.line, activeDoc, activeContent])

  const scrollToCursor = (doc: Document) => {
    if (!contentRef.current || doc.cursor.line <= 0) return

    // Find the line element and scroll to it
    const lineElements = contentRef.current.querySelectorAll('.code-line')
    const targetLine = lineElements[doc.cursor.line - 1]

    if (targetLine) {
      targetLine.scrollIntoView({ block: 'center', behavior: 'smooth' })
    }
  }

  const renderContent = () => {
    if (!activeDoc) return (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', height: '100%', color: 'var(--muted)', fontStyle: 'italic' }}>
        Select a file to open
      </div>
    )

    if (loading[activeDoc.path]) return (
      <div style={{ padding: '20px', color: 'var(--muted)' }}>Loading {activeDoc.path}...</div>
    )

    const content = contentCache[activeDoc.path] || ''
    const lines = content.split('\n')

    return (
      <pre ref={contentRef} style={{ margin: 0, padding: '16px 0', fontFamily: 'var(--font-mono)', fontSize: '13px', lineHeight: 1.5, overflowX: 'auto' }}>
        {lines.map((line, idx) => {
          const lineNum = idx + 1
          const isTargetLine = activeDoc.cursor.line === lineNum
          return (
            <div
              key={idx}
              className="code-line"
              style={{
                display: 'flex',
                background: isTargetLine ? 'rgba(87, 212, 255, 0.15)' : 'transparent',
                borderLeft: isTargetLine ? '3px solid var(--accent)' : '3px solid transparent',
              }}
            >
              <div style={{ width: '40px', flexShrink: 0, textAlign: 'right', paddingRight: '12px', color: isTargetLine ? 'var(--text)' : 'var(--muted)', userSelect: 'none' }}>
                {lineNum}
              </div>
              <div style={{ whiteSpace: 'pre', color: 'var(--text)' }}>
                {line}
              </div>
            </div>
          )
        })}
      </pre>
    )
  }

  return (
    <div style={{ display: 'flex', flexDirection: 'column', height: '100%', background: 'rgba(5, 10, 19, 0.95)', border: '1px solid rgba(148, 163, 184, 0.1)', borderRadius: '12px', overflow: 'hidden' }}>
      {/* Tabs */}
      <div style={{ display: 'flex', background: 'rgba(255,255,255,0.02)', borderBottom: '1px solid rgba(148, 163, 184, 0.1)', overflowX: 'auto' }}>
        {documents.map(doc => {
          const isActive = doc.path === activeDocPath
          const filename = doc.path.split('/').pop() || doc.path
          return (
            <div
              key={doc.path}
              onClick={() => onTabClick(doc.path)}
              style={{
                display: 'flex', alignItems: 'center', gap: '8px', padding: '10px 16px',
                cursor: 'pointer', fontSize: '0.85rem',
                background: isActive ? 'transparent' : 'rgba(0,0,0,0.2)',
                color: isActive ? 'var(--text)' : 'var(--muted)',
                borderBottom: isActive ? '2px solid var(--accent)' : '2px solid transparent',
                borderRight: '1px solid rgba(148, 163, 184, 0.05)',
              }}
            >
              <span style={{ fontFamily: 'var(--font-mono)' }}>{filename}</span>
              <button
                onClick={(e) => { e.stopPropagation(); onCloseTab(doc.path) }}
                style={{ background: 'transparent', border: 'none', color: 'var(--muted)', cursor: 'pointer', padding: '2px' }}
                title="Close"
              >×</button>
            </div>
          )
        })}
      </div>

      {/* Breadcrumbs */}
      {activeDoc && (
        <div style={{ display: 'flex', padding: '8px 16px', fontSize: '0.8rem', color: 'var(--muted)', borderBottom: '1px solid rgba(148, 163, 184, 0.05)' }}>
          {activeDoc.path.replace(/\//g, ' › ')}
        </div>
      )}

      {/* Editor Content */}
      <div style={{ flex: 1, overflowY: 'auto' }}>
        {renderContent()}
      </div>
    </div>
  )
}
