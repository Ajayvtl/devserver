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

interface ProviderConfig {
  id: string
  name: string
  enabled: boolean
  secretRef: string
  metadata: any
}

export function ProvidersPanel() {
  const { currentOrgId, can } = useAuth()
  const { push } = useToast()
  const [providers, setProviders] = useState<ProviderConfig[]>([])
  const [loading, setLoading] = useState(true)
  const [showAdd, setShowAdd] = useState(false)
  
  const [newName, setNewName] = useState('')
  const [newSecretRef, setNewSecretRef] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const loadProviders = async () => {
    if (!currentOrgId) return
    setLoading(true)
    try {
      const data = await requestOrFallback<ProviderConfig[]>(`/api/v1/providers?scope=organization&ownerId=${currentOrgId}`, [], {
        headers: { 'X-Org-ID': currentOrgId }
      })
      setProviders(data)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadProviders()
  }, [currentOrgId])

  const handleAdd = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!currentOrgId) return
    setIsSubmitting(true)
    try {
      await request('/api/v1/providers', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: {
          scope: 'organization',
          ownerId: currentOrgId,
          name: newName,
          enabled: true,
          secretRef: newSecretRef,
          metadata: {}
        }
      })
      push({ title: 'Provider Added', message: 'Successfully configured new AI provider.', tone: 'success' })
      setShowAdd(false)
      loadProviders()
    } catch (err: any) {
      push({ title: 'Error', message: err.message, tone: 'danger' })
    } finally {
      setIsSubmitting(false)
    }
  }

  const hasWriteAccess = can('providers.write')

  return (
    <div className="stack" style={{ gap: '2rem' }}>
      <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Button 
          variant="primary" 
          onClick={() => setShowAdd(!showAdd)}
          disabled={!hasWriteAccess}
        >
          {showAdd ? 'Cancel' : 'Add Provider'}
        </Button>
      </div>

      {showAdd && (
        <Card className="glass-panel" style={{ animation: 'slideIn 0.3s ease-out' }}>
          <div className="card__eyebrow">Configure New Provider</div>
          <form onSubmit={handleAdd} className="stack" style={{ marginTop: '1rem' }}>
            <div className="page-grid--two">
              <Input 
                label="Provider Name (e.g. openai, gemini, groq)" 
                value={newName} 
                onChange={e => setNewName(e.target.value)}
                required
              />
              <Input 
                label="Secret Reference ID" 
                value={newSecretRef} 
                onChange={e => setNewSecretRef(e.target.value)}
                required
              />
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '1rem' }}>
              <Button type="submit" variant="primary" disabled={isSubmitting}>
                {isSubmitting ? 'Saving...' : 'Save Configuration'}
              </Button>
            </div>
          </form>
        </Card>
      )}

      <Card>
        <div className="card__eyebrow">Active Integrations</div>
        {loading ? (
          <div style={{ padding: '2rem', textAlign: 'center', opacity: 0.7 }}>Loading providers...</div>
        ) : providers.length === 0 ? (
          <div style={{ padding: '3rem', textAlign: 'center', opacity: 0.5 }}>
            No providers configured for this organization yet.
          </div>
        ) : (
          <DataTable
            columns={[
              { header: 'Provider' },
              { header: 'Status' },
              { header: 'Secret Ref' },
              { header: 'Actions' }
            ]}
            rows={providers.map(p => [
              <strong key="name">{p.name}</strong>,
              <Badge key="status" tone={p.enabled ? 'success' : 'neutral'}>{p.enabled ? 'Active' : 'Disabled'}</Badge>,
              <span key="ref" style={{ fontFamily: 'monospace', fontSize: '0.875rem' }}>{p.secretRef}</span>,
              <Button key="action" variant="ghost" disabled={!hasWriteAccess}>Edit</Button>
            ])}
          />
        )}
      </Card>
    </div>
  )
}
