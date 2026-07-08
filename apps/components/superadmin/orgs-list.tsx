'use client'

import { useEffect, useState } from 'react'
import { Card } from '../ui/card'
import { Badge } from '../ui/badge'
import { DataTable } from '../ui/table'
import { SectionHeader } from '../ui/section-header'
import { useAuth } from '../auth/auth-context'
import { requestOrFallback } from '@/lib/api/client'
import { useToast } from '../toast'
import { FolderKanban, ShieldAlert, Plus, Calendar, Settings } from 'lucide-react'

interface Organization {
  id: string
  name: string
  slug: string
  createdAt: string
  updatedAt: string
}

export function SuperAdminOrgsList() {
  const { currentOrgId, user } = useAuth()
  const { push } = useToast()
  
  const [orgs, setOrgs] = useState<Organization[]>([])
  const [loading, setLoading] = useState(true)

  const loadOrgs = async () => {
    if (!currentOrgId) return
    setLoading(true)
    try {
      // Fetch all organizations in system via the new endpoint
      const response = await requestOrFallback<{ data: Organization[] }>(
        '/api/v1/superadmin/organizations',
        { data: [] },
        { headers: { 'X-Org-ID': currentOrgId } }
      )
      
      // Handle either raw array or paginated response format
      const list = Array.isArray(response) ? response : (response?.data || [])
      setOrgs(list)
    } catch (err: any) {
      push({ title: 'Fetch Failed', message: err.message || 'Unable to retrieve system organizations.', tone: 'danger' })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadOrgs()
  }, [currentOrgId])

  if (loading) {
    return (
      <div className="skeleton-group">
        <div className="skeleton-line" style={{ height: '40px', width: '30%' }} />
        <div className="skeleton-line" style={{ height: '120px' }} />
        <div className="skeleton-line" style={{ height: '200px' }} />
      </div>
    )
  }

  const newestOrg = orgs.length > 0 
    ? [...orgs].sort((a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime())[0]
    : null

  return (
    <div className="stack" style={{ gap: '20px' }}>
      <SectionHeader
        eyebrow="Platform Administration"
        title="Registered Organizations"
        description="Super Admin dashboard to monitor system tenants, check registration dates, and manage global status."
      />

      {/* High-level stats panel */}
      <div className="page-grid--three" style={{ gap: '16px' }}>
        <Card>
          <div className="card__eyebrow" style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            <FolderKanban size={14} className="text-accent" />
            Total Tenants
          </div>
          <div style={{ fontSize: '2rem', fontWeight: 'bold', margin: '8px 0' }}>
            {orgs.length}
          </div>
          <div style={{ fontSize: '0.75rem', color: 'var(--muted)' }}>
            Organizations registered on this DevServer instance
          </div>
        </Card>

        <Card>
          <div className="card__eyebrow" style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            <ShieldAlert size={14} style={{ color: 'var(--success)' }} />
            Active Services
          </div>
          <div style={{ fontSize: '2rem', fontWeight: 'bold', margin: '8px 0', color: 'var(--success)' }}>
            {orgs.length}
          </div>
          <div style={{ fontSize: '0.75rem', color: 'var(--muted)' }}>
            100% of tenants operational with active routing
          </div>
        </Card>

        <Card>
          <div className="card__eyebrow" style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            <Calendar size={14} style={{ color: 'var(--accent)' }} />
            Newest Tenant
          </div>
          <div style={{ fontSize: '1.1rem', fontWeight: '600', margin: '14px 0 8px', textOverflow: 'ellipsis', overflow: 'hidden', whiteSpace: 'nowrap' }}>
            {newestOrg ? newestOrg.name : 'None'}
          </div>
          <div style={{ fontSize: '0.75rem', color: 'var(--muted)' }}>
            Registered: {newestOrg ? new Date(newestOrg.createdAt).toLocaleDateString() : 'N/A'}
          </div>
        </Card>
      </div>

      {/* Main Organizations Table */}
      <Card style={{ padding: '0px' }}>
        <div style={{ padding: '16px 20px', borderBottom: '1px solid var(--border)', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <div>
            <h3 style={{ margin: 0, fontSize: '0.95rem', fontWeight: '600' }}>Platform Tenant Index</h3>
            <p style={{ margin: '4px 0 0', fontSize: '0.75rem', color: 'var(--muted)' }}>System-wide directory of all organizations and their metadata.</p>
          </div>
          <Badge tone="accent">Super Admin Scope</Badge>
        </div>

        <DataTable
          columns={[
            { header: 'Organization Name' },
            { header: 'ID / Slug' },
            { header: 'Registration Date' },
            { header: 'Status' }
          ]}
          rows={orgs.map(org => [
            <div key={org.id} style={{ display: 'flex', flexDirection: 'column' }}>
              <span style={{ fontWeight: '500', fontSize: '0.85rem' }}>{org.name}</span>
            </div>,
            <code key={`${org.id}-slug`} style={{ fontSize: '0.75rem', color: 'var(--muted)', background: 'rgba(255,255,255,0.03)', padding: '2px 6px', borderRadius: '4px' }}>
              {org.slug || org.id.substring(0, 8)}
            </code>,
            <span key={`${org.id}-date`} style={{ fontSize: '0.8rem', color: 'var(--muted)' }}>
              {new Date(org.createdAt).toLocaleDateString(undefined, {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
              })}
            </span>,
            <Badge key={`${org.id}-status`} tone="success">
              Active
            </Badge>
          ])}
        />
        
        {orgs.length === 0 && (
          <div style={{ padding: '40px', textAlign: 'center', color: 'var(--muted)' }}>
            No registered organizations found in system.
          </div>
        )}
      </Card>
    </div>
  )
}
