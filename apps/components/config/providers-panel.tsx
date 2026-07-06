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
  metadata: Record<string, string>
}

export function ProvidersPanel() {
  const { currentOrgId, can } = useAuth()
  const { push } = useToast()
  const [providers, setProviders] = useState<ProviderConfig[]>([])
  const [loading, setLoading] = useState(true)
  const [showAdd, setShowAdd] = useState(false)
  const [editingProvider, setEditingProvider] = useState<ProviderConfig | null>(null)
  
  const [formName, setFormName] = useState('openai')
  const [formSecretRef, setFormSecretRef] = useState('')
  const [formEnabled, setFormEnabled] = useState(true)
  const [formModel, setFormModel] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [testingConnection, setTestingConnection] = useState(false)

  const loadProviders = async () => {
    if (!currentOrgId) return
    setLoading(true)
    try {
      const data = await requestOrFallback<ProviderConfig[]>(`/api/v1/providers?scope=organization&ownerId=${currentOrgId}`, [], {
        headers: { 'X-Org-ID': currentOrgId }
      })
      setProviders(data || [])
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadProviders()
  }, [currentOrgId])

  const openForm = (provider?: ProviderConfig) => {
    if (provider) {
      setEditingProvider(provider)
      setFormName(provider.name)
      setFormSecretRef(provider.secretRef)
      setFormEnabled(provider.enabled)
      setFormModel(provider.metadata?.defaultModel || '')
    } else {
      setEditingProvider(null)
      setFormName('openai')
      setFormSecretRef('')
      setFormEnabled(true)
      setFormModel('')
    }
    setShowAdd(true)
  }

  const closeForm = () => {
    setShowAdd(false)
    setEditingProvider(null)
  }

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!currentOrgId) return
    setIsSubmitting(true)
    try {
      await request('/api/v1/providers', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: {
          id: editingProvider?.id,
          scope: 'organization',
          ownerId: currentOrgId,
          name: formName,
          enabled: formEnabled,
          secretRef: formSecretRef,
          metadata: { defaultModel: formModel }
        }
      })
      push({ title: 'Provider Saved', message: `Successfully configured ${formName}.`, tone: 'success' })
      closeForm()
      loadProviders()
    } catch (err: any) {
      push({ title: 'Error', message: err.message, tone: 'danger' })
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleTestConnection = async () => {
    if (!currentOrgId) return
    setTestingConnection(true)
    try {
      if (!formSecretRef) throw new Error('Secret Reference is required for connection testing.')
      
      await request('/api/v1/providers/test', {
        method: 'POST',
        headers: { 'X-Org-ID': currentOrgId },
        body: {
          ownerId: currentOrgId,
          name: formName,
          secretRef: formSecretRef
        }
      })
      
      push({ title: 'Connection Successful', message: `Successfully authenticated with ${formName}.`, tone: 'success' })
    } catch (err: any) {
      push({ title: 'Connection Failed', message: err.message, tone: 'danger' })
    } finally {
      setTestingConnection(false)
    }
  }

  const hasWriteAccess = can('providers.write')

  return (
    <div className="stack" style={{ gap: '2rem' }}>
      <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
        <Button 
          variant="primary" 
          onClick={() => showAdd ? closeForm() : openForm()}
          disabled={!hasWriteAccess}
        >
          {showAdd ? 'Cancel' : 'Add AI Provider'}
        </Button>
      </div>

      {showAdd && (
        <Card className="glass-panel" style={{ animation: 'slideIn 0.3s ease-out' }}>
          <div className="card__eyebrow">{editingProvider ? 'Edit Provider' : 'Configure New Provider'}</div>
          <form onSubmit={handleSave} className="stack" style={{ marginTop: '1rem' }}>
            <div className="page-grid--two">
              <div className="input-group">
                <label className="input-label">Provider Type</label>
                <select 
                  value={formName} 
                  onChange={e => setFormName(e.target.value)} 
                  disabled={!!editingProvider}
                  style={{ padding: '0.625rem', borderRadius: '6px', border: '1px solid var(--border)', background: 'var(--bg-elevated)', color: 'var(--text-primary)', width: '100%' }}
                >
                  <option value="openai">OpenAI</option>
                  <option value="gemini">Google Gemini</option>
                  <option value="groq">Groq</option>
                  <option value="ollama">Ollama (Local)</option>
                </select>
              </div>
              <Input 
                label="Environment Secret Reference ID" 
                value={formSecretRef} 
                onChange={e => setFormSecretRef(e.target.value)}
                placeholder="e.g. SEC_OPENAI_API_KEY"
                required
              />
              <Input 
                label="Default Model (Optional)" 
                value={formModel} 
                onChange={e => setFormModel(e.target.value)}
                placeholder="e.g. gpt-4o, gemini-1.5-pro"
              />
              <label className="checkbox-card" style={{ marginTop: '1.5rem' }}>
                <input 
                  type="checkbox" 
                  checked={formEnabled} 
                  onChange={e => setFormEnabled(e.target.checked)}
                />
                <span>Enable this provider for AI requests</span>
              </label>
            </div>
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '1rem', marginTop: '1rem' }}>
              <Button type="button" variant="secondary" onClick={handleTestConnection} disabled={testingConnection || !formSecretRef}>
                {testingConnection ? 'Testing...' : 'Test Connection'}
              </Button>
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
              { header: 'Default Model' },
              { header: 'Actions' }
            ]}
            rows={providers.map(p => [
              <strong key="name" style={{ textTransform: 'capitalize' }}>{p.name}</strong>,
              <Badge key="status" tone={p.enabled ? 'success' : 'neutral'}>{p.enabled ? 'Active' : 'Disabled'}</Badge>,
              <span key="ref" style={{ fontFamily: 'monospace', fontSize: '0.875rem' }}>{p.secretRef}</span>,
              <span key="model">{p.metadata?.defaultModel || '-'}</span>,
              <Button key="action" variant="ghost" disabled={!hasWriteAccess} onClick={() => openForm(p)}>Edit</Button>
            ])}
          />
        )}
      </Card>
    </div>
  )
}
