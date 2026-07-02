'use client'

import type { ReactNode } from 'react'

import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Input } from '../ui/input'
import { navigation } from '@/lib/navigation'

interface Props {
  children: ReactNode
  server: string
  notifications: number
}

export function DashboardShell({ children, server, notifications }: Props) {
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
              {group.items.map((item) => (
                <a key={item.href} href={item.href} className="dashboard-nav__item">
                  <span>{item.label}</span>
                  <small>{item.subtitle}</small>
                </a>
              ))}
            </div>
          ))}
        </nav>
      </aside>

      <div className="dashboard-main">
        <header className="dashboard-topbar">
          <div className="dashboard-topbar__search">
            <Input placeholder="Search projects, services, tasks..." />
          </div>
          <div className="dashboard-topbar__actions">
            <Badge tone="accent">{notifications} notifications</Badge>
            <Badge tone="neutral">{server}</Badge>
            <Button variant="secondary" href="#settings">
              Profile
            </Button>
          </div>
        </header>

        <main className="dashboard-content">{children}</main>
      </div>
    </div>
  )
}
