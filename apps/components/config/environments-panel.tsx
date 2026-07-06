'use client'

import { useEffect, useState } from 'react'
import { request, requestOrFallback } from '@/lib/api/client'
import { useAuth } from '@/components/auth/auth-context'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { DataTable } from '@/components/ui/table'
import { useToast } from '@/components/toast'

export function EnvironmentsPanel() {
  const { currentOrgId, can } = useAuth()
  const { push } = useToast()
  const [environments, setEnvironments] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  const loadEnvironments = async () => {
    if (!currentOrgId) return
    setLoading(true)
    try {
      const data = await requestOrFallback<any[]>(`/api/v1/environments?ownerId=${currentOrgId}`, [], {
        headers: { 'X-Org-ID': currentOrgId }
      })
      setEnvironments(data)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadEnvironments()
  }, [currentOrgId])

  const hasWriteAccess = can('environments.write')

  return (
    <div className="stack" style={{ gap: '2rem' }}>
      <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Button variant="primary" disabled={!hasWriteAccess}>
          Create Environment
        </Button>
      </div>

      <Card>
        <div className="card__eyebrow">Environments</div>
        {loading ? (
          <div style={{ padding: '2rem', textAlign: 'center', opacity: 0.7 }}>Loading...</div>
        ) : environments.length === 0 ? (
          <div style={{ padding: '3rem', textAlign: 'center', opacity: 0.5 }}>
            No environments found. Start by creating a Development environment.
          </div>
        ) : (
          <DataTable
            columns={[
              { header: 'Name' },
              { header: 'Type' },
              { header: 'Actions' }
            ]}
            rows={environments.map(e => [
              <strong key="name">{e.name}</strong>,
              <Badge key="type" tone="info">{e.type}</Badge>,
              <Button key="action" variant="ghost" disabled={!hasWriteAccess}>Manage</Button>
            ])}
          />
        )}
      </Card>
    </div>
  )
}
