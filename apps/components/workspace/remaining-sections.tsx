import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import type { WorkspaceDomainInfo, WorkspaceLogEntry, AIContextInfo, MCPInfo, WorkspaceDoctorData, WorkspaceSettingsData } from '@/lib/types'

// Database section
export function DatabaseSection({ data }: { data: { kind: string; files: string[]; migrations: string[]; environment: string[] } }) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Database</div>
          <h2 className="ws-section__title">Database</h2>
          <p className="ws-section__subtitle">{data.kind || 'No database detected'}</p>
        </div>
        {data.kind && <Badge tone="accent">{data.kind}</Badge>}
      </div>
      <div className="ws-grid-2">
        <Card>
          <div className="section-header-lite"><h3>Schema Files</h3></div>
          {data.files?.length ? (
            <div className="ws-detail-list">{data.files.map((f) => <div key={f} className="ws-detail-row"><code className="ws-mono">{f}</code></div>)}</div>
          ) : <EmptyState title="No schema files" description="No .sql files found." />}
        </Card>
        <Card>
          <div className="section-header-lite"><h3>Migrations</h3></div>
          {data.migrations?.length ? (
            <div className="ws-detail-list">{data.migrations.map((m) => <div key={m} className="ws-detail-row"><code className="ws-mono">{m}</code></div>)}</div>
          ) : <EmptyState title="No migrations" description="No migration files detected." />}
        </Card>
      </div>
    </div>
  )
}

// Domains section
export function DomainsSection({ data }: { data: WorkspaceDomainInfo[] }) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Domains</div>
          <h2 className="ws-section__title">Domains</h2>
          <p className="ws-section__subtitle">{data?.length || 0} domains configured</p>
        </div>
      </div>
      {data?.length ? (
        <Card>
          <div className="ws-detail-list">
            {data.map((d) => (
              <div key={d.host} className="ws-detail-row">
                <strong>{d.host}</strong>
                <Badge tone={d.ssl === 'Valid' ? 'success' : 'warning'}>{d.ssl}</Badge>
                <code className="ws-mono">{d.target}</code>
              </div>
            ))}
          </div>
        </Card>
      ) : <EmptyState title="No domains" description="Domain detection reads nginx configurations." />}
    </div>
  )
}

// Logs section
export function LogsSection({ data }: { data: WorkspaceLogEntry[] }) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Logs</div>
          <h2 className="ws-section__title">Logs</h2>
          <p className="ws-section__subtitle">{data?.length || 0} log entries</p>
        </div>
      </div>
      {data?.length ? (
        <Card>
          <div className="ws-log-list">
            {data.map((entry, i) => (
              <div key={i} className="ws-log-row">
                <span className={`ws-log-level ws-log-level--${entry.level}`}>{entry.level}</span>
                <span className="ws-log-source">{entry.source}</span>
                <span className="ws-log-msg">{entry.message}</span>
              </div>
            ))}
          </div>
        </Card>
      ) : <EmptyState title="No logs" description="Log files from .devserver/logs/ will appear here." />}
    </div>
  )
}

// AI Context Builder section
export function AISection({ data }: { data: AIContextInfo }) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">AI Context Builder</div>
          <h2 className="ws-section__title">AI</h2>
          <p className="ws-section__subtitle">Structured context for AI consumption</p>
        </div>
        <Badge tone="accent">Generated</Badge>
      </div>
      <div className="ws-grid-2">
        <Card>
          <div className="section-header-lite"><h3>Context Summary</h3></div>
          <div className="ws-detail-list">
            <div className="ws-detail-row"><span>Workspace</span><strong>{data.workspace?.name || '—'}</strong></div>
            <div className="ws-detail-row"><span>Framework</span><strong>{data.architecture?.framework || '—'}</strong></div>
            <div className="ws-detail-row"><span>Languages</span><span>{data.architecture?.languages?.join(', ') || '—'}</span></div>
            <div className="ws-detail-row"><span>Routes</span><strong>{data.routes?.api?.length || 0}</strong></div>
            <div className="ws-detail-row"><span>Database</span><strong>{data.database?.kind || 'None'}</strong></div>
            <div className="ws-detail-row"><span>Generated</span><span>{data.generatedAt ? new Date(data.generatedAt).toLocaleString() : '—'}</span></div>
          </div>
        </Card>
        <Card>
          <div className="section-header-lite"><h3>Dependencies</h3></div>
          <div className="ws-tag-list">
            {data.dependencies?.frontend?.map((d) => <Badge key={d} tone="accent">{d}</Badge>)}
            {data.dependencies?.backend?.map((d) => <Badge key={d} tone="info">{d}</Badge>)}
            {data.dependencies?.database?.map((d) => <Badge key={d} tone="warning">{d}</Badge>)}
          </div>
        </Card>
      </div>
      {data.tasks?.suggested?.length > 0 && (
        <Card>
          <div className="section-header-lite"><h3>Suggested Tasks</h3></div>
          <div className="ws-detail-list">
            {data.tasks.suggested.map((t) => <div key={t} className="ws-detail-row"><span>{t}</span></div>)}
          </div>
        </Card>
      )}
    </div>
  )
}

