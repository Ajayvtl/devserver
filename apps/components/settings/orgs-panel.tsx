'use client'

import { useEffect, useState } from 'react'
import { Card } from '../ui/card'
import { Button } from '../ui/button'
import { DataTable } from '../ui/table'
import { Input } from '../ui/input'
import { useAuth } from '../auth/auth-context'
import { requestOrFallback, request } from '@/lib/api/client'
import { useToast } from '../toast'

export function OrgsPanel() {
	const { currentOrgId, setCurrentOrgId } = useAuth()
	const { push } = useToast()

  const [orgs, setOrgs] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [newOrgName, setNewOrgName] = useState('')
  const [creating, setCreating] = useState(false)

  const loadOrgs = async () => {
    setLoading(true)
    try {
      const data = await requestOrFallback('/api/v1/organizations', [])
      setOrgs(data)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadOrgs()
  }, [])

  const handleCreate = async () => {
    if (!newOrgName.trim()) return
    setCreating(true)
    try {
      await request('/api/v1/organizations', {
        method: 'POST',
        body: { name: newOrgName }
      })
      push({ title: 'Success', message: 'Organization created.', tone: 'success' })
      setNewOrgName('')
      await loadOrgs()
      window.location.reload()
    } catch (err: any) {
      push({ title: 'Error', message: err.message, tone: 'danger' })
    } finally {
      setCreating(false)
    }
  }

  if (loading) return <div>Loading organizations...</div>

  return (
    <div className="stack" style={{ gap: '2rem' }}>
      <Card>
        <div className="card__eyebrow">Your Organizations</div>
        <div style={{ margin: '1rem 0' }}>
          <DataTable 
            columns={[{ header: 'Name' }, { header: 'Role' }, { header: 'Action' }]}
            rows={orgs.map(o => [
              o.name,
              o.id === currentOrgId ? 'Active' : 'Member',
              <Button 
                variant={o.id === currentOrgId ? 'secondary' : 'primary'} 
                disabled={o.id === currentOrgId}
                onClick={() => setCurrentOrgId(o.id)}
              >
                {o.id === currentOrgId ? 'Current' : 'Switch'}
              </Button>
            ])}
          />
        </div>
      </Card>
      
      <Card>
        <div className="card__eyebrow">Create Organization</div>
        <div className="flex" style={{ gap: '1rem', marginTop: '1rem', alignItems: 'flex-end' }}>
          <div style={{ flex: 1 }}>
            <Input 
              label="Organization Name"
              value={newOrgName}
              onChange={e => setNewOrgName(e.target.value)}
              placeholder="e.g. Acme Corp"
            />
          </div>
          <Button variant="primary" onClick={handleCreate} disabled={creating || !newOrgName.trim()}>
            {creating ? 'Creating...' : 'Create'}
          </Button>
        </div>
      </Card>
    </div>
  )
}
