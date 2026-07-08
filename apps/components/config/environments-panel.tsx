import { useEffect, useState, useCallback } from 'react'
import { request, requestOrFallback } from '@/lib/api/client'
import { useAuth } from '@/components/auth/auth-context'
import { PermissionEnvironmentsWrite } from '@/lib/rbac/permissions'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { DataTable } from '@/components/ui/table'
import { useToast } from '@/components/toast'

interface EnvData {
  id: string
  name: string
  type: string
}

interface Variable {
  key: string
  value: string
}

interface SecretRef {
  key: string
}

export function EnvironmentsPanel() {
  const { currentOrgId, can } = useAuth()
  const { push } = useToast()

  const [environments, setEnvironments] = useState<EnvData[]>([])
  const [loading, setLoading] = useState(true)
  const [showAdd, setShowAdd] = useState(false)
  const [newName, setNewName] = useState('')
  const [newType, setNewType] = useState('development')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const [activeEnv, setActiveEnv] = useState<EnvData | null>(null)
  const [variables, setVariables] = useState<Variable[]>([])
  const [secrets, setSecrets] = useState<SecretRef[]>([])

  const [varKey, setVarKey] = useState('')
  const [varValue, setVarValue] = useState('')
  const [secKey, setSecKey] = useState('')
  const [secValue, setSecValue] = useState('')

  const loadEnvironments = useCallback(async () => {
    if (!currentOrgId) return
    setLoading(true)
    try {
      const data = await requestOrFallback<EnvData[]>(`/api/v1/environments?ownerId=${currentOrgId}`, [], {
        headers: { 'X-Org-ID': currentOrgId }
      })
      setEnvironments(data || [])
    } finally {
      setLoading(false)
    }
  }, [currentOrgId])

  const loadDetails = useCallback(async (env: EnvData) => {
    setActiveEnv(env)
    const [vars, secs] = await Promise.all([
      requestOrFallback<Variable[]>(`/api/v1/variables?envId=${env.id}`, [], { headers: { 'X-Org-ID': currentOrgId || '' } }),
      requestOrFallback<SecretRef[]>(`/api/v1/secrets?envId=${env.id}`, [], { headers: { 'X-Org-ID': currentOrgId || '' } })
    ])
    setVariables(vars || [])
    setSecrets(secs || [])
  }, [currentOrgId])

  useEffect(() => {
    loadEnvironments()
  }, [loadEnvironments])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!currentOrgId) return
    setIsSubmitting(true)
    try {
      await request('/api/v1/environments', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: { ownerId: currentOrgId, name: newName, type: newType }
      })
      push({ title: 'Environment Created', message: `Created ${newName}.`, tone: 'success' })
      setShowAdd(false)
      loadEnvironments()
    } catch (err: any) {
      push({ title: 'Error', message: err.message, tone: 'danger' })
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleAddVar = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!activeEnv) return
    try {
      await request('/api/v1/variables', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId || '' },
        body: { envId: activeEnv.id, key: varKey, value: varValue }
      })
      setVarKey('')
      setVarValue('')
      loadDetails(activeEnv)
    } catch (err: any) {
      push({ title: 'Error', message: err.message, tone: 'danger' })
    }
  }

  const handleAddSec = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!activeEnv) return
    try {
      await request('/api/v1/secrets', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId || '' },
        body: { envId: activeEnv.id, key: secKey, plaintext: secValue }
      })
      setSecKey('')
      setSecValue('')
      loadDetails(activeEnv)
      push({ title: 'Secret Saved', message: `Secret ${secKey} encrypted and stored safely.`, tone: 'success' })
    } catch (err: any) {
      push({ title: 'Error', message: err.message, tone: 'danger' })
    }
  }

  const hasWriteAccess = can(PermissionEnvironmentsWrite)

  return (
    <div className="stack" style={{ gap: '2rem' }}>
      {!activeEnv && (
        <>
          <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
            <Button variant="primary" onClick={() => setShowAdd(!showAdd)} disabled={!hasWriteAccess}>
              {showAdd ? 'Cancel' : 'Create Environment'}
            </Button>
          </div>

          {showAdd && (
            <Card className="glass-panel" style={{ animation: 'slideIn 0.3s ease-out' }}>
              <div className="card__eyebrow">New Environment</div>
              <form onSubmit={handleCreate} className="stack" style={{ marginTop: '1rem' }}>
                <div className="page-grid--two">
                  <Input
                    label="Environment Name"
                    value={newName}
                    onChange={e => setNewName(e.target.value)}
                    required
                  />
                  <div className="input-group">
                    <label className="input-label">Environment Type</label>
                    <select
                      value={newType}
                      onChange={e => setNewType(e.target.value)}
                      style={{ padding: '0.625rem', borderRadius: '6px', border: '1px solid var(--border)', background: 'var(--bg-elevated)', color: 'var(--text-primary)', width: '100%' }}
                    >
                      <option value="development">Development</option>
                      <option value="staging">Staging</option>
                      <option value="production">Production</option>
                      <option value="custom">Custom</option>
                    </select>
                  </div>
                </div>
                <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '1rem' }}>
                  <Button type="submit" variant="primary" disabled={isSubmitting}>
                    {isSubmitting ? 'Creating...' : 'Create Environment'}
                  </Button>
                </div>
              </form>
            </Card>
          )}

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
                  <Button key="action" variant="ghost" onClick={() => loadDetails(e)}>Manage</Button>
                ])}
              />
            )}
          </Card>
        </>
      )}

      {activeEnv && (
        <div className="stack" style={{ gap: '2rem', animation: 'slideIn 0.2s ease-out' }}>
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <h2 style={{ fontSize: '1.25rem', fontWeight: 600 }}>{activeEnv.name} <Badge tone="info">{activeEnv.type}</Badge></h2>
            <Button variant="secondary" onClick={() => setActiveEnv(null)}>Back to Environments</Button>
          </div>

          <div className="page-grid--two">
            <Card>
              <div className="card__eyebrow">Variables</div>
              <form onSubmit={handleAddVar} style={{ display: 'flex', gap: '0.5rem', margin: '1rem 0' }}>
                <Input placeholder="Key" value={varKey} onChange={e => setVarKey(e.target.value)} required />
                <Input placeholder="Value" value={varValue} onChange={e => setVarValue(e.target.value)} required />
                <Button type="submit" variant="secondary" disabled={!hasWriteAccess}>Add</Button>
              </form>
              <DataTable
                columns={[{ header: 'Key' }, { header: 'Value' }]}
                rows={variables.map(v => [
                  <span key="k" style={{ fontFamily: 'monospace' }}>{v.key}</span>,
                  <span key="v">{v.value}</span>
                ])}
              />
            </Card>

            <Card>
              <div className="card__eyebrow">Secrets (Zero-Trust)</div>
              <form onSubmit={handleAddSec} style={{ display: 'flex', gap: '0.5rem', margin: '1rem 0' }}>
                <Input placeholder="Key" value={secKey} onChange={e => setSecKey(e.target.value)} required />
                <Input placeholder="Secret Value" type="password" value={secValue} onChange={e => setSecValue(e.target.value)} required />
                <Button type="submit" variant="secondary" disabled={!hasWriteAccess}>Store</Button>
              </form>
              <DataTable
                columns={[{ header: 'Key' }, { header: 'Status' }]}
                rows={secrets.map(s => [
                  <span key="k" style={{ fontFamily: 'monospace' }}>{s.key}</span>,
                  <Badge key="s" tone="success">Encrypted</Badge>
                ])}
              />
            </Card>
          </div>
        </div>
      )}
    </div>
  )
}
