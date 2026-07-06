'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import type { ReactNode } from 'react'

import { navigation } from '@/lib/navigation'
import { useAuth } from '@/components/auth/auth-context'

import { Badge } from './ui/badge'
import { Button } from './ui/button'
import { Input } from './ui/input'

interface Props {
  children: ReactNode
  server: string
  notifications: number
}

export function AppShell({ children, server, notifications }: Props) {
  const pathname = usePathname()
  const { user, organizations, currentOrgId, setCurrentOrgId, signOut, isLoading } = useAuth()

  if (isLoading) {
    return <div className="dashboard-shell"><div className="dashboard-main" style={{ padding: '2rem' }}>Loading authentication...</div></div>
  }

  return (
    <div className="dashboard-shell">
      <aside className="dashboard-sidebar">
        <div className="brand-lockup">
          <div className="brand-lockup__mark">D</div>
          <div>
            <div className="brand-lockup__title">DevServer</div>
            <div className="brand-lockup__subtitle">Control plane</div>
          </div>
        </div>

        <nav className="dashboard-nav" aria-label="Sidebar">
          {navigation.map((group) => (
            <div key={group.title} className="dashboard-nav__group">
              <div className="dashboard-nav__title">{group.title}</div>
              {group.items.map((item) => {
                const current = isActive(pathname, item.href)
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={`dashboard-nav__item${current ? ' is-active' : ''}`}
                    aria-current={current ? 'page' : undefined}
                  >
                    <span>{item.label}</span>
                    <small>{item.subtitle}</small>
                  </Link>
                )
              })}
            </div>
          ))}
        </nav>
      </aside>

      <div className="dashboard-main">
        <header className="dashboard-topbar">
          <div className="dashboard-topbar__search">
            <Input placeholder="Search projects, tasks, docs..." aria-label="Search" />
          </div>
          <div className="dashboard-topbar__actions" style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            {organizations.length > 0 && (
              <select
                value={currentOrgId || ''}
                onChange={(e) => setCurrentOrgId(e.target.value)}
                style={{ padding: '0.375rem', borderRadius: '4px', border: '1px solid #ccc', backgroundColor: 'var(--bg-card, #fff)' }}
              >
                {organizations.map((org) => (
                  <option key={org.id} value={org.id}>
                    {org.name}
                  </option>
                ))}
              </select>
            )}
            
            <Badge tone="accent">{notifications} notifications</Badge>
            <Badge tone="neutral">{server}</Badge>
            <span style={{ fontSize: '0.875rem', fontWeight: 500 }}>{user?.id ? 'Active' : ''}</span>
            <Button variant="secondary" onClick={signOut}>
              Sign out
            </Button>
          </div>
        </header>

        <main className="dashboard-content">{children}</main>
      </div>
    </div>
  )
}

function isActive(pathname: string, href: string) {
  const base = href.split('#')[0]
  if (!base) {
    return false
  }
  if (base === '/') {
    return pathname === '/'
  }
  return pathname === base || pathname.startsWith(`${base}/`)
}
