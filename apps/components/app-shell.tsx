'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useState, useEffect, useRef, type ReactNode } from 'react'
import { 
  LayoutDashboard, 
  FolderKanban, 
  Terminal, 
  Activity, 
  History, 
  Database, 
  Lock, 
  Cpu, 
  Wrench, 
  BookOpen, 
  ShieldAlert, 
  Settings,
  Bell,
  Wifi,
  WifiOff,
  Search,
  Menu,
  X
} from 'lucide-react'

import { navigation } from '@/lib/navigation'
import { useAuth } from '@/components/auth/auth-context'
import { usePreferences, type Density, type FontScale } from '@/components/providers'

import { Badge } from './ui/badge'
import { Button } from './ui/button'

interface Props {
  children: ReactNode
  server: string
  notifications?: number
}

function getNavIcon(href: string) {
  switch (href) {
    case '/dashboard': return LayoutDashboard
    case '/projects': return FolderKanban
    case '/workspace/devserver': return Terminal
    case '/monitor': return Activity
    case '/deploy': return History
    case '/backup': return Database
    case '/config/environments': return Lock
    case '/config/providers': return Cpu
    case '/devcenter/project-doctor': return Wrench
    case '/devcenter/architecture': return BookOpen
    case '/devcenter/ai-usage': return ShieldAlert
    case '/settings': return Settings
    default: return FolderKanban
  }
}

