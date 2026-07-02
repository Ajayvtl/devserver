'use client'

import { useState } from 'react'
import { Badge } from '../ui/badge'
import { Input } from '../ui/input'
import { EmptyState } from '../ui/empty-state'
import type { WorkspaceFileInfo } from '@/lib/types'

interface Props { files: WorkspaceFileInfo[]; total: number }

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

function buildTree(files: WorkspaceFileInfo[]) {
  const tree: Record<string, WorkspaceFileInfo[]> = {}
  for (const f of files) {
    const parts = f.path.split('/')
    const dir = parts.length > 1 ? parts.slice(0, -1).join('/') : '.'
    if (!tree[dir]) tree[dir] = []
    tree[dir].push(f)
  }
  return tree
}

export function FilesSection({ files, total }: Props) {
  const [query, setQuery] = useState('')
  const filtered = query
    ? files.filter((f) => f.path.toLowerCase().includes(query.toLowerCase()))
    : files
  const tree = buildTree(filtered)
  const dirs = Object.keys(tree).sort()

  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">File Explorer</div>
          <h2 className="ws-section__title">Files</h2>
          <p className="ws-section__subtitle">{total} indexed files</p>
        </div>
        <Badge tone="info">Read-only</Badge>
      </div>

      <div style={{ maxWidth: 520 }}>
        <Input
          placeholder="Search files..."
          value={query}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setQuery(e.target.value)}
          aria-label="Search files"
        />
      </div>

      {filtered.length === 0 ? (
        <EmptyState title="No files found" description={query ? 'Try a different search term.' : 'No files have been indexed yet.'} />
      ) : (
        <div className="ws-file-tree">
          {dirs.map((dir) => (
            <div key={dir} className="ws-file-group">
              <div className="ws-file-group__header">
                <span className="ws-file-group__icon">📁</span>
                <span className="ws-file-group__name">{dir}</span>
                <Badge tone="neutral">{tree[dir].length}</Badge>
              </div>
              <div className="ws-file-group__files">
                {tree[dir].map((f) => {
                  const name = f.path.split('/').pop() || f.path
                  return (
                    <div key={f.path} className="ws-file-row">
                      <span className="ws-file-row__name">{name}</span>
                      <span className="ws-file-row__size">{formatSize(f.size)}</span>
                      <span className="ws-file-row__kind">{f.kind}</span>
                    </div>
                  )
                })}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
