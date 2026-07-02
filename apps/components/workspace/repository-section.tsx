import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { EmptyState } from '../ui/empty-state'
import type { WorkspaceGitInfo } from '@/lib/types'

interface Props { data: WorkspaceGitInfo }

export function RepositorySection({ data }: Props) {
  return (
    <div className="ws-section">
      <div className="ws-section__header">
        <div>
          <div className="card__eyebrow">Git Manager</div>
          <h2 className="ws-section__title">Repository</h2>
          <p className="ws-section__subtitle">Real git integration — no heuristics</p>
        </div>
        <Badge tone={data.clean ? 'success' : 'warning'}>{data.clean ? 'Clean' : 'Dirty'}</Badge>
      </div>

      <div className="ws-grid-2">
        <Card>
          <div className="section-header-lite"><h3>Branch & Commit</h3></div>
          <div className="ws-detail-list">
            <div className="ws-detail-row"><span>Branch</span><Badge tone="accent">{data.branch || '—'}</Badge></div>
            <div className="ws-detail-row"><span>Commit</span><code className="ws-mono">{data.commit ? data.commit.slice(0, 12) : '—'}</code></div>
            <div className="ws-detail-row"><span>Ahead</span><strong>{data.ahead}</strong></div>
            <div className="ws-detail-row"><span>Behind</span><strong>{data.behind}</strong></div>
          </div>
        </Card>
        <Card>
          <div className="section-header-lite"><h3>Changed Files</h3></div>
          {data.changedFiles?.length ? (
            <div className="ws-detail-list">
              {data.changedFiles.map((f) => (
                <div key={f} className="ws-detail-row">
                  <code className="ws-mono" style={{ fontSize: '0.85rem' }}>{f}</code>
                </div>
              ))}
            </div>
          ) : (
            <EmptyState title="No changes" description="Working tree is clean." />
          )}
        </Card>
      </div>

      <Card>
        <div className="section-header-lite"><h3>Recent Commits</h3></div>
        {data.recentCommits?.length ? (
          <div className="ws-commit-list">
            {data.recentCommits.map((c, i) => {
              const parts = c.split(' ')
              const hash = parts[0]
              const msg = parts.slice(1).join(' ')
              return (
                <div key={i} className="ws-commit-row">
                  <code className="ws-mono ws-commit-hash">{hash}</code>
                  <span className="ws-commit-msg">{msg}</span>
                </div>
              )
            })}
          </div>
        ) : (
          <EmptyState title="No commits" description="No recent commit history available." />
        )}
      </Card>
    </div>
  )
}