// Knowledge section
export function KnowledgeSection({ data }: { data: { files: string[]; topics: string[] } }) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Knowledge</div>
          <h2 className="ws-section__title">Knowledge</h2>
          <p className="ws-section__subtitle">{data.files?.length || 0} docs · {data.topics?.length || 0} topics</p>
        </div>
      </div>
      <div className="ws-grid-2">
        <Card>
          <div className="section-header-lite"><h3>Files</h3></div>
          {data.files?.length ? (
            <div className="ws-detail-list">{data.files.map((f) => <div key={f} className="ws-detail-row"><code className="ws-mono">{f}</code></div>)}</div>
          ) : <EmptyState title="No docs" description="Add markdown files or use .devserver/knowledge/." />}
        </Card>
        <Card>
          <div className="section-header-lite"><h3>Topics</h3></div>
          {data.topics?.length ? (
            <div className="ws-tag-list">{data.topics.map((t) => <Badge key={t} tone="info">{t}</Badge>)}</div>
          ) : <EmptyState title="No topics" description="Topics are derived from documentation files." />}
        </Card>
      </div>
    </div>
  )
}

// Doctor section
export function DoctorSection({ data }: { data: WorkspaceDoctorData }) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Project Doctor</div>
          <h2 className="ws-section__title">Doctor</h2>
          <p className="ws-section__subtitle">Health score: {data.health?.score || 0}</p>
        </div>
        <Badge tone={data.health?.score >= 60 ? 'success' : 'warning'}>Score {data.health?.score || 0}</Badge>
      </div>
      {data.checks?.length > 0 && (
        <Card>
          <div className="section-header-lite"><h3>Health Checks</h3></div>
          <div className="ws-detail-list">
            {data.checks.map((c) => (
              <div key={c.label} className="ws-detail-row">
                <span>{c.label}</span>
                <Badge tone={c.status === 'passed' ? 'success' : c.status === 'warning' ? 'warning' : 'danger'}>{c.status}</Badge>
                <span style={{ color: 'var(--muted)', fontSize: '0.88rem' }}>{c.detail}</span>
              </div>
            ))}
          </div>
        </Card>
      )}
      {data.recommendations?.length > 0 && (
        <Card>
          <div className="section-header-lite"><h3>Recommendations</h3></div>
          <ul className="bullet-list">{data.recommendations.map((r) => <li key={r}>{r}</li>)}</ul>
        </Card>
      )}
    </div>
  )
}

// Settings section
export function SettingsSection({ data }: { data: WorkspaceSettingsData }) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Workspace Settings</div>
          <h2 className="ws-section__title">Settings</h2>
          <p className="ws-section__subtitle">Configuration and index metadata</p>
        </div>
      </div>
      <div className="ws-grid-2">
        <Card>
          <div className="section-header-lite"><h3>Workspace</h3></div>
          <div className="ws-detail-list">
            <div className="ws-detail-row"><span>ID</span><code className="ws-mono">{data.workspace?.id}</code></div>
            <div className="ws-detail-row"><span>Name</span><strong>{data.workspace?.name}</strong></div>
            <div className="ws-detail-row"><span>Root</span><code className="ws-mono">{data.workspace?.root}</code></div>
            <div className="ws-detail-row"><span>Framework</span><strong>{data.workspace?.framework}</strong></div>
            <div className="ws-detail-row"><span>Runtime</span><strong>{data.workspace?.runtime}</strong></div>
          </div>
        </Card>
        <Card>
          <div className="section-header-lite"><h3>Index</h3></div>
          <div className="ws-detail-list">
            <div className="ws-detail-row"><span>Version</span><strong>{data.index?.version}</strong></div>
            <div className="ws-detail-row"><span>Generated</span><span>{data.index?.generatedAt ? new Date(data.index.generatedAt).toLocaleString() : '—'}</span></div>
            <div className="ws-detail-row"><span>Fingerprint</span><code className="ws-mono" style={{ fontSize: '0.8rem' }}>{data.cache?.fingerprint?.slice(0, 16) || '—'}</code></div>
            <div className="ws-detail-row"><span>Changed</span><strong>{data.cache?.changed?.length || 0}</strong></div>
          </div>
        </Card>
      </div>
    </div>
  )
}

// MCP Manager section
export function MCPSection({ data }: { data: MCPInfo }) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">MCP Manager</div>
          <h2 className="ws-section__title">MCP Providers</h2>
          <p className="ws-section__subtitle">{data.providers?.length || 0} providers configured</p>
        </div>
      </div>
      {data.providers?.length ? (
        <div className="ws-service-grid">
          {data.providers.map((p) => (
            <Card key={p.name}>
              <div className="ws-infra-card__top">
                <strong>{p.name}</strong>
                <Badge tone={p.enabled ? (p.healthy ? 'success' : 'warning') : 'neutral'}>
                  {p.enabled ? (p.healthy ? 'Healthy' : 'Unhealthy') : 'Disabled'}
                </Badge>
              </div>
              <div className="ws-detail-list" style={{ marginTop: 8 }}>
                <div className="ws-detail-row"><span>Protocol</span><Badge tone="info">{p.protocol}</Badge></div>
                <div className="ws-detail-row"><span>Endpoint</span><code className="ws-mono" style={{ fontSize: '0.82rem' }}>{p.endpoint}</code></div>
              </div>
            </Card>
          ))}
        </div>
      ) : <EmptyState title="No providers" description="Configure MCP providers in .devserver/mcp.json" />}
    </div>
  )
}
