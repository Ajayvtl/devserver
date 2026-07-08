'use client'

import { useEffect, useState, useCallback } from 'react'
import { Badge } from '../ui/badge'
import { Button } from '../ui/button'
import { Card } from '../ui/card'
import { DataTable } from '../ui/table'
import { SectionHeader } from '../ui/section-header'
import { Input } from '../ui/input'
import { useAuth } from '@/components/auth/auth-context'
import { PermissionSettingsWrite } from '@/lib/rbac/permissions'
import { request, requestOrFallback } from '@/lib/api/client'
import { useToast } from '../toast'
import { OrgsPanel } from './orgs-panel'
import { MembersPanel } from './members-panel'

export function SettingsPage() {
  const { currentOrgId, can } = useAuth()
  const { push } = useToast()

  const [settings, setSettings] = useState<{ key: string; value: string }[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  // Local state for edits
  const [theme, setTheme] = useState('dark')
  const [language, setLanguage] = useState('en-US')
  const [mfaEnabled, setMfaEnabled] = useState('false')
  const [activeTab, setActiveTab] = useState('general')

  const loadSettings = useCallback(async () => {
    if (!currentOrgId) return
    setLoading(true)
    try {
      const data = await requestOrFallback<{ key: string; value: string }[]>(`/api/v1/settings?scope=organization&ownerId=${currentOrgId}`, [], {
        headers: { 'X-Org-ID': currentOrgId }
      })
      setSettings(data)

      const t = data.find(s => s.key === 'theme')?.value
      if (t) setTheme(t)

      const l = data.find(s => s.key === 'language')?.value
      if (l) setLanguage(l)

      const m = data.find(s => s.key === 'mfa')?.value
      if (m) setMfaEnabled(m)

    } finally {
      setLoading(false)
    }
  }, [currentOrgId])

  useEffect(() => {
    loadSettings()
  }, [loadSettings])

  const handleSave = async () => {
    if (!currentOrgId) return
    setSaving(true)
    try {
      await request('/api/v1/settings', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: { scope: 'organization', ownerId: currentOrgId, key: 'theme', value: theme }
      })
      await request('/api/v1/settings', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: { scope: 'organization', ownerId: currentOrgId, key: 'language', value: language }
      })
      await request('/api/v1/settings', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: { scope: 'organization', ownerId: currentOrgId, key: 'mfa', value: mfaEnabled }
      })
      push({ title: 'Settings Saved', message: 'Organization preferences updated successfully.', tone: 'success' })
      loadSettings()
    } catch (err: any) {
      push({ title: 'Save Failed', message: err.message, tone: 'danger' })
    } finally {
      setSaving(false)
    }
  }

  const hasWriteAccess = can(PermissionSettingsWrite)

  if (loading) {
    return (
      <div className="page">
        <div className="skeleton-group">
          <div className="skeleton-line"></div>
          <div className="skeleton-line"></div>
          <div className="skeleton-line"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="page">
      <SectionHeader
        eyebrow="Settings"
        title="Organization Settings"
        description="Control platform preferences, security posture, and integrations."
        actions={
          <Button variant="primary" onClick={handleSave} disabled={saving || !hasWriteAccess || activeTab !== 'general'}>
            {saving ? 'Saving...' : 'Save changes'}
          </Button>
        }
      />

      <div className="tabs">
        <div className="tabs__list">
          <button
            type="button"
            onClick={() => setActiveTab('general')}
            className={`tabs__tab ${activeTab === 'general' ? 'is-active' : ''}`}
          >
            General
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('orgs')}
            className={`tabs__tab ${activeTab === 'orgs' ? 'is-active' : ''}`}
          >
            Organizations
          </button>
          <button
            type="button"
            onClick={() => setActiveTab('members')}
            className={`tabs__tab ${activeTab === 'members' ? 'is-active' : ''}`}
          >
            Members
          </button>
        </div>
      </div>

      {activeTab === 'general' && (
        <>
          <div className="page-grid--two">
            <Card>
              <div className="card__eyebrow">Preferences</div>
              <div className="stack">
                <Input
                  label="Theme"
                  value={theme}
                  onChange={e => setTheme(e.target.value)}
                  disabled={!hasWriteAccess}
                />
                <Input
                  label="Language"
                  value={language}
                  onChange={e => setLanguage(e.target.value)}
                  disabled={!hasWriteAccess}
                />
              </div>
            </Card>

            <Card>
              <div className="card__eyebrow">Security Posture</div>
              <div className="stack">
                <label className={`checkbox-card ${!hasWriteAccess ? 'is-disabled' : ''}`}>
                  <input
                    type="checkbox"
                    checked={mfaEnabled === 'true'}
                    onChange={e => setMfaEnabled(e.target.checked ? 'true' : 'false')}
                    disabled={!hasWriteAccess}
                  />
                  <span>Require Multi-Factor Authentication (MFA) for all members</span>
                </label>
              </div>
            </Card>
          </div>

          <div className="page-grid--two">
            <Card>
              <div className="card__eyebrow">Active Integrations</div>
              <DataTable
                columns={[
                  { header: 'Integration' },
                  { header: 'Status' }
                ]}
                rows={[
                  ['GitHub', <Badge key="github" tone="success">Connected</Badge>],
                  ['Slack', <Badge key="slack" tone="info">Pending</Badge>]
                ]}
              />
            </Card>
          </div>
        </>
      )}

      {activeTab === 'orgs' && (
        <OrgsPanel />
      )}

      {activeTab === 'members' && (
        <MembersPanel />
      )}
    </div>
  )
}