export function AppShell({ children, server, notifications = 3 }: Props) {
  const pathname = usePathname()
  const router = useRouter()
  const { user, organizations, currentOrgId, setCurrentOrgId, signOut, isLoading } = useAuth()
  const { density, setDensity, fontScale, setFontScale } = usePreferences()

  // State managers
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)
  const [isCommandPaletteOpen, setIsCommandPaletteOpen] = useState(false)
  const [isNotificationOpen, setIsNotificationOpen] = useState(false)
  const [isOnline, setIsOnline] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')

  // Refs
  const paletteInputRef = useRef<HTMLInputElement>(null)
  const paletteRef = useRef<HTMLDivElement>(null)
  const notificationRef = useRef<HTMLDivElement>(null)

  // System status mock logs
  const systemNotifications = [
    { id: 1, type: 'danger', message: 'Deployment failed on executor SSH-Remote', time: '5m ago' },
    { id: 2, type: 'warning', message: 'Ollama AI provider connection latency > 2.4s', time: '18m ago' },
    { id: 3, type: 'success', message: 'Workspace index complete. Model ready.', time: '1h ago' }
  ]

  // Online / Offline state tracking
  useEffect(() => {
    setIsOnline(navigator.onLine)
    const handleOnline = () => setIsOnline(true)
    const handleOffline = () => setIsOnline(false)

    window.addEventListener('online', handleOnline)
    window.addEventListener('offline', handleOffline)
    return () => {
      window.removeEventListener('online', handleOnline)
      window.removeEventListener('offline', handleOffline)
    }
  }, [])

  // Keyboard listener for command palette (Ctrl+K or Cmd+K)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
        e.preventDefault()
        setIsCommandPaletteOpen((prev) => !prev)
      }
      if (e.key === 'Escape') {
        setIsCommandPaletteOpen(false)
        setIsNotificationOpen(false)
      }
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [])

  // Autofocus input when command palette opens
  useEffect(() => {
    if (isCommandPaletteOpen && paletteInputRef.current) {
      paletteInputRef.current.focus()
    }
  }, [isCommandPaletteOpen])

  // Click outside handlers
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (paletteRef.current && !paletteRef.current.contains(e.target as Node)) {
        setIsCommandPaletteOpen(false)
      }
      if (notificationRef.current && !notificationRef.current.contains(e.target as Node)) {
        setIsNotificationOpen(false)
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  if (isLoading) {
    return (
      <div className="dashboard-shell">
        <div className="dashboard-main" style={{ display: 'grid', placeItems: 'center', height: '100vh', backgroundColor: '#07111d' }}>
          <div style={{ textAlign: 'center' }}>
            <div className="brand-lockup__mark" style={{ margin: '0 auto 1rem', width: '64px', height: '64px', fontSize: '1.5rem' }}>D</div>
            <div style={{ color: 'var(--muted)', fontSize: '0.875rem', letterSpacing: '0.05em' }}>AUTHENTICATING ACCOUNT...</div>
          </div>
        </div>
      </div>
    )
  }

  // Filter shortcuts
  const allShortcuts = navigation.flatMap(g => g.items.map(i => ({ ...i, group: g.title })))
  const filteredShortcuts = allShortcuts.filter(item => 
    item.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.subtitle.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const handleShortcutClick = (href: string) => {
    setIsCommandPaletteOpen(false)
    setSearchQuery('')
    router.push(href)
  }

  return (
    <div className="shell">
      {/* Offline Status Alert Banner */}
      {!isOnline && (
        <div className="offline-banner" role="alert" aria-live="assertive">
          <WifiOff size={16} className="offline-banner__icon" />
          <span className="offline-banner__text">
            <strong>Offline Mode</strong> — Access is restricted to cached read-only data. All operations will resume automatically when connectivity is restored.
          </span>
        </div>
      )}

      {/* Sidebar Navigation - Desktop (Compact IDE style) */}
      <aside className="sidebar sidebar--compact" aria-label="Main Navigation">
        <div className="sidebar__brand">
          <div className="sidebar__mark">D</div>
          <div>
            <div className="sidebar__name">DevServer</div>
            <div className="sidebar__tag">IDE Control Deck</div>
          </div>
        </div>

        {/* Workspace Quick Selector */}
        <div className="sidebar__panel">
          <div className="card__eyebrow" style={{ marginBottom: '6px', fontSize: '0.72rem' }}>Active Context</div>
          <select 
            style={{ 
              width: '100%', 
              padding: '6px 10px', 
              borderRadius: '8px', 
              border: '1px solid var(--border)', 
              background: 'rgba(255,255,255,0.03)', 
              color: 'var(--text)',
              fontSize: '0.8rem',
              fontWeight: '500'
            }}
            defaultValue="devserver"
            onChange={(e) => router.push(`/workspace/${e.target.value}`)}
          >
            <option value="devserver">devserver</option>
            <option value="auth-service">auth-service</option>
            <option value="db-indexer">db-indexer</option>
          </select>
        </div>

        <nav className="sidebar__nav">
          {navigation.map((group) => (
            <div key={group.title}>
              <h4 className="sidebar__group-title">{group.title}</h4>
              <div className="sidebar__links" style={{ gap: '4px' }}>
                {group.items.map((item) => {
                  const current = isActive(pathname, item.href)
                  const Icon = getNavIcon(item.href)
                  return (
                    <Link
                      key={item.href}
                      href={item.href}
                      className={`sidebar__link${current ? ' is-active' : ''}`}
                      aria-current={current ? 'page' : undefined}
                    >
                      <Icon size={14} style={{ opacity: current ? 1 : 0.65 }} />
                      <span className="sidebar__link-label">{item.label}</span>
                    </Link>
                  )
                })}
              </div>
            </div>
          ))}
        </nav>
      </aside>

      {/* Mobile Drawer menu - Slides in */}
      <div className={`mobile-drawer${isMobileMenuOpen ? ' is-open' : ''}`} aria-hidden={!isMobileMenuOpen}>
        <div className="mobile-drawer__overlay" onClick={() => setIsMobileMenuOpen(false)} />
        <div className="mobile-drawer__content">
          <div className="mobile-drawer__header">
            <div className="sidebar__brand">
              <div className="sidebar__mark">D</div>
              <div>
                <div className="sidebar__name">DevServer</div>
                <div className="sidebar__tag">IDE Control Deck</div>
              </div>
            </div>
            <button className="mobile-drawer__close" aria-label="Close menu" onClick={() => setIsMobileMenuOpen(false)}>
              <X size={20} />
            </button>
          </div>
          <nav className="mobile-drawer__nav">
            {navigation.map((group) => (
              <div key={group.title} style={{ marginBottom: '16px' }}>
                <h4 className="sidebar__group-title">{group.title}</h4>
                <div className="sidebar__links" style={{ gap: '4px' }}>
                  {group.items.map((item) => {
                    const current = isActive(pathname, item.href)
                    const Icon = getNavIcon(item.href)
                    return (
                      <Link
                        key={item.href}
                        href={item.href}
                        onClick={() => setIsMobileMenuOpen(false)}
                        className={`sidebar__link${current ? ' is-active' : ''}`}
                      >
                        <Icon size={14} />
                        <span className="sidebar__link-label">{item.label}</span>
                      </Link>
                    )
                  })}
                </div>
              </div>
            ))}
          </nav>
        </div>
      </div>

      {/* Main Workspace Area */}
      <div className="shell__main">
        {/* Universal Topbar */}
        <header className="topbar" style={{ padding: '12px 20px' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
            <button 
              className="mobile-hamburger-btn" 
              aria-label="Toggle Navigation Drawer" 
              onClick={() => setIsMobileMenuOpen(true)}
            >
              <Menu size={20} />
            </button>
            <div>
              <div className="topbar__eyebrow" style={{ fontSize: '0.7rem' }}>Developer Platform Control Deck</div>
              <h2 className="topbar__title" style={{ fontSize: '1.05rem', marginTop: '2px' }}>Overview / Workspace</h2>
            </div>
          </div>

          <div className="topbar__right" style={{ gap: '8px' }}>
            {/* Display Preference Controls (Density & Font Scaling) */}
            <div style={{ display: 'flex', alignItems: 'center', gap: '4px', borderRight: '1px solid var(--border)', paddingRight: '8px' }}>
              <select
                value={density}
                onChange={(e) => setDensity(e.target.value as Density)}
                style={{
                  padding: '4px 8px',
                  borderRadius: '6px',
                  border: '1px solid var(--border)',
                  backgroundColor: 'rgba(15,25,41,0.95)',
                  color: 'var(--text)',
                  fontSize: '0.75rem',
                  fontWeight: 600
                }}
                aria-label="Select Display Density"
              >
                <option value="comfortable">Comfortable</option>
                <option value="compact">Compact</option>
                <option value="ultra">Ultra Compact</option>
              </select>

              <select
                value={fontScale}
                onChange={(e) => setFontScale(Number(e.target.value) as FontScale)}
                style={{
                  padding: '4px 8px',
                  borderRadius: '6px',
                  border: '1px solid var(--border)',
                  backgroundColor: 'rgba(15,25,41,0.95)',
                  color: 'var(--text)',
                  fontSize: '0.75rem',
                  fontWeight: 600
                }}
                aria-label="Select Font Scaling"
              >
                <option value={14}>14px</option>
                <option value={15}>15px</option>
                <option value={16}>16px</option>
              </select>
            </div>

            {/* Search Input triggering command palette */}
            <div style={{ position: 'relative' }}>
              <button 
                onClick={() => setIsCommandPaletteOpen(true)}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: '8px',
                  width: '200px',
                  padding: '6px 12px',
                  borderRadius: '8px',
                  border: '1px solid var(--border)',
                  background: 'rgba(255,255,255,0.03)',
                  color: 'var(--muted)',
                  fontSize: '0.78rem',
                  textAlign: 'left',
                  cursor: 'pointer'
                }}
              >
                <Search size={12} />
                <span>Search (Ctrl+K)...</span>
              </button>
            </div>

            {/* Organizations Switcher */}
            {organizations.length > 0 && (
              <select
                value={currentOrgId || ''}
                onChange={(e) => setCurrentOrgId(e.target.value)}
                style={{
                  padding: '6px 10px',
                  borderRadius: '8px',
                  border: '1px solid var(--border)',
                  backgroundColor: 'rgba(15,25,41,0.95)',
                  color: 'var(--text)',
                  fontSize: '0.8rem',
                  fontWeight: 600
                }}
              >
                {organizations.map((org) => (
                  <option key={org.id} value={org.id}>
                    {org.name}
                  </option>
                ))}
              </select>
            )}

            {/* Notification Badge Trigger */}
            <div style={{ position: 'relative' }} ref={notificationRef}>
              <button 
                onClick={() => setIsNotificationOpen(!isNotificationOpen)}
                style={{
                  padding: '6px 10px',
                  borderRadius: '8px',
                  border: '1px solid var(--border)',
                  background: isNotificationOpen ? 'var(--border-strong)' : 'rgba(255,255,255,0.03)',
                  color: 'var(--text)',
                  fontSize: '0.8rem',
                  fontWeight: 600,
                  display: 'flex',
                  alignItems: 'center',
                  gap: '6px',
                  cursor: 'pointer'
                }}
              >
                <Bell size={13} />
                <span style={{ color: 'var(--danger)', fontWeight: 'bold' }}>{notifications}</span>
              </button>

              {/* Notification Popover Panel */}
              {isNotificationOpen && (
                <div className="notification-panel" role="dialog" aria-label="System Notifications Inbox">
                  <div className="notification-panel__header">
                    <h4>Notifications Inbox</h4>
                    <button onClick={() => setIsNotificationOpen(false)}>Close</button>
                  </div>
                  <div className="notification-panel__list">
                    {systemNotifications.map(n => (
                      <div key={n.id} className={`notification-item notification-item--${n.type}`}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '4px' }}>
                          <span className={`notification-badge notification-badge--${n.type}`}>
                            {n.type.toUpperCase()}
                          </span>
                          <span className="notification-item__time">{n.time}</span>
                        </div>
                        <p className="notification-item__message">{n.message}</p>
                      </div>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* Server adapter badge */}
            <Badge tone="accent">Host: {server}</Badge>

            <Button variant="secondary" onClick={signOut}>
              Sign out
            </Button>
          </div>
        </header>

        {/* Content Node */}
        <main className="shell__content">{children}</main>
      </div>

      {/* Command Palette Modal */}
      {isCommandPaletteOpen && (
        <div className="command-palette-backdrop" role="dialog" aria-modal="true" aria-label="Universal Command Palette Launcher">
          <div className="command-palette" ref={paletteRef}>
            <div className="command-palette__search">
              <Search size={16} className="command-palette__search-icon" />
              <input
                ref={paletteInputRef}
                type="text"
                placeholder="Search commands or sitemap paths (e.g. Settings, Telemetry)..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <span className="command-palette__esc-hint">ESC to close</span>
            </div>

            <div className="command-palette__results">
              <div className="command-palette__section-title">Navigation Shortcuts</div>
              {filteredShortcuts.length ? (
                filteredShortcuts.map((item) => {
                  const Icon = getNavIcon(item.href)
                  return (
                    <button
                      key={item.href}
                      onClick={() => handleShortcutClick(item.href)}
                      className="command-palette__item"
                    >
                      <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                        <Icon size={14} style={{ color: 'var(--accent)' }} />
                        <div>
                          <div className="command-palette__item-label">{item.label}</div>
                          <div className="command-palette__item-sub">{item.subtitle}</div>
                        </div>
                      </div>
                      <span className="command-palette__item-group">{item.group}</span>
                    </button>
                  )
                })
              ) : (
                <div style={{ padding: '24px', textAlign: 'center', color: 'var(--muted)' }}>
                  No matching shortcuts found.
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Mobile Bottom Tab Navigation */}
      <nav className="mobile-bottom-nav" aria-label="Mobile Navigation Bar">
        <Link href="/dashboard" className={`mobile-bottom-nav__item${pathname === '/dashboard' ? ' is-active' : ''}`}>
          <LayoutDashboard size={18} />
          <span className="mobile-bottom-nav__label">Home</span>
        </Link>
        <Link href="/projects" className={`mobile-bottom-nav__item${pathname === '/projects' ? ' is-active' : ''}`}>
          <FolderKanban size={18} />
          <span className="mobile-bottom-nav__label">Projects</span>
        </Link>
        <Link href="/deploy" className={`mobile-bottom-nav__item${pathname === '/deploy' ? ' is-active' : ''}`}>
          <History size={18} />
          <span className="mobile-bottom-nav__label">Deploy</span>
        </Link>
        <Link href="/settings" className={`mobile-bottom-nav__item${pathname === '/settings' ? ' is-active' : ''}`}>
          <Settings size={18} />
          <span className="mobile-bottom-nav__label">Settings</span>
        </Link>
      </nav>
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
