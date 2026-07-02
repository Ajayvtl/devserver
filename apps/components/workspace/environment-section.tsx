'use client'

import { useState } from 'react'
import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import type { WorkspaceEnvironmentInfo } from '@/lib/types'

interface Props { data: WorkspaceEnvironmentInfo }

export function EnvironmentSection({ data }: Props) {
  const [showSecrets, setShowSecrets] = useState(false)

  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Environment Manager</div>
          <h2 className="ws-section__title">Environment</h2>
          <p className="ws-section__subtitle">{data.files?.length || 0} env files · {data.keys?.length || 0} variables</p>
        </div>
        <button className="button button--secondary" onClick={() => setShowSecrets(!showSecrets)}>
          {showSecrets ? 'Hide secrets' : 'Show secrets'}
        </button>
      </div>

      {data.files?.length > 0 && (
        <Card>
          <div className="section-header-lite"><h3>Environment Files</h3></div>
          <div className="ws-tag-list">
            {data.files.map((f) => <Badge key={f} tone="info">{f}</Badge>)}
          </div>
        </Card>
      )}

      {data.secrets?.length > 0 && (
        <Card>
          <div className="section-header-lite"><h3>Secrets ({data.secrets.length})</h3></div>
          <div className="ws-tag-list">
            {data.secrets.map((s) => <Badge key={s} tone="danger">{s}</Badge>)}
          </div>
        </Card>
      )}

      <Card>
        <div className="section-header-lite"><h3>Variables</h3></div>
        {data.keys?.length ? (
          <div className="ws-env-table">
            {data.keys.map((key) => {
              const isSecret = data.secrets?.includes(key)
              const val = data.preview?.[key] || ''
              return (
                <div key={key} className="ws-env-row">
                  <code className="ws-mono ws-env-key">{key}</code>
                  <span className="ws-env-val">
                    {isSecret && !showSecrets ? '••••••••' : val || '(empty)'}
                  </span>
                  {isSecret && <Badge tone="danger">Secret</Badge>}
                </div>
              )
            })}
          </div>
        ) : (
          <EmptyState title="No variables" description="No environment variables detected." />
        )}
      </Card>
    </div>
  )
}
