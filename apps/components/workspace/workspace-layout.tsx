'use client'

import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import type { ReactNode } from 'react'

import type { WorkspaceSection } from '@/lib/types'

const sections: Array<{ key: WorkspaceSection; label: string; icon: string }> = [
  { key: 'overview', label: 'Overview', icon: '◎' },
  { key: 'files', label: 'Files', icon: '⊞' },
  { key: 'repository', label: 'Repository', icon: '⎇' },
  { key: 'environment', label: 'Environment', icon: '⚙' },
  { key: 'infrastructure', label: 'Infrastructure', icon: '⬡' },
  { key: 'services', label: 'Services', icon: '◈' },
  { key: 'deployments', label: 'Deployments', icon: '▲' },
  { key: 'database', label: 'Database', icon: '⊟' },
  { key: 'domains', label: 'Domains', icon: '◉' },
  { key: 'logs', label: 'Logs', icon: '≡' },
  { key: 'ai', label: 'AI', icon: '✦' },
  { key: 'knowledge', label: 'Knowledge', icon: '◆' },
  { key: 'doctor', label: 'Doctor', icon: '♥' },
  { key: 'settings', label: 'Settings', icon: '⊕' },
]

interface Props {
  id: string
  children: ReactNode
  active: WorkspaceSection
}

export function WorkspaceLayout({ id, children, active }: Props) {
  return (
    <div className="ws-layout">
      <aside className="ws-sidebar">
        <div className="ws-sidebar__header">
          <div className="ws-sidebar__icon">W</div>
          <div>
            <div className="ws-sidebar__name">Workspace</div>
            <div className="ws-sidebar__id">{id}</div>
          </div>
        </div>
        <nav className="ws-nav" aria-label="Workspace sections">
          {sections.map((s) => (
            <Link
              key={s.key}
              href={`/workspace/${id}?section=${s.key}`}
              className={`ws-nav__item${active === s.key ? ' is-active' : ''}`}
              aria-current={active === s.key ? 'page' : undefined}
            >
              <span className="ws-nav__icon">{s.icon}</span>
              <span className="ws-nav__label">{s.label}</span>
            </Link>
          ))}
        </nav>
      </aside>
      <main className="ws-main">{children}</main>
    </div>
  )
}
